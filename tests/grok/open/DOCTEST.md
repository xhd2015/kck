# `kck grok open`

Thin CLI over `agent-pro` `sessions.Open`: focus the hosting iTerm tab when
present, otherwise resume `grok --resume <session-id>` in a new window.
Supports `--tab` / `--tab-index` via shared `ResolveSessionSource` /
`ResolveFromTab`.

L2 only — injectable `GrokHome` + `GrokOpenOpts` (OpenFake). No live iTerm.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — harness via `run.MainWith` with `Args: ["grok","open",…]`.
- **`kck/run`** — dispatches `grok open`, prints kck-flavored help, prefixes
  errors with `Error:`.
- **`sessions.RunOpen`** — core focus-or-resume (agent-pro).

## Behaviors

- Root help mentions `grok open`.
- `kck grok --help` / `kck grok open --help` document tab source + resume.
- Unknown session → `Error: grok session not found: …`.
- Missing session id → usage error with `Error:`.
- One hosting tab → `focused: window W, tab T`.
- No live host → `opened: new window; resuming <id>` via injectable opener.
- Agent-run-managed exited → `opened: new window; agent-run resume <ar-id>`.
- `--tab N` → focus resolved tab (no resume).

## Version

0.0.1

## Decision Tree

```text
grok/open/
├── help/
│   ├── root-lists-grok/
│   ├── grok-usage/
│   └── open-usage/
├── missing-session-id/
├── unknown-session/
├── focus-exactly-one/
├── resume-no-live/
├── agent-run-resume/
└── tab-focus/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-grok/` | `kck -h` mentions `grok open`. |
| `help/grok-usage/` | `kck grok --help` lists `open`. |
| `help/open-usage/` | `kck grok open --help` documents `--tab` + resume + `--index`. |
| `missing-session-id/` | `Error:` usage; no focus/open. |
| `unknown-session/` | `Error: grok session not found`. |
| `focus-exactly-one/` | Focuses; opener not called. |
| `resume-no-live/` | Opens resume window; no focus. |
| `agent-run-resume/` | Prefer agent-run resume window when managed. |
| `tab-focus/` | `--tab 2` focuses resolved tab. |

## How to Run

```sh
doctest vet ./tests/grok/open
doctest test ./tests/grok/open
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
	NoAgentRun       bool
	AgentRunByID     map[string]*sessions.AgentRunOpenResult
	AgentRunErr      error
}

type Response struct {
	Stdout        string
	Stderr        string
	ErrText       string
	ExitCode      int
	Focused       []string
	Opened        []string
	AgentRunCalls []string
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
		AgentRunByID:     req.AgentRunByID,
		AgentRunErr:      req.AgentRunErr,
	}
	var stdout, stderr bytes.Buffer
	opts := fake.OpenOpts()
	opts.Stderr = &stderr
	opts.NoAgentRun = req.NoAgentRun
	err := run.MainWith(run.Options{
		Args:         append([]string(nil), req.Args...),
		Stdout:       &stdout,
		Stderr:       &stderr,
		GrokHome:     req.GrokHome,
		GrokOpenOpts: opts,
	})
	resp := &Response{
		Stdout:        stdout.String(),
		Stderr:        stderr.String(),
		Focused:       append([]string(nil), fake.Focused...),
		Opened:        append([]string(nil), fake.Opened...),
		AgentRunCalls: append([]string(nil), fake.AgentRunCalls...),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
