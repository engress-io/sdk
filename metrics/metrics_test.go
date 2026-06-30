package metrics

import "testing"

func TestClientNoOpWithoutAgent(t *testing.T) {
	c := Init("test")
	c.Gauge("tunnels.active", 1)
	c.Count("tunnels.created", 1, "protocol:tcp")
	c.Histogram("requests.latency", 12.5, "status:200")
}
