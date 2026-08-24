# `kck grok list`

Thin CLI over `agent-pro` `sessions.RunListLive`: list Grok session ids
hosted in iTerm tabs. L2 — injectable `GrokListLiveOpts`.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — `run.MainWith` with `Args: ["grok","list",…]`.
- **`kck/run`** — dispatches, prints kck-flavored help, prefixes `Error:`.
- **`sessions.RunListLive`** — core list-live.

## Behaviors

- Root / grok help mention `list`.
- `kck grok list --help` documents `--json` / `--limit`.
- One host → table row with sid.
- Empty → `0 sessions`.
- `--json` machine envelope.

## Version

0.0.1

## Decision Tree

```text
grok/list/
├── help/
│   ├── root-lists-list/
│   ├── grok-usage/
│   └── list-usage/
├── one-host/
├── empty/
└── json/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-list/` | `kck -h` mentions `grok list`. |
| `help/grok-usage/` | `kck grok --help` lists `list`. |
| `help/list-usage/` | `kck grok list --help` documents flags. |
| `one-host/` | Table contains sid + ITERM. |
| `empty/` | `0 sessions`. |
| `json/` | JSON session_id. |

## How to Run

```sh
doctest vet ./tests/grok/list
doctest test ./tests/grok/list
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
	Args         []string
	GrokHome     string
	Procs        []sessions.FocusProc
	OpenFiles    map[int][]string
	ITerm        []iterm2.SessionRef
	PaneByTTY    map[string]sessions.LivePaneInfo
	CwdBySession map[string]string
}

type Response struct {
	Stdout   string
	Stderr   string
	ErrText  string
	ExitCode int
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	fake := &sessions.ListLiveFake{
		FocusFake: sessions.FocusFake{
			Procs:     req.Procs,
			OpenFiles: req.OpenFiles,
			ITerm:     req.ITerm,
		},
		PaneByTTY:    req.PaneByTTY,
		CwdBySession: req.CwdBySession,
	}
	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:             append([]string(nil), req.Args...),
		Stdout:           &stdout,
		Stderr:           &stderr,
		GrokHome:         req.GrokHome,
		GrokListLiveOpts: fake.ListLiveOpts(),
	})
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String()}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
