// Package ports maps well-known port numbers to human-readable service names
// and subdomain prefixes. Shared by the agent (for display) and the edge
// (for validation).
package ports

// Service describes a recognized network service.
type Service struct {
	// Name is the human-readable service name (e.g. "HTTPS", "RDP").
	Name string
	// Prefix is the subdomain prefix used in tunnel URLs (e.g. "https-", "rdp-").
	Prefix string
}

// table maps a port number to its service description.
// Covers IANA well-known ports plus common game, database, and app ports.
var table = map[int]Service{
	// Web
	80:   {"HTTP", "http-"},
	443:  {"HTTPS", "https-"},
	8080: {"HTTP Alt", "http-"},
	8443: {"HTTPS Alt", "https-"},
	8000: {"HTTP Dev", "http-"},
	9000: {"HTTP Dev", "http-"},
	3000: {"Web App", "web-"},
	5000: {"Web App", "web-"},
	5173: {"Vite Dev", "web-"},
	8081: {"Web App", "web-"},
	8888: {"Web App", "web-"},
	9090: {"Web App", "web-"},

	// Remote access
	22:   {"SSH", "ssh-"},
	3389: {"RDP", "rdp-"},
	5900: {"VNC", "vnc-"},
	5901: {"VNC", "vnc-"},
	5938: {"TeamViewer", "tv-"},

	// Databases
	3306:  {"MySQL", "mysql-"},
	5432:  {"PostgreSQL", "pg-"},
	27017: {"MongoDB", "mongo-"},
	6379:  {"Redis", "redis-"},
	9200:  {"Elasticsearch", "es-"},
	5601:  {"Kibana", "kibana-"},
	11211: {"Memcached", "cache-"},

	// Mail
	25:  {"SMTP", "smtp-"},
	465: {"SMTPS", "smtp-"},
	587: {"SMTP Submission", "smtp-"},
	993: {"IMAPS", "imap-"},
	995: {"POP3S", "pop3-"},

	// DNS / NTP
	53:  {"DNS", "dns-"},
	123: {"NTP", "ntp-"},
	161: {"SNMP", "snmp-"},

	// Directory
	389: {"LDAP", "ldap-"},
	636: {"LDAPS", "ldaps-"},
	992: {"Telnet TLS", "telnet-"},

	// Games
	25565: {"Minecraft", "mc-"},
	25575: {"Minecraft RCON", "mcrcon-"},
	27015: {"Source Engine", "game-"},
	27016: {"Source Engine", "game-"},
	7777:  {"Game Server", "game-"},
	7778:  {"Game Server", "game-"},
	3074:  {"Xbox Live", "xbox-"},
	3724:  {"World of Warcraft", "wow-"},
	6112:  {"Battle.net", "bnet-"},
	8085:  {"World of Warcraft", "wow-"},
}

// Lookup returns the service for a given port, plus a boolean indicating
// whether the port is recognized. Unrecognized ports get ("Service", "tunnel-").
func Lookup(port int) (Service, bool) {
	if s, ok := table[port]; ok {
		return s, true
	}
	return Service{Name: "Service", Prefix: "tunnel-"}, false
}

// PrefixFor returns the subdomain prefix for a port.
func PrefixFor(port int) string {
	s, _ := Lookup(port)
	return s.Prefix
}
