package engine

import (
	"errors"
	"fmt"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestEngine_SetAndGet(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)
	eng := NewEngine(logger)

	if err := eng.Set("k", "v"); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	if err := eng.Set("k", "v2"); err != nil {
		t.Fatalf("Set overwrite returned error: %v", err)
	}

	got, err := eng.Get("k")
	if err != nil {
		t.Fatalf("Get existing returned error: %v", err)
	}
	if got != "v2" {
		t.Fatalf("Get mismatch: want %q, got %q", "v2", got)
	}

	if recorded.Len() != 0 {
		t.Fatalf("unexpected logs on success: %d", recorded.Len())
	}
}

func TestEngine_GetAndDelete_NotFound(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)
	eng := NewEngine(logger)

	_, err := eng.Get("missing")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("want ErrKeyNotFound, got %v", err)
	}
	if recorded.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", recorded.Len())
	}
	entry := recorded.All()[0]
	if entry.Message != "Get: entry not found" {
		t.Fatalf("unexpected log message: %q", entry.Message)
	}
	if got := fmt.Sprint(entry.ContextMap()["key"]); got != "missing" {
		t.Fatalf("logged key mismatch: want %q, got %q", "missing", got)
	}

	recorded.TakeAll()
	err = eng.Delete("absent")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("want ErrKeyNotFound, got %v", err)
	}
	if recorded.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", recorded.Len())
	}
	entry = recorded.All()[0]
	if entry.Message != "Delete: entry not found" {
		t.Fatalf("unexpected log message: %q", entry.Message)
	}
	if got := fmt.Sprint(entry.ContextMap()["key"]); got != "absent" {
		t.Fatalf("logged key mismatch: want %q, got %q", "absent", got)
	}
}

func TestEngine_DeleteExisting(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)
	eng := NewEngine(logger)

	if err := eng.Set("k", "v"); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}
	if err := eng.Delete("k"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	_, err := eng.Get("k")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound after delete, got %v", err)
	}

	found := false
	for _, e := range recorded.All() {
		if e.Message == "Delete: entry not found" {
			found = true
		}
	}
	if found {
		t.Fatalf("should not log on successful delete")
	}
}
