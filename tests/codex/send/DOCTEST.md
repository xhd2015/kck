# `kck codex send`

Thin CLI over `agent-pro` `codex/sessions.RunSend`: type into hosting pane
(agent-run prefer for managed `--session-id`).

L2 only — injectable `CodexHome` + `CodexSendOpts`. No live iTerm.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — harness via `run.MainWith` with `Args: ["codex","send",…]`.
- **`kck/run`** — dispatches `codex send`, prints kck-flavored help.
- **`sessions.RunSend`** — core send (agent-pro).

## Behaviors

- Help surfaces list send.
- Hosting tab + `--no-agent-run` → SendText + `sent to session`.
- Missing payload → Error.
- No host → Error.

## Version

0.0.1

## Decision Tree

```text
codex/send/
├── help/
│   ├── root-lists-send/
│   ├── codex-usage/
│   └── send-usage/
├── iterm-send/
├── missing-payload/
├── no-host/
└── cron/
    └── invalid-expr/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-send/` | `kck -h` mentions `codex send`. |
| `help/codex-usage/` | `kck codex --help` lists `send`. |
| `help/send-usage/` | Documents `--session-id` / `--open` / keys / `--cron`. |
| `iterm-send/` | `sent to session`; type then two Enter writes (Codex submit). |
| `missing-payload/` | Error about missing text/key. |
| `no-host/` | Error: no hosting iTerm tab. |
| `cron/invalid-expr/` | `Error: invalid --cron:`; no SendText. |

## How to Run

```sh
doctest vet ./tests/codex/send
doctest test ./tests/codex/send
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
	Stdout    string
	Stderr    string
	ErrText   string
	ExitCode  int
	SendCalls []sessions.SendCall
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	fake := &sessions.SendFake{
		FocusFake: sessions.FocusFake{
			Procs:     req.Procs,
			OpenFiles: req.OpenFiles,
			ITerm:     req.ITerm,
		},
		CurrentSessionID: req.CurrentSessionID,
		ControllingTTY:   req.ControllingTTY,
	}
	opts := fake.SendOpts()
	opts.NoAgentRun = req.NoAgentRun
	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:          append([]string(nil), req.Args...),
		Stdout:        &stdout,
		Stderr:        &stderr,
		CodexHome:     req.CodexHome,
		CodexSendOpts: opts,
	})
	resp := &Response{
		Stdout:    stdout.String(),
		Stderr:    stderr.String(),
		SendCalls: append([]sessions.SendCall(nil), fake.SendCalls...),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
