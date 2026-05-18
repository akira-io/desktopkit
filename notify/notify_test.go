package notify

import (
	"errors"
	"strings"
	"testing"
)

func TestShow_RejectsEmptyTitle(t *testing.T) {
	if err := Show("", "body"); !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected ErrEmptyTitle, got %v", err)
	}
	if err := Show("   ", "body"); !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected ErrEmptyTitle for whitespace title, got %v", err)
	}
}

func TestQuoteAppleScript_EscapesQuotesAndBackslashes(t *testing.T) {
	got := quoteAppleScript(`he said "hi" \\path`)
	want := `"he said \"hi\" \\\\path"`
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestQuotePowerShell_DoublesSingleQuotes(t *testing.T) {
	got := quotePowerShell(`it's fine`)
	want := `'it''s fine'`
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestShow_DoesNotPanicOnBestEffort(t *testing.T) {
	err := Show("onyx-test", "body")
	if err != nil && !errors.Is(err, ErrNotifyUnavailable) && !strings.Contains(err.Error(), "notify") {
		t.Fatalf("unexpected error shape: %v", err)
	}
}
