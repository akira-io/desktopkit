# notify

Show desktop notifications. Per-platform backends, no third-party deps.

## API

| Symbol | Kind | Summary |
| --- | --- | --- |
| `Show(title, body string) error` | func | Displays a notification. Title must be non-empty. |
| `ErrEmptyTitle` | error | Returned when title is empty or whitespace-only. |
| `ErrNotifyUnavailable` | error | Returned when no supported backend is reachable. |

## Platform behavior

| Platform | Backend |
| --- | --- |
| macOS | `osascript -e 'display notification ...'`. Always present. |
| Linux | `notify-send` (from `libnotify-bin`). |
| Windows | PowerShell `BurntToast` module (modern toast). Falls back to `msg.exe`. |

Linux requires `libnotify-bin` installed (`apt install libnotify-bin`). Windows toasts require the `BurntToast` PowerShell module (`Install-Module BurntToast`); without it the call falls back to `msg.exe`, which delivers a basic message box.

## Examples

```go
import "github.com/akira-io/onyx/notify"

if err := notify.Show("Hyperion", "Export finished."); err != nil {
    log.Println(err)
}
```

## Errors

- `ErrEmptyTitle` — title cannot be empty.
- `ErrNotifyUnavailable` — every backend failed (Windows when neither BurntToast nor `msg.exe` is reachable; unrecognized platform).
- Wrapped exec errors when a backend was found but the invocation itself failed.

## Dependencies

- `osinfo` for platform detection.

## Related modules

- `files` — opening files with the default app.
- `clipboard` — clipboard read/write.
