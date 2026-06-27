package stats

import (
	"sync"
	"sync/atomic"
	"time"
)

var nextRequestID atomic.Uint64

func NextRequestID() uint64 {
	return nextRequestID.Add(1)
}

const maxLogEntries = 200

// RequestStatusRunning marks an in-flight request (not yet completed).
const RequestStatusRunning = 0

type Collector struct {
	startedAt    time.Time
	requestCount atomic.Uint64
	errorCount   atomic.Uint64
	bytesIn      atomic.Uint64
	bytesOut     atomic.Uint64
	totalLatency atomic.Uint64

	logMu sync.RWMutex
	log   []RequestEntry
}

type RequestEntry struct {
	ID          uint64            `json:"id"`
	Time        time.Time         `json:"time"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Subdomain   string            `json:"subdomain,omitempty"`
	Host        string            `json:"host,omitempty"`
	RemoteAddr  string            `json:"remote_addr,omitempty"`
	RequestURI  string            `json:"request_uri,omitempty"`
	UserAgent   string            `json:"user_agent,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	ReqHeaders  map[string]string `json:"req_headers,omitempty"`
	ReqBody     string            `json:"req_body,omitempty"`
	Status      int               `json:"status"`
	LatencyMs   float64           `json:"latency_ms"`
	BytesIn     uint64            `json:"bytes_in"`
	BytesOut    uint64            `json:"bytes_out"`
	Error       string            `json:"error,omitempty"`
	Streaming   bool              `json:"streaming,omitempty"`
}

func New() *Collector {
	return &Collector{startedAt: time.Now()}
}

func (c *Collector) Record(entry RequestEntry) {
	c.CompleteRequest(entry)
}

// BeginRequest tracks a request as soon as it arrives (before the response completes).
func (c *Collector) BeginRequest(entry RequestEntry) {
	entry.Status = RequestStatusRunning
	c.logMu.Lock()
	if i := c.indexByID(entry.ID); i >= 0 {
		c.log[i] = entry
	} else {
		c.log = append(c.log, entry)
		if len(c.log) > maxLogEntries {
			c.log = c.log[len(c.log)-maxLogEntries:]
		}
	}
	c.logMu.Unlock()
}

// CompleteRequest finalizes a request and updates aggregate counters.
func (c *Collector) CompleteRequest(entry RequestEntry) {
	c.logMu.Lock()
	if i := c.indexByID(entry.ID); i >= 0 {
		c.log[i] = entry
	} else {
		c.log = append(c.log, entry)
		if len(c.log) > maxLogEntries {
			c.log = c.log[len(c.log)-maxLogEntries:]
		}
	}
	c.logMu.Unlock()

	c.requestCount.Add(1)
	c.bytesIn.Add(entry.BytesIn)
	c.bytesOut.Add(entry.BytesOut)
	c.totalLatency.Add(uint64(entry.LatencyMs))
	if entry.Status >= 400 || entry.Error != "" {
		c.errorCount.Add(1)
	}
}

func (c *Collector) indexByID(id uint64) int {
	if id == 0 {
		return -1
	}
	for i := range c.log {
		if c.log[i].ID == id {
			return i
		}
	}
	return -1
}

type Snapshot struct {
	Name          string         `json:"name"`
	Domain        string         `json:"domain"`
	IP            string         `json:"ip"`
	LocalAddr     string         `json:"local_addr"`
	ProxyPort     int            `json:"proxy_port"`
	DashboardPort int            `json:"dashboard_port"`
	Uptime        string         `json:"uptime"`
	RequestCount  uint64         `json:"request_count"`
	ErrorCount    uint64         `json:"error_count"`
	BytesIn       uint64         `json:"bytes_in"`
	BytesOut      uint64         `json:"bytes_out"`
	AvgLatencyMs  float64        `json:"avg_latency_ms"`
	Recent        []RequestEntry `json:"recent"`
}

type Info struct {
	Name          string
	Domain        string
	IP            string
	LocalAddr     string
	ProxyPort     int
	DashboardPort int
}

var infoMu sync.RWMutex
var info Info

func SetInfo(i Info) {
	infoMu.Lock()
	info = i
	infoMu.Unlock()
}

func (c *Collector) Snapshot() Snapshot {
	infoMu.RLock()
	i := info
	infoMu.RUnlock()

	count := c.requestCount.Load()
	var avg float64
	if count > 0 {
		avg = float64(c.totalLatency.Load()) / float64(count)
	}

	c.logMu.RLock()
	recent := make([]RequestEntry, len(c.log))
	copy(recent, c.log)
	c.logMu.RUnlock()

	// newest first for dashboard
	for i, j := 0, len(recent)-1; i < j; i, j = i+1, j-1 {
		recent[i], recent[j] = recent[j], recent[i]
	}

	return Snapshot{
		Name:          i.Name,
		Domain:        i.Domain,
		IP:            i.IP,
		LocalAddr:     i.LocalAddr,
		ProxyPort:     i.ProxyPort,
		DashboardPort: i.DashboardPort,
		Uptime:        time.Since(c.startedAt).Round(time.Second).String(),
		RequestCount:  count,
		ErrorCount:    c.errorCount.Load(),
		BytesIn:       c.bytesIn.Load(),
		BytesOut:      c.bytesOut.Load(),
		AvgLatencyMs:  avg,
		Recent:        recent,
	}
}
