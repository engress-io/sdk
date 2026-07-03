package subdomain

import (
	"context"
	"regexp"
	"strings"
	"testing"
)

type mockStore struct {
	taken map[string]bool
}

func (m *mockStore) SubdomainExists(_ context.Context, subdomain string) (bool, error) {
	if m.taken == nil {
		return false, nil
	}
	return m.taken[subdomain], nil
}

type collisionStore struct {
	attempts int
}

func (c *collisionStore) SubdomainExists(_ context.Context, _ string) (bool, error) {
	c.attempts++
	return c.attempts == 1, nil
}

func TestSemanticAllocate_Format(t *testing.T) {
	st := &mockStore{}
	sub, err := SemanticAllocate(context.Background(), st, "https-")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(sub, "https-") {
		t.Fatalf("got %q", sub)
	}
	suffix := strings.TrimPrefix(sub, "https-")
	if len(suffix) != 8 {
		t.Fatalf("suffix len %d", len(suffix))
	}
	if !regexp.MustCompile(`^[a-z0-9]{8}$`).MatchString(suffix) {
		t.Fatalf("invalid suffix %q", suffix)
	}
}

func TestSemanticAllocate_EmptyPrefix(t *testing.T) {
	st := &mockStore{}
	sub, err := SemanticAllocate(context.Background(), st, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(sub, "tunnel-") {
		t.Fatalf("got %q", sub)
	}
	suffix := strings.TrimPrefix(sub, "tunnel-")
	if len(suffix) != 8 {
		t.Fatalf("suffix len %d", len(suffix))
	}
}

func TestSemanticAllocate_Collision(t *testing.T) {
	st := &collisionStore{}
	sub, err := SemanticAllocate(context.Background(), st, "https-")
	if err != nil {
		t.Fatal(err)
	}
	if st.attempts < 2 {
		t.Fatalf("expected at least 2 attempts, got %d", st.attempts)
	}
	if !strings.HasPrefix(sub, "https-") {
		t.Fatalf("got %q", sub)
	}
}

func TestRenameWithPrefix(t *testing.T) {
	st := &mockStore{}
	sub, err := RenameWithPrefix(context.Background(), st, "http-")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(sub, "http-") {
		t.Fatalf("got %q", sub)
	}
}