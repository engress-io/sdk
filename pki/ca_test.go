package pki

import (
	"testing"
	"time"
)

func TestMintAndParseClientCert(t *testing.T) {
	dir := t.TempDir()
	certPath := dir + "/ca.crt"
	keyPath := dir + "/ca.key"
	if err := EnsureDevCA(certPath, keyPath); err != nil {
		t.Fatal(err)
	}
	ca, err := LoadCA(certPath, keyPath)
	if err != nil {
		t.Fatal(err)
	}
	certPEM, _, expiresAt, err := ca.MintClientCert("ep-123", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if len(certPEM) == 0 || !expiresAt.After(time.Now()) {
		t.Fatal("expected valid cert")
	}
}
