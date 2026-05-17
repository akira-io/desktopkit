package files

import (
	"errors"
	"testing"
)

func TestOpenPath_RejectsEmptyPath(t *testing.T) {
	if err := OpenPath(""); !errors.Is(err, ErrPathRequired) {
		t.Fatalf("expected ErrPathRequired, got %v", err)
	}
}

func TestOpenURL_RejectsEmptyURL(t *testing.T) {
	if err := OpenURL(""); !errors.Is(err, ErrPathRequired) {
		t.Fatalf("expected ErrPathRequired, got %v", err)
	}
}

func TestRevealInFileManager_RejectsEmptyPath(t *testing.T) {
	if err := RevealInFileManager(""); !errors.Is(err, ErrPathRequired) {
		t.Fatalf("expected ErrPathRequired, got %v", err)
	}
}
