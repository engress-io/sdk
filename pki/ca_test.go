package pki

import (
	"crypto/tls"
	"crypto/x509"
	"io"
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

func TestServerTLSConfigIncludesSAN(t *testing.T) {
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
	serverTLS, err := ca.ServerTLSConfig(tls.NoClientCert)
	if err != nil {
		t.Fatal(err)
	}
	ln, err := tls.Listen("tcp", "127.0.0.1:0", serverTLS)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	accepted := make(chan struct{}, 1)
	go func() {
		c, err := ln.Accept()
		if err == nil {
			accepted <- struct{}{}
			_, _ = io.Copy(io.Discard, c)
			_ = c.Close()
		}
	}()
	pool := x509.NewCertPool()
	_ = pool.AppendCertsFromPEM(ca.CACertPEM())
	client := &tls.Config{RootCAs: pool, ServerName: TunnelServerCN, MinVersion: tls.VersionTLS13}
	conn, err := tls.Dial("tcp", ln.Addr().String(), client)
	if err != nil {
		t.Fatal(err)
	}
	_ = conn.Close()
	select {
	case <-accepted:
	default:
		t.Fatal("server did not accept")
	}
}
