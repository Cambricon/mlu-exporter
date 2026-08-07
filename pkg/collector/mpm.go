package collector

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func init() {
	registerCollector(Mpm, NewMpmCollector)
}

type mpmCollector struct {
	baseInfo   BaseInfo
	client     cndev.Cndev
	fnMap      map[string]interface{}
	metrics    metrics.CollectorMetrics
	sharedInfo *MLUStatMap

	mpmIntervalMS uint
	metricIDs     []cndev.MpmMetricID
	cachedResults map[uint]map[string]float64
	cacheMu       sync.RWMutex
	stopCh        chan struct{}
	doneCh        chan struct{}
}

func NewMpmCollector(m metrics.CollectorMetrics, bi BaseInfo) Collector {
	c := &mpmCollector{
		baseInfo:      bi,
		client:        bi.cndevClient,
		metrics:       m,
		mpmIntervalMS: bi.mpmIntervalMS,
		metricIDs:     buildMetricIDsFrom(m),
	}

	c.fnMap = map[string]interface{}{
		MpmIPUUtil:              c.collectMpmIPUUtil,
		MpmMLUUtil:              c.collectMpmMLUUtil,
		MpmTensorUtil:           c.collectMpmTensorUtil,
		MpmPCIeTxPerSec:         c.collectMpmPCIeTxPerSec,
		MpmPCIeRxPerSec:         c.collectMpmPCIeRxPerSec,
		MpmMLULinkTotalTxPerSec: c.collectMpmMLULinkTotalTxPerSec,
		MpmMLULinkTotalRxPerSec: c.collectMpmMLULinkTotalRxPerSec,
		MpmMLULinkTxPerSec:      c.collectMpmMLULinkTxBandwidth,
		MpmMLULinkRxPerSec:      c.collectMpmMLULinkRxBandwidth,
	}

	return c
}

var metricKeyToCIdMap = map[string]cndev.MpmMetricID{
	MpmIPUUtil:              cndev.MpmMetricIPUUtil,
	MpmMLUUtil:              cndev.MpmMetricMLUUtil,
	MpmTensorUtil:           cndev.MpmMetricTensorUtil,
	MpmPCIeTxPerSec:         cndev.MpmMetricPCIeTxPerSec,
	MpmPCIeRxPerSec:         cndev.MpmMetricPCIeRxPerSec,
	MpmMLULinkTotalTxPerSec: cndev.MpmMetricMLULinkTotalTxPerSec,
	MpmMLULinkTotalRxPerSec: cndev.MpmMetricMLULinkTotalRxPerSec,
}

var perLinkMetricKeys = map[string]bool{
	MpmMLULinkTxPerSec: true,
	MpmMLULinkRxPerSec: true,
}

func buildMetricIDsFrom(m metrics.CollectorMetrics) []cndev.MpmMetricID {
	var ids []cndev.MpmMetricID
	hasPerLink := false
	for key := range m {
		if id, ok := metricKeyToCIdMap[key]; ok {
			ids = append(ids, id)
			continue
		}
		if perLinkMetricKeys[key] {
			hasPerLink = true
		}
	}
	if hasPerLink {
		for i := 0; i < cndev.MpmMaxLinks; i++ {
			ids = append(ids, cndev.MpmLinkTxMetricID(i))
			ids = append(ids, cndev.MpmLinkRxMetricID(i))
		}
	}
	return ids
}

func (c *mpmCollector) init(info *MLUStatMap) error {
	c.sharedInfo = info

	metricIDs := c.getCachedMetricIDs()
	if len(metricIDs) == 0 {
		return nil
	}

	// Take initial sample for each device so the first Prometheus scrape
	// (which arrives >100ms later) can produce valid MPM metrics.
	info.Range(func(_ string, stat MLUStat) bool {
		if stat.mpmDisabled {
			return true
		}
		if _, err := c.client.MpmCollect(stat.slot, metricIDs); err != nil {
			log.Debugf("Initial MpmCollect for slot %d: %v", stat.slot, err)
		}
		return true
	})

	return nil
}

func (c *mpmCollector) updateMetrics(m metrics.CollectorMetrics) {
	c.cacheMu.Lock()
	c.metrics = m
	c.metricIDs = buildMetricIDsFrom(m)
	c.cacheMu.Unlock()
}

func (c *mpmCollector) getCachedMetricIDs() []cndev.MpmMetricID {
	c.cacheMu.RLock()
	ids := c.metricIDs
	c.cacheMu.RUnlock()
	return ids
}

func (c *mpmCollector) start() {
	if c.mpmIntervalMS == 0 {
		return
	}
	metricIDs := c.getCachedMetricIDs()
	if len(metricIDs) == 0 {
		return
	}

	// Collect immediately to seed the cache so the first scrape has data.
	c.collectAndCache(metricIDs)

	c.stopCh = make(chan struct{})
	c.doneCh = make(chan struct{})
	ticker := time.NewTicker(time.Duration(c.mpmIntervalMS) * time.Millisecond)
	log.Infof("MPM background timer started with interval %dms", c.mpmIntervalMS)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				metricIDs := c.getCachedMetricIDs()
				if len(metricIDs) > 0 {
					c.collectAndCache(metricIDs)
				}
			case <-c.stopCh:
				log.Info("MPM background timer stopped")
				c.doneCh <- struct{}{}
				return
			}
		}
	}()
}

func (c *mpmCollector) stop() {
	if c.stopCh != nil {
		close(c.stopCh)
		<-c.doneCh
		c.stopCh = nil
		c.doneCh = nil
	}
}

func (c *mpmCollector) collectAndCache(metricIDs []cndev.MpmMetricID) {
	c.sharedInfo.CndevMu.RLock()
	defer c.sharedInfo.CndevMu.RUnlock()

	newResults := make(map[uint]map[string]float64)

	c.sharedInfo.Range(func(_ string, stat MLUStat) bool {
		if stat.mpmDisabled {
			return true
		}

		results, err := c.client.MpmCollect(stat.slot, metricIDs)
		if err != nil {
			log.Errorf("MpmCollect failed for slot %d: %v", stat.slot, err)
			return true
		}
		if results == nil {
			return true
		}

		metricResults := make(map[string]float64)
		for _, r := range results {
			if r.Ret != 0 {
				log.Debugf("MPM metric ID %d ret=%d for slot %d", r.MetricID, r.Ret, stat.slot)
				continue
			}
			if math.IsNaN(r.Value) {
				log.Debugf("MPM metric ID %d returned NaN for slot %d, skipping", r.MetricID, stat.slot)
				continue
			}
			metricKey, linkIdx := cIDToMetricKeyAndLink(r.MetricID)
			if metricKey != "" {
				if linkIdx >= 0 {
					metricResults[fmt.Sprintf("%s:%d", metricKey, linkIdx)] = r.Value
				} else {
					metricResults[metricKey] = r.Value
				}
			}
		}
		newResults[stat.slot] = metricResults
		return true
	})

	c.cacheMu.Lock()
	c.cachedResults = newResults
	c.cacheMu.Unlock()
}

func (c *mpmCollector) collect(ch chan<- prometheus.Metric) {
	c.sharedInfo.CndevMu.RLock()
	defer c.sharedInfo.CndevMu.RUnlock()

	metricIDs := c.getCachedMetricIDs()
	if len(metricIDs) == 0 {
		return
	}

	var deviceResults map[uint]map[string]float64

	if c.mpmIntervalMS > 0 {
		c.cacheMu.RLock()
		deviceResults = c.cachedResults
		c.cacheMu.RUnlock()
		if deviceResults == nil {
			return
		}
	} else {
		deviceResults = make(map[uint]map[string]float64)

		c.sharedInfo.Range(func(_ string, stat MLUStat) bool {
			if stat.mpmDisabled {
				return true
			}

			results, err := c.client.MpmCollect(stat.slot, metricIDs)
			if err != nil {
				log.Errorf("MpmCollect failed for slot %d: %v", stat.slot, err)
				return true
			}
			if results == nil {
				return true
			}

			metricResults := make(map[string]float64)
			for _, r := range results {
				if r.Ret != 0 {
					log.Debugf("MPM metric ID %d ret=%d for slot %d", r.MetricID, r.Ret, stat.slot)
					continue
				}
				if math.IsNaN(r.Value) {
					log.Debugf("MPM metric ID %d returned NaN for slot %d, skipping", r.MetricID, stat.slot)
					continue
				}
				metricKey, linkIdx := cIDToMetricKeyAndLink(r.MetricID)
				if metricKey != "" {
					if linkIdx >= 0 {
						metricResults[fmt.Sprintf("%s:%d", metricKey, linkIdx)] = r.Value
					} else {
						metricResults[metricKey] = r.Value
					}
				}
			}
			deviceResults[stat.slot] = metricResults
			return true
		})
	}

	for name, m := range c.metrics {
		fn, ok := c.fnMap[name]
		if !ok {
			continue
		}
		f, ok := fn.(func(chan<- prometheus.Metric, metrics.Metric, map[uint]map[string]float64))
		if !ok {
			log.Warnf("type assertion for fn %s failed, skip", name)
		} else {
			f(ch, m, deviceResults)
		}
	}
}

func cIDToMetricKeyAndLink(id cndev.MpmMetricID) (string, int) {
	for key, mid := range metricKeyToCIdMap {
		if mid == id {
			return key, -1
		}
	}
	// Per-link MLULink metric IDs start at 62 (L0Tx) and follow the pattern:
	// L{n}Tx = 62 + 2*n, L{n}Rx = 62 + 2*n + 1 (see cndev.h CNDEV_MPM_METRIC_MLULINK_L0_TX_PER_SEC = 62)
	linkIdx := int(id-62) / 2
	if linkIdx >= 0 && linkIdx < cndev.MpmMaxLinks {
		if int(id-62)%2 == 0 {
			return MpmMLULinkTxPerSec, linkIdx
		}
		return MpmMLULinkRxPerSec, linkIdx
	}
	return "", -1
}

func (c *mpmCollector) collectMpmMetric(ch chan<- prometheus.Metric, m metrics.Metric, deviceResults map[uint]map[string]float64, metricKey string) {
	c.sharedInfo.Range(func(_ string, stat MLUStat) bool {
		if stat.mpmDisabled {
			return true
		}
		results, ok := deviceResults[stat.slot]
		if !ok {
			return true
		}
		value, ok := results[metricKey]
		if !ok {
			return true
		}
		lv := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, hostIP: c.baseInfo.hostIP})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, value, lv...)
		return true
	})
}

func (c *mpmCollector) collectMpmIPUUtil(ch chan<- prometheus.Metric, m metrics.Metric, deviceResults map[uint]map[string]float64) {
	c.collectMpmMetric(ch, m, deviceResults, MpmIPUUtil)
}

func (c *mpmCollector) collectMpmMLUUtil(ch chan<- prometheus.Metric, m metrics.Metric, deviceResults map[uint]map[string]float64) {
	c.collectMpmMetric(ch, m, deviceResults, MpmMLUUtil)
}

func (c *mpmCollector) collectMpmTensorUtil(ch chan<- prometheus.Metric, m metrics.Metric, deviceResults map[uint]map[string]float64) {
	c.collectMpmMetric(ch, m, deviceResults, MpmTensorUtil)
}

func (c *mpmCollector) collectMpmPCIeTxPerSec(ch chan<- prometheus.Metric, m metrics.Metric, deviceResults map[uint]map[string]float64) {
	c.collectMpmMetric(ch, m, deviceResults, MpmPCIeTxPerSec)
}

func (c *mpmCollector) collectMpmPCIeRxPerSec(ch chan<- prometheus.Metric, m metrics.Metric, deviceResults map[uint]map[string]float64) {
	c.collectMpmMetric(ch, m, deviceResults, MpmPCIeRxPerSec)
}

func (c *mpmCollector) collectMpmMLULinkTotalTxPerSec(ch chan<- prometheus.Metric, m metrics.Metric, deviceResults map[uint]map[string]float64) {
	c.collectMpmMetric(ch, m, deviceResults, MpmMLULinkTotalTxPerSec)
}

func (c *mpmCollector) collectMpmMLULinkTotalRxPerSec(ch chan<- prometheus.Metric, m metrics.Metric, deviceResults map[uint]map[string]float64) {
	c.collectMpmMetric(ch, m, deviceResults, MpmMLULinkTotalRxPerSec)
}

func (c *mpmCollector) collectMpmMLULinkBandwidth(ch chan<- prometheus.Metric, m metrics.Metric, deviceResults map[uint]map[string]float64) {
	c.sharedInfo.Range(func(_ string, stat MLUStat) bool {
		if stat.mpmDisabled {
			return true
		}
		results, ok := deviceResults[stat.slot]
		if !ok {
			return true
		}
		for i := 0; i < stat.link; i++ {
			valueKey := fmt.Sprintf("%s:%d", m.Name, i)
			value, ok := results[valueKey]
			if !ok {
				continue
			}
			ppi := ""
			if stat.linkPPI != nil {
				ppi = stat.linkPPI[i]
			}
			lv := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, hostIP: c.baseInfo.hostIP, link: i, ppi: ppi})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, value, lv...)
		}
		return true
	})
}

func (c *mpmCollector) collectMpmMLULinkTxBandwidth(ch chan<- prometheus.Metric, m metrics.Metric, deviceResults map[uint]map[string]float64) {
	c.collectMpmMLULinkBandwidth(ch, m, deviceResults)
}

func (c *mpmCollector) collectMpmMLULinkRxBandwidth(ch chan<- prometheus.Metric, m metrics.Metric, deviceResults map[uint]map[string]float64) {
	c.collectMpmMLULinkBandwidth(ch, m, deviceResults)
}
