package collector

import (
	"sync"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type XIDEventHandler interface {
	HandleXIDEvent(event cndev.XIDInfoWithTimestamp)
}

type XIDEventManager struct {
	cndevcli cndev.Cndev
	handlers []XIDEventHandler
	mu       sync.RWMutex
	started  bool
	slots    []int
	ch       chan cndev.XIDInfoWithTimestamp
}

var globalXIDManager *XIDEventManager
var managerOnce sync.Once

func GetXIDEventManager(cli cndev.Cndev) *XIDEventManager {
	managerOnce.Do(func() {
		globalXIDManager = &XIDEventManager{
			cndevcli: cli,
			handlers: make([]XIDEventHandler, 0),
			ch:       make(chan cndev.XIDInfoWithTimestamp, 10),
		}
	})
	return globalXIDManager
}

func (m *XIDEventManager) RegisterHandler(handler XIDEventHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers = append(m.handlers, handler)
}

func (m *XIDEventManager) SetSlots(slots []int) {
	if len(m.slots) > 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.slots = slots
}

func (m *XIDEventManager) Start() error {
	var err error
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.started {
		log.Debug("XIDEventManager already started")
		return nil
	}

	if len(m.slots) == 0 {
		log.Debug("XIDEventManager no slots to monitor")
		return nil
	}

	if err = m.cndevcli.RegisterEventsHandleAndWait(m.slots, m.ch); err != nil {
		log.Errorln(errors.Wrap(err, "register event handle"))
		return err
	}
	log.Debug("XIDEventManager register event handle finished")
	m.started = true
	go m.eventLoop()
	return err
}

func (m *XIDEventManager) eventLoop() {
	for event := range m.ch {
		m.mu.RLock()
		handlers := make([]XIDEventHandler, len(m.handlers))
		copy(handlers, m.handlers)
		m.mu.RUnlock()

		for _, handler := range handlers {
			handler.HandleXIDEvent(event)
		}
	}
}
