package clipboard

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/akira-io/onyx/osinfo"
)

var ErrClipboardUnavailable = errors.New("clipboard: no supported backend available")

func Read() (string, error) {
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return runReader("pbpaste")
	case platform.IsWindows():
		return readWindows()
	case platform.IsLinux():
		for _, b := range linuxReaders() {
			out, err := runReader(b.cmd, b.args...)
			if err == nil {
				return out, nil
			}
		}
		return "", ErrClipboardUnavailable
	}
	return "", ErrClipboardUnavailable
}

func Write(text string) error {
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return runWriter(text, "pbcopy")
	case platform.IsWindows():
		return writeWindows(text)
	case platform.IsLinux():
		for _, b := range linuxWriters() {
			if err := runWriter(text, b.cmd, b.args...); err == nil {
				return nil
			}
		}
		return ErrClipboardUnavailable
	}
	return ErrClipboardUnavailable
}

func readWindows() (string, error) {
	script := `Add-Type -AssemblyName System.Windows.Forms; ` +
		`[System.Windows.Forms.Clipboard]::GetText()`
	cmd := exec.Command("powershell", "-NoProfile", "-STA", "-Command", script)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("clipboard read via powershell: %w", err)
	}
	return strings.TrimRight(string(out), "\r\n"), nil
}

func writeWindows(text string) error {
	script := `Add-Type -AssemblyName System.Windows.Forms; ` +
		`$v = $env:ONYX_CLIP_TEXT; ` +
		`if ([string]::IsNullOrEmpty($v)) { [System.Windows.Forms.Clipboard]::Clear() } ` +
		`else { [System.Windows.Forms.Clipboard]::SetText($v) }`
	cmd := exec.Command("powershell", "-NoProfile", "-STA", "-Command", script)
	cmd.Env = append(cmd.Environ(), "ONYX_CLIP_TEXT="+text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard write via powershell: %w", err)
	}
	return nil
}

type backend struct {
	cmd  string
	args []string
}

func linuxReaders() []backend {
	return []backend{
		{"wl-paste", []string{"--no-newline"}},
		{"xclip", []string{"-selection", "clipboard", "-o"}},
		{"xsel", []string{"--clipboard", "--output"}},
	}
}

func linuxWriters() []backend {
	return []backend{
		{"wl-copy", nil},
		{"xclip", []string{"-selection", "clipboard"}},
		{"xsel", []string{"--clipboard", "--input"}},
	}
}

func runReader(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", fmt.Errorf("clipboard read via %s: %w", name, err)
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func runWriter(text, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard write via %s: %w", name, err)
	}
	return nil
}
