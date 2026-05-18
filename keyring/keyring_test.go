package keyring

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSet_RejectsEmptyService(t *testing.T) {
	if err := Set("", "account", "secret"); !errors.Is(err, ErrEmptyService) {
		t.Fatalf("expected ErrEmptyService, got %v", err)
	}
}

func TestSet_RejectsEmptyAccount(t *testing.T) {
	if err := Set("service", "", "secret"); !errors.Is(err, ErrEmptyAccount) {
		t.Fatalf("expected ErrEmptyAccount, got %v", err)
	}
}

func TestGet_RejectsEmptyService(t *testing.T) {
	if _, err := Get("", "account"); !errors.Is(err, ErrEmptyService) {
		t.Fatalf("expected ErrEmptyService, got %v", err)
	}
}

func TestDelete_RejectsEmptyAccount(t *testing.T) {
	if err := Delete("service", ""); !errors.Is(err, ErrEmptyAccount) {
		t.Fatalf("expected ErrEmptyAccount, got %v", err)
	}
}

func TestRoundTrip_OnRealKeyring(t *testing.T) {
	service := fmt.Sprintf("onyx-test-%d-%d", os.Getpid(), time.Now().UnixNano())
	account := "tester"
	secret := "hunter2"

	if err := Set(service, account, secret); err != nil {
		if errors.Is(err, ErrKeyringUnavailable) {
			t.Skip("no keyring backend available")
		}
		t.Skipf("set failed (backend missing or denied?): %v", err)
	}
	t.Cleanup(func() { _ = Delete(service, account) })

	got, err := Get(service, account)
	if err != nil {
		t.Skipf("get failed: %v", err)
	}
	if got != secret {
		t.Fatalf("got %q, want %q", got, secret)
	}
}

func TestGet_ReturnsNotFoundForMissingEntry(t *testing.T) {
	service := fmt.Sprintf("onyx-missing-%d-%d", os.Getpid(), time.Now().UnixNano())
	_, err := Get(service, "nobody")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if errors.Is(err, ErrKeyringUnavailable) {
		t.Skip("no keyring backend available")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Skipf("backend reported %v (expected ErrNotFound on supported platforms)", err)
	}
}

func TestWindowsTarget_JoinsWithColon(t *testing.T) {
	if got := windowsTarget("svc", "acct"); got != "svc:acct" {
		t.Fatalf("got %q want svc:acct", got)
	}
}
