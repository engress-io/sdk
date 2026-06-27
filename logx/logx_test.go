package logx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestInitJSONFormat(t *testing.T) {
	t.Setenv("FLUX_LOG_FORMAT", "json")
	var buf bytes.Buffer
	InitWriter(&buf)
	slog.Info("hello", "key", "value")
	line := strings.TrimSpace(buf.String())
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("not json: %q err=%v", line, err)
	}
	if m["msg"] != "hello" {
		t.Fatalf("msg = %v", m["msg"])
	}
}

func TestInitTextFormat(t *testing.T) {
	t.Setenv("FLUX_LOG_FORMAT", "text")
	var buf bytes.Buffer
	InitWriter(&buf)
	slog.Info("plain")
	if !strings.Contains(buf.String(), "plain") {
		t.Fatalf("output = %q", buf.String())
	}
}

func TestInitDefaultsToText(t *testing.T) {
	t.Setenv("FLUX_LOG_FORMAT", "")
	var buf bytes.Buffer
	InitWriter(&buf)
	slog.Warn("warn")
	if buf.Len() == 0 {
		t.Fatal("expected output")
	}
}

func TestQuietDiscardsOutput(t *testing.T) {
	t.Setenv("FLUX_LOG_FORMAT", "text")
	var buf bytes.Buffer
	InitWriter(&buf)
	before := slog.Default()
	Quiet()
	slog.Info("hidden")
	if buf.Len() != 0 {
		t.Fatalf("expected quiet to discard, got %q", buf.String())
	}
	Unquiet()
	if slog.Default() != before {
		t.Fatal("Unquiet did not restore logger")
	}
	slog.Info("visible")
	if buf.Len() == 0 {
		t.Fatal("expected output after Unquiet")
	}
}
