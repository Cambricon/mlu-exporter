package collector

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/pkg/errors"
	"github.com/prometheus/common/expfmt"
	log "github.com/sirupsen/logrus"
)

type HTTPClient struct {
	url        *url.URL
	client     *http.Client
	sharedInfo map[string]MLUStat
	labels     []string
	host       string
	metricName string
	jobName    string
}

func NewClient(ur string, sharedInfo map[string]MLUStat, labels []string, host, metricName, jobName string) (c *HTTPClient, err error) {
	u, err := url.Parse(ur)
	if err != nil {
		return
	}
	c = &HTTPClient{
		url:        u,
		client:     &http.Client{},
		sharedInfo: sharedInfo,
		labels:     labels,
		host:       host,
		metricName: metricName,
		jobName:    jobName,
	}
	return
}

func (c *HTTPClient) PushWithRetries(event cndev.XIDInfoWithTimestamp, retries int) (err error) {
	slotInfo := map[uint]MLUStat{}
	for _, stat := range c.sharedInfo {
		slotInfo[stat.slot] = stat
	}

	labels := []string{}
	labelValues := getLabelValues(c.labels, labelInfo{
		stat:    slotInfo[event.Device],
		host:    c.host,
		xidInfo: event,
	})

	labels = append(labels, fmt.Sprintf(`job="%s"`, c.jobName))
	for i, label := range c.labels {
		labels = append(labels, fmt.Sprintf(`%s="%s"`, label, labelValues[i]))
	}

	// event.Timestamp is microseconds, we need to convert it to milliseconds
	metrics := fmt.Sprintf("%s{%s} %d %d", c.metricName, strings.Join(labels, ","), 1, event.Timestamp/1000)
	return c.sendWithRetries(c.url.String(), metrics, retries)
}

func (c *HTTPClient) sendWithRetries(url, body string, retries int) error {
	var err error
	for i := 0; i < retries; i++ {
		var req *http.Request
		req, err = http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
		if err != nil {
			log.Errorln(errors.Wrap(err, "new request"))
			continue
		}
		req.Header.Set("Content-Type", string(expfmt.FmtText))

		resp, err := c.client.Do(req)
		if err != nil {
			log.Errorln(errors.Wrap(err, "do request"))
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode/100 != 2 {
			var b []byte
			b, err = io.ReadAll(resp.Body)
			if err != nil {
				log.Errorln(errors.Wrap(err, "read body"))
				continue
			}
			log.Errorln(fmt.Sprintf("send failed: %s", b))
		}
		return nil
	}
	return err
}
