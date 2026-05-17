package paths

import (
	"errors"
	"strings"
	"testing"
)

func TestFor_RememberApplicationName(t *testing.T) {
	app := For("Hyperion")
	if app.Name() != "Hyperion" {
		t.Fatalf("expected name %q, got %q", "Hyperion", app.Name())
	}
}

func TestConfig_RequiresApplicationName(t *testing.T) {
	app := For("")
	_, err := app.Config()
	if !errors.Is(err, ErrMissingApplicationName) {
		t.Fatalf("expected ErrMissingApplicationName, got %v", err)
	}
}

func TestData_RequiresApplicationName(t *testing.T) {
	app := For("")
	_, err := app.Data()
	if !errors.Is(err, ErrMissingApplicationName) {
		t.Fatalf("expected ErrMissingApplicationName, got %v", err)
	}
}

func TestCache_RequiresApplicationName(t *testing.T) {
	app := For("")
	_, err := app.Cache()
	if !errors.Is(err, ErrMissingApplicationName) {
		t.Fatalf("expected ErrMissingApplicationName, got %v", err)
	}
}

func TestLogs_RequiresApplicationName(t *testing.T) {
	app := For("")
	_, err := app.Logs()
	if !errors.Is(err, ErrMissingApplicationName) {
		t.Fatalf("expected ErrMissingApplicationName, got %v", err)
	}
}

func TestConfig_IncludesApplicationNameInPath(t *testing.T) {
	app := For("Hyperion")
	got, err := app.Config()
	if err != nil {
		t.Fatalf("Config: %v", err)
	}
	if !strings.Contains(got, "Hyperion") {
		t.Fatalf("expected path to contain application name, got %q", got)
	}
}

func TestData_IncludesApplicationNameInPath(t *testing.T) {
	app := For("Hyperion")
	got, err := app.Data()
	if err != nil {
		t.Fatalf("Data: %v", err)
	}
	if !strings.Contains(got, "Hyperion") {
		t.Fatalf("expected path to contain application name, got %q", got)
	}
}

func TestCache_IncludesApplicationNameInPath(t *testing.T) {
	app := For("Hyperion")
	got, err := app.Cache()
	if err != nil {
		t.Fatalf("Cache: %v", err)
	}
	if !strings.Contains(got, "Hyperion") {
		t.Fatalf("expected path to contain application name, got %q", got)
	}
}

func TestLogs_IncludesApplicationNameInPath(t *testing.T) {
	app := For("Hyperion")
	got, err := app.Logs()
	if err != nil {
		t.Fatalf("Logs: %v", err)
	}
	if !strings.Contains(got, "Hyperion") {
		t.Fatalf("expected path to contain application name, got %q", got)
	}
}
