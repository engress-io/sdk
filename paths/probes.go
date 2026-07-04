package paths

import (
	"path"
	"strings"
)

// probeRule names a blocked scanner path pattern for logging.
type probeRule struct {
	name string
	kind int // 0 exact, 1 prefix, 2 suffix
	pat  string
}

// probeRules are paths commonly requested by automated vulnerability scanners.
// Legitimate AI tunnel paths (/v1/*, /api/tags, /api/generate) are not listed.
var probeRules = []probeRule{
	{name: ".env", kind: 1, pat: "/.env"},
	{name: "config.js", kind: 0, pat: "/config.js"},
	{name: "settings.js", kind: 0, pat: "/settings.js"},
	{name: "env.js", kind: 0, pat: "/js/env.js"},
	{name: "api/config", kind: 0, pat: "/api/config"},
	{name: "wp-admin", kind: 1, pat: "/wp-admin"},
	{name: "wp-login", kind: 0, pat: "/wp-login.php"},
	{name: "phpmyadmin", kind: 1, pat: "/phpmyadmin"},
	{name: "git", kind: 1, pat: "/.git/"},
	{name: "actuator", kind: 1, pat: "/actuator/"},
	{name: "server-status", kind: 0, pat: "/server-status"},
	{name: "debug", kind: 1, pat: "/debug/"},
	{name: "aws-config", kind: 1, pat: "/.aws/"},
	{name: "docker-config", kind: 1, pat: "/.docker/"},
	{name: "DS_Store", kind: 0, pat: "/.DS_Store"},
}

const (
	matchExact  = 0
	matchPrefix = 1
)

// IsProbe reports whether path looks like an automated scanner probe.
// path should be the URL path only (no query string).
func IsProbe(raw string) (bool, string) {
	p := path.Clean(raw)
	if p == "." {
		p = "/"
	}
	lower := strings.ToLower(p)
	for _, rule := range probeRules {
		switch rule.kind {
		case matchExact:
			if lower == rule.pat {
				return true, rule.name
			}
		case matchPrefix:
			if lower == rule.pat || strings.HasPrefix(lower, rule.pat) {
				return true, rule.name
			}
		}
	}
	return false, ""
}
