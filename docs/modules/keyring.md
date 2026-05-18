# keyring

Store, retrieve, and delete secrets in the system credential store via per-platform backends.

## API

| Symbol | Kind | Summary |
| --- | --- | --- |
| `Set(service, account, secret string) error` | func | Stores or updates a secret. |
| `Get(service, account string) (string, error)` | func | Returns the stored secret, or `ErrNotFound`. |
| `Delete(service, account string) error` | func | Removes a stored secret. |
| `ErrEmptyService`, `ErrEmptyAccount` | error | Input validation failures. |
| `ErrNotFound` | error | No matching entry exists. |
| `ErrKeyringUnavailable` | error | No supported backend reachable. |

## Platform behavior

| Platform | Backend |
| --- | --- |
| macOS | `security add/find/delete-generic-password` (Keychain). |
| Linux | `secret-tool store/lookup/clear` (libsecret, Secret Service). |
| Windows | `cmdkey` for write/delete; PowerShell `CredentialManager` module for read. |

Linux requires `libsecret-tools` installed (`apt install libsecret-tools`). A running Secret Service provider (GNOME Keyring, KWallet via gnome-keyring-bridge, etc.) is also required. Windows read returns `ErrKeyringUnavailable` if the `CredentialManager` PowerShell module is not installed (`Install-Module CredentialManager`).

## Examples

```go
import "github.com/akira-io/onyx/keyring"

if err := keyring.Set("hyperion", "github_pat", token); err != nil {
    return err
}

stored, err := keyring.Get("hyperion", "github_pat")
if errors.Is(err, keyring.ErrNotFound) {
    // first-run prompt path
}
_ = keyring.Delete("hyperion", "github_pat")
```

## Errors

- `ErrEmptyService`, `ErrEmptyAccount` — input validation.
- `ErrNotFound` — entry does not exist.
- `ErrKeyringUnavailable` — backend missing (Linux without libsecret, Windows without CredentialManager module).
- Wrapped exec errors when a backend is found but the call itself fails (locked keychain, user cancelled prompt, etc.).

## Dependencies

- `osinfo` for platform detection.

## Related modules

- `clipboard` — clipboard read/write.
- `files` — opening files with the default app.
