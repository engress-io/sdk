package metrics

import (
	"os"
	"sync"

	"github.com/DataDog/datadog-go/v5/statsd"
)

// Client emits DogStatsD metrics when DD_AGENT_HOST is set; otherwise no-op.
type Client struct {
	statsd *statsd.Client
}

var (
	global   *Client
	initOnce sync.Once
)

func Init(service string) *Client {
	initOnce.Do(func() {
		host := os.Getenv("DD_AGENT_HOST")
		if host == "" {
			global = &Client{}
			return
		}
		port := os.Getenv("DD_DOGSTATSD_PORT")
		if port == "" {
			port = "8125"
		}
		env := os.Getenv("DD_ENV")
		if env == "" {
			env = "production"
		}
		tags := []string{"service:" + service, "env:" + env}
		c, err := statsd.New(host+":"+port, statsd.WithTags(tags))
		if err != nil {
			global = &Client{}
			return
		}
		global = &Client{statsd: c}
	})
	return global
}

func Global() *Client {
	if global == nil {
		return Init("")
	}
	return global
}

func (c *Client) enabled() bool {
	return c != nil && c.statsd != nil
}

func (c *Client) Gauge(name string, value float64, tags ...string) {
	if !c.enabled() {
		return
	}
	_ = c.statsd.Gauge("engress."+name, value, tags, 1)
}

func (c *Client) Count(name string, delta int64, tags ...string) {
	if !c.enabled() {
		return
	}
	_ = c.statsd.Count("engress."+name, delta, tags, 1)
}

func (c *Client) Histogram(name string, value float64, tags ...string) {
	if !c.enabled() {
		return
	}
	_ = c.statsd.Histogram("engress."+name, value, tags, 1)
}

func (c *Client) Close() error {
	if !c.enabled() {
		return nil
	}
	return c.statsd.Close()
}
