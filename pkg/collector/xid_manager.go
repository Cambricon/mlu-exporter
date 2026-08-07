package collector

import (
	"sort"
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
	stopCh   chan struct{}
	mluInfo  *MLUStatMap
	done     chan struct{}
}

var globalXIDManager *XIDEventManager
var managerOnce sync.Once

func GetXIDEventManager(cli cndev.Cndev, mluInfo *MLUStatMap) *XIDEventManager {
	managerOnce.Do(func() {
		globalXIDManager = &XIDEventManager{
			cndevcli: cli,
			handlers: make([]XIDEventHandler, 0),
			ch:       make(chan cndev.XIDInfoWithTimestamp, 10),
			stopCh:   make(chan struct{}),
			mluInfo:  mluInfo,
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
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.started {
		return
	}
	m.slots = slots
}

// UpdateSlotsFromMLUInfo recalculates slots from the current mluInfo,
// filtering out devices where xidCallback is disabled.
func (m *XIDEventManager) UpdateSlotsFromMLUInfo() {
	slots := []int{}
	m.mluInfo.Range(func(_ string, stat MLUStat) bool {
		if !stat.cndevInterfaceDisabled["xidCallbackDisabled"] {
			slots = append(slots, int(stat.slot))
		}
		return true
	})
	sort.Ints(slots)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.slots = slots
}

func (m *XIDEventManager) Start() error {
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

	// wait until cndev handle is valid (CndevMu write lock released means Init completed)
	// safe to hold mu.Lock while taking CndevMu.RLock: no code path does CndevMu.Lock→mu.Lock
	m.mluInfo.CndevMu.RLock()
	err := m.cndevcli.RegisterEventsHandleAndWait(m.slots, m.ch, m.stopCh)
	m.mluInfo.CndevMu.RUnlock()

	if err != nil {
		log.Errorln(errors.Wrap(err, "register event handle"))
		return err
	}
	log.Debug("XIDEventManager register event handle finished")
	m.started = true
	m.done = make(chan struct{})
	go m.eventLoop()
	return nil
}

func (m *XIDEventManager) ReRegister() {
	// update slots from current device topology before re-registering
	m.UpdateSlotsFromMLUInfo()

	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.slots) == 0 {
		return
	}
	// stop existing waitEvents and eventLoop goroutines via stopCh
	// must release mu before waiting: old eventLoop uses mu.RLock to copy handlers,
	// holding mu.Lock would deadlock
	if m.started {
		close(m.stopCh) // signal both waitEvents and eventLoop to stop
		m.mu.Unlock()
		<-m.done // wait for eventLoop to finish
		// waitEvents exits asynchronously (within eventWaitTimeout),
		// it checks stopCh before processing events, so no risk of
		// writing to old ch or conflicting with new goroutines
		m.mu.Lock()
	}

	m.started = false
	m.ch = make(chan cndev.XIDInfoWithTimestamp, 10)
	m.stopCh = make(chan struct{})

	// safe to hold mu.Lock while taking CndevMu.RLock: no code path does CndevMu.Lock→mu.Lock
	m.mluInfo.CndevMu.RLock()
	err := m.cndevcli.RegisterEventsHandleAndWait(m.slots, m.ch, m.stopCh)
	m.mluInfo.CndevMu.RUnlock()

	if err != nil {
		log.Errorln(errors.Wrap(err, "re-register event handle"))
		return
	}
	log.Debug("XIDEventManager re-register event handle finished")
	m.started = true
	m.done = make(chan struct{})
	go m.eventLoop()
}

func (m *XIDEventManager) eventLoop() {
	for {
		select {
		case <-m.stopCh:
			close(m.done)
			return
		case event, ok := <-m.ch:
			if !ok {
				close(m.done)
				return
			}
			m.mu.RLock()
			handlers := make([]XIDEventHandler, len(m.handlers))
			copy(handlers, m.handlers)
			m.mu.RUnlock()

			for _, handler := range handlers {
				handler.HandleXIDEvent(event)
			}
		}
	}
}
