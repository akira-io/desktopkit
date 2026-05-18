package keyring

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/akira-io/onyx/osinfo"
)

var (
	ErrKeyringUnavailable = errors.New("keyring: no supported backend available")
	ErrNotFound           = errors.New("keyring: secret not found")
	ErrEmptyService       = errors.New("keyring: service must not be empty")
	ErrEmptyAccount       = errors.New("keyring: account must not be empty")
)

func Set(service, account, secret string) error {
	if err := validate(service, account); err != nil {
		return err
	}
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return setDarwin(service, account, secret)
	case platform.IsLinux():
		return setLinux(service, account, secret)
	case platform.IsWindows():
		return setWindows(service, account, secret)
	}
	return ErrKeyringUnavailable
}

func Get(service, account string) (string, error) {
	if err := validate(service, account); err != nil {
		return "", err
	}
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return getDarwin(service, account)
	case platform.IsLinux():
		return getLinux(service, account)
	case platform.IsWindows():
		return getWindows(service, account)
	}
	return "", ErrKeyringUnavailable
}

func Delete(service, account string) error {
	if err := validate(service, account); err != nil {
		return err
	}
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return deleteDarwin(service, account)
	case platform.IsLinux():
		return deleteLinux(service, account)
	case platform.IsWindows():
		return deleteWindows(service, account)
	}
	return ErrKeyringUnavailable
}

func validate(service, account string) error {
	if strings.TrimSpace(service) == "" {
		return ErrEmptyService
	}
	if strings.TrimSpace(account) == "" {
		return ErrEmptyAccount
	}
	return nil
}

func setDarwin(service, account, secret string) error {
	cmd := exec.Command("security", "add-generic-password", "-U", "-s", service, "-a", account, "-w", secret)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("keyring set via security: %w", err)
	}
	return nil
}

func getDarwin(service, account string) (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-s", service, "-a", account, "-w")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 44 {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("keyring get via security: %w", err)
	}
	return strings.TrimRight(stdout.String(), "\n"), nil
}

func deleteDarwin(service, account string) error {
	cmd := exec.Command("security", "delete-generic-password", "-s", service, "-a", account)
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 44 {
			return ErrNotFound
		}
		return fmt.Errorf("keyring delete via security: %w", err)
	}
	return nil
}

func setLinux(service, account, secret string) error {
	cmd := exec.Command("secret-tool", "store", "--label="+service, "service", service, "account", account)
	cmd.Stdin = strings.NewReader(secret)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("keyring set via secret-tool: %w", err)
	}
	return nil
}

func getLinux(service, account string) (string, error) {
	cmd := exec.Command("secret-tool", "lookup", "service", service, "account", account)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("keyring get via secret-tool: %w", err)
	}
	if len(out) == 0 {
		return "", ErrNotFound
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func deleteLinux(service, account string) error {
	cmd := exec.Command("secret-tool", "clear", "service", service, "account", account)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("keyring delete via secret-tool: %w", err)
	}
	return nil
}

func windowsTarget(service, account string) string {
	return service + ":" + account
}

func setWindows(service, account, secret string) error {
	target := windowsTarget(service, account)
	cmd := exec.Command("cmdkey", "/generic:"+target, "/user:"+account, "/pass:"+secret)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("keyring set via cmdkey: %w", err)
	}
	return nil
}

func getWindows(service, account string) (string, error) {
	target := windowsTarget(service, account)
	script := fmt.Sprintf(
		`if (-not (Get-Module -ListAvailable -Name CredentialManager)) { exit 2 }; Import-Module CredentialManager; $c = Get-StoredCredential -Target '%s'; if ($null -eq $c) { exit 3 }; [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($c.Password))`,
		strings.ReplaceAll(target, `'`, `''`),
	)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			switch exitErr.ExitCode() {
			case 2:
				return "", ErrKeyringUnavailable
			case 3:
				return "", ErrNotFound
			}
		}
		return "", fmt.Errorf("keyring get via powershell: %w", err)
	}
	return strings.TrimRight(string(out), "\r\n"), nil
}

func deleteWindows(service, account string) error {
	target := windowsTarget(service, account)
	cmd := exec.Command("cmdkey", "/delete:"+target)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("keyring delete via cmdkey: %w", err)
	}
	return nil
}
