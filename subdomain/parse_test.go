package subdomain

import "testing"

func TestParseSubdomain_Valid(t *testing.T) {
	prefix, suffix, ok := ParseSubdomain("https-a1b2c3d4")
	if !ok {
		t.Fatal("expected ok")
	}
	if prefix != "https-" {
		t.Fatalf("prefix %q", prefix)
	}
	if suffix != "a1b2c3d4" {
		t.Fatalf("suffix %q", suffix)
	}
}

func TestParseSubdomain_Invalid(t *testing.T) {
	cases := []string{
		"",
		"short",
		"https-a1b2c3",
		"https-A1B2C3D4",
		"https-a1b2c3d!",
	}
	for _, label := range cases {
		_, _, ok := ParseSubdomain(label)
		if ok {
			t.Fatalf("expected invalid for %q", label)
		}
	}
}