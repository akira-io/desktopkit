package shell

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestResolve_FailsWhenNothingMatches(t *testing.T) {
	candidates := NewCandidates().
		WithName("definitely-not-a-real-binary-xyz").
		WithCandidate("/definitely/not/a/path/binary")

	_, err := candidates.Resolve()
	if !errors.Is(err, ErrBinaryNotFound) {
		t.Fatalf("expected ErrBinaryNotFound, got %v", err)
	}
}

func TestResolve_FindsExplicitCandidateFile(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "fakebin")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	candidates := NewCandidates().
		WithName("definitely-not-a-real-binary-xyz").
		WithCandidate(binary)

	resolved, err := candidates.Resolve()
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if resolved.AbsolutePath() != binary {
		t.Fatalf("expected %q, got %q", binary, resolved.AbsolutePath())
	}
	if resolved.Source() != SourceCandidate {
		t.Fatalf("expected SourceCandidate, got %s", resolved.Source())
	}
}

func TestResolve_PrefersPATHOverCandidates(t *testing.T) {
	candidates := NewCandidates().
		WithName("sh").
		WithCandidate("/definitely/not/a/path/binary")

	resolved, err := candidates.Resolve()
	if err != nil {
		t.Skipf("sh not available on this system: %v", err)
	}
	if resolved.Source() != SourcePath {
		t.Fatalf("expected SourcePath, got %s", resolved.Source())
	}
}

func TestCandidates_IgnoresEmptyInputs(t *testing.T) {
	candidates := NewCandidates().
		WithName("").
		WithCandidate("")

	_, err := candidates.Resolve()
	if !errors.Is(err, ErrBinaryNotFound) {
		t.Fatalf("expected ErrBinaryNotFound, got %v", err)
	}
}
