package notify

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/akira-io/onyx/osinfo"
)

var (
	ErrNotifyUnavailable = errors.New("notify: no supported backend available")
	ErrEmptyTitle        = errors.New("notify: title must not be empty")
)

func Show(title, body string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return showDarwin(title, body)
	case platform.IsLinux():
		return showLinux(title, body)
	case platform.IsWindows():
		return showWindows(title, body)
	}
	return ErrNotifyUnavailable
}

func showDarwin(title, body string) error {
	script := fmt.Sprintf(
		`display notification %s with title %s`,
		quoteAppleScript(body), quoteAppleScript(title),
	)
	if err := exec.Command("osascript", "-e", script).Run(); err != nil {
		return fmt.Errorf("notify via osascript: %w", err)
	}
	return nil
}

func showLinux(title, body string) error {
	if err := exec.Command("notify-send", title, body).Run(); err != nil {
		return fmt.Errorf("notify via notify-send: %w", err)
	}
	return nil
}

func showWindows(title, body string) error {
	burnt := fmt.Sprintf(
		`if (Get-Module -ListAvailable -Name BurntToast) { Import-Module BurntToast; New-BurntToastNotification -Text %s, %s; exit 0 } else { exit 1 }`,
		quotePowerShell(title), quotePowerShell(body),
	)
	if err := exec.Command("powershell", "-NoProfile", "-Command", burnt).Run(); err == nil {
		return nil
	}
	if err := exec.Command("msg", "*", fmt.Sprintf("%s\n%s", title, body)).Run(); err == nil {
		return nil
	}
	return ErrNotifyUnavailable
}

func quoteAppleScript(s string) string {
	return `"` + strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `"`, `\"`) + `"`
}

func quotePowerShell(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `''`) + `'`
}
