package paths

import "testing"

func TestIsProbe_blocked(t *testing.T) {
	blocked := []struct {
		path string
		want string
	}{
		{"/.env", ".env"},
		{"/.env.local", ".env"},
		{"/.env.production", ".env"},
		{"//.env", ".env"},
		{"/config.js", "config.js"},
		{"/settings.js", "settings.js"},
		{"/js/env.js", "env.js"},
		{"/api/config", "api/config"},
		{"/wp-admin", "wp-admin"},
		{"/wp-admin/setup.php", "wp-admin"},
		{"/wp-login.php", "wp-login"},
		{"/phpmyadmin", "phpmyadmin"},
		{"/.git/config", "git"},
		{"/actuator/health", "actuator"},
		{"/server-status", "server-status"},
		{"/debug/pprof", "debug"},
	}
	for _, tc := range blocked {
		t.Run(tc.path, func(t *testing.T) {
			ok, name := IsProbe(tc.path)
			if !ok {
				t.Fatalf("IsProbe(%q) = false, want true", tc.path)
			}
			if name != tc.want {
				t.Fatalf("name = %q, want %q", name, tc.want)
			}
		})
	}
}

func TestIsProbe_allowed(t *testing.T) {
	allowed := []string{
		"/",
		"/v1/models",
		"/v1/chat/completions",
		"/api/tags",
		"/api/generate",
		"/api/chat",
		"/api/health",
		"/my.env.js",
		"/env",
		"/debugging",
		"/healthz",
	}
	for _, p := range allowed {
		t.Run(p, func(t *testing.T) {
			if ok, _ := IsProbe(p); ok {
				t.Fatalf("IsProbe(%q) = true, want false", p)
			}
		})
	}
}
