# `kck codex open`

Thin CLI over `agent-pro` `codex/sessions.RunOpen`: focus hosting tab or
`codex resume <id>` (agent-run prefer when managed).

L2 only — injectable `CodexHome` + `CodexOpenOpts`. No live iTerm.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — harness via `run.MainWith` with `Args: ["codex","open",…]`.
- **`kck/run`** — dispatches `codex open`, prints kck-flavored help, prefixes errors.
- **`sessions.RunOpen`** — core open (agent-pro).

## Behaviors

- Help surfaces list open.
- One host → focus.
- Unknown id → Error.
- No host + dry-run + --no-agent-run → would resume with `codex resume`.

## Version

0.0.1

## Decision Tree

```text
codex/open/
├── help/
│   ├── root-lists-open/
│   ├── codex-usage/
│   └── open-usage/
├── focus-exactly-one/
├── unknown-session/
└── dry-run-resume/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-open/` | `kck -h` mentions `codex open`. |
| `help/codex-usage/` | `kck codex --help` lists `open`. |
| `help/open-usage/` | Documents resume + `--no-agent-run` + `--tab`. |
| `focus-exactly-one/` | focused window/tab; FocusITerm called. |
| `unknown-session/` | `Error: codex session not found`. |
| `dry-run-resume/` | Would open; command contains `resume`. |

## How to Run

```sh
doctest vet ./tests/codex/open
doctest test ./tests/codex/open
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/agent-pro/agent/codex/sessions"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"

	"kck/run"
)

type Request struct {
	Args             []string
	CodexHome        string
	TempDir          string
	ProjectDir       string
	SessionID        string
	Procs            []sessions.FocusProc
	OpenFiles        map[int][]string
	ITerm            []iterm2.SessionRef
	CurrentSessionID string
	ControllingTTY   string
	NoAgentRun       bool
}

type Response struct {
	Stdout   string
	Stderr   string
	ErrText  string
	ExitCode int
	Focused  []string
	Opened   []string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	fake := &sessions.OpenFake{
		FocusFake: sessions.FocusFake{
			Procs:     req.Procs,
			OpenFiles: req.OpenFiles,
			ITerm:     req.ITerm,
		},
		CurrentSessionID: req.CurrentSessionID,
		ControllingTTY:   req.ControllingTTY,
	}
	opts := fake.OpenOpts()
	opts.NoAgentRun = req.NoAgentRun
	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:          append([]string(nil), req.Args...),
		Stdout:        &stdout,
		Stderr:        &stderr,
		CodexHome:     req.CodexHome,
		CodexOpenOpts: opts,
	})
	resp := &Response{
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Focused: append([]string(nil), fake.Focused...),
		Opened:  append([]string(nil), fake.Opened...),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
