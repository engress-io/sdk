package ports

// reason documents why a port is blocked.
type reason struct {
	name   string
	reason string
}

// blockedPorts are never allowed through the platform. Enforced at the
// edge during tunnel registration.
var blockedPorts = map[int]reason{
	// VPN tunneling
	1194:  {"OpenVPN", "VPN tunneling"},
	1195:  {"OpenVPN", "VPN tunneling"},
	51820: {"WireGuard", "VPN tunneling"},
	51821: {"WireGuard", "VPN tunneling"},
	500:   {"IPsec", "VPN tunneling"},
	4500:  {"IPsec NAT-T", "VPN tunneling"},
	1701:  {"L2TP", "VPN tunneling"},
	1723:  {"PPTP", "VPN tunneling"},

	// Anonymizers
	9001: {"Tor", "anonymizer"},
	9030: {"Tor Directory", "anonymizer"},
	9050: {"Tor SOCKS", "anonymizer"},
	9051: {"Tor Control", "anonymizer"},
	9150: {"Tor Browser", "anonymizer"},

	// P2P file sharing
	6881:  {"BitTorrent", "P2P file sharing"},
	6882:  {"BitTorrent", "P2P file sharing"},
	6883:  {"BitTorrent", "P2P file sharing"},
	6884:  {"BitTorrent", "P2P file sharing"},
	6885:  {"BitTorrent", "P2P file sharing"},
	6886:  {"BitTorrent", "P2P file sharing"},
	6887:  {"BitTorrent", "P2P file sharing"},
	6888:  {"BitTorrent", "P2P file sharing"},
	6889:  {"BitTorrent", "P2P file sharing"},
	6969:  {"BitTorrent Tracker", "P2P file sharing"},
	49001: {"BitTorrent DHT", "P2P file sharing"},

	// Spam relay
	25:  {"SMTP", "spam relay"},
	465: {"SMTPS", "spam relay"},
	587: {"SMTP Submission", "spam relay"},

	// Amplification / unauthenticated data stores
	11211: {"Memcached", "amplification risk"},
	6379:  {"Redis", "unauthenticated data"},
	27017: {"MongoDB", "unauthenticated data"},
	9200:  {"Elasticsearch", "unauthenticated data"},
	3306:  {"MySQL", "unauthenticated data"},
	5432:  {"PostgreSQL", "unauthenticated data"},
}

// IsBlocked reports whether a port is blocked and returns the service name
// and reason for blocking.
func IsBlocked(port int) (bool, string) {
	if r, ok := blockedPorts[port]; ok {
		return true, r.name
	}
	return false, ""
}

// BlockedPorts returns a copy of the blocked ports map for admin display.
func BlockedPorts() map[int]reason {
	out := make(map[int]reason, len(blockedPorts))
	for p, r := range blockedPorts {
		out[p] = r
	}
	return out
}
