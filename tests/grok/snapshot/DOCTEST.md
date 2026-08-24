# `kck grok snapshot`

Thin CLI over `agent-pro` `sessions.Snapshot`: capture visible pane text for a
live Grok host (agent-run TTY when managed, else iTerm Contents). Supports
`--tab` / `--tab-index` via shared `ResolveSessionSource` / `ResolveFromTab`.
No resume when no host.

L2 only — injectable `GrokHome` + `GrokSnapshotOpts` (SnapshotFake). No live iTerm.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — harness via `run.MainWith` with `Args: ["grok","snapshot",…]`.
- **`kck/run`** — dispatches `grok snapshot`, prints kck-flavored help, prefixes
  errors with `Error:`.
- **`sessions.RunSnapshot`** — core capture (agent-pro).

## Behaviors

- Root help mentions `grok snapshot`.
- `kck grok --help` / `kck grok snapshot --help` document tab source + Contents.
- Unknown session → `Error: grok session not found: …`.
- Missing session id → usage error with `Error:`.
- One hosting tab → pane text on stdout.
- No live host → `Error: no hosting iTerm tab …`.
- `--tab N` → capture resolved tab pane.

## Version

0.0.1

## Decision Tree

```text
grok/snapshot/
├── help/
│   ├── root-lists-snapshot/
│   ├── grok-usage/
│   └── snapshot-usage/
├── missing-session-id/
├── unknown-session/
├── capture-exactly-one/
├── no-live-fails/
└── tab-capture/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-snapshot/` | `kck -h` mentions `grok snapshot`. |
| `help/grok-usage/` | `kck grok --help` lists `snapshot`. |
| `help/snapshot-usage/` | `kck grok snapshot --help` documents `--tab`, `--index`, `--iterm`, agent-run. |
| `missing-session-id/` | `Error:` usage; no Contents. |
| `unknown-session/` | `Error: grok session not found`. |
| `capture-exactly-one/` | Pane text on stdout; Contents called. |
| `no-live-fails/` | Hard error; Contents not called. |
| `tab-capture/` | `--tab 2` captures resolved tab pane. |

## How to Run

```sh
doctest vet ./tests/grok/snapshot
doctest test ./tests/grok/snapshot
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"

	"kck/run"
)

type Request struct {
	Args             []string
	GrokHome         string
	TempDir          string
	ProjectDir       string
	SessionID        string
	Procs            []sessions.FocusProc
	OpenFiles        map[int][]string
	ITerm            []iterm2.SessionRef
	CurrentSessionID string
	ControllingTTY   string
	ContentsByID     map[string]iterm2.ContentsResult
}

type Response struct {
	Stdout        string
	Stderr        string
	ErrText       string
	ExitCode      int
	ContentsCalls []string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	fake := &sessions.SnapshotFake{
		FocusFake: sessions.FocusFake{
			Procs:     req.Procs,
			OpenFiles: req.OpenFiles,
			ITerm:     req.ITerm,
		},
		CurrentSessionID: req.CurrentSessionID,
		ControllingTTY:   req.ControllingTTY,
		ContentsByID:     req.ContentsByID,
	}
	var stdout, stderr bytes.Buffer
	opts := fake.SnapshotOpts()
	err := run.MainWith(run.Options{
		Args:             append([]string(nil), req.Args...),
		Stdout:           &stdout,
		Stderr:           &stderr,
		GrokHome:         req.GrokHome,
		GrokSnapshotOpts: opts,
	})
	resp := &Response{
		Stdout:        stdout.String(),
		Stderr:        stderr.String(),
		ContentsCalls: append([]string(nil), fake.ContentsCalls...),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
