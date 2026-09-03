# `kck grok new`

Open a new empty Grok session via agent-run. Default: new terminal + wait for
session id. `--here` / `--no-new-terminal`: current terminal, silent.

L2 only — injectable `GrokNewOpts`. No live iTerm.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — harness via `run.MainWith` with `Args: ["grok","new",…]`.
- **`kck/run`** — builds brainstorm-prefixed prompt, allocates session id,
  launches `agent-run run --open …`.
- **`agent-run`** — TTY session (`grok-tty`); `--new-terminal` ForceNew, or
  foreground `--here`.

## Behaviors

- Default / `--new-terminal`: launch with `--new-terminal`, wait for **provider**
  (Grok) session id via registry PID → open-files, print `opened:` + that id.
- `--here` / `--no-new-terminal`: foreground launch; **no kck stdout/stderr**.
- Prompt wrap: `/brainstorm <msg>`; runner `grok-tty`.
- Default stages (`--no-submit`); `--submit` omits it.
- No `-n` / `--new-window` / `--no-agent-run`.

## Version

0.0.1

## Decision Tree

```text
grok/new/
├── help/
│   ├── root-lists-new/
│   ├── grok-usage/
│   └── new-usage/
├── missing-message/
├── here-new-terminal-conflict/
├── dry-run/
├── agent-run-new-terminal/
├── agent-run-here/
├── no-new-terminal-alias/
└── agent-run-submit/
```

## How to Run

```sh
doctest vet ./tests/grok/new
doctest test ./tests/grok/new
```

```go
import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"

	"kck/run"
)

type Request struct {
	Args       []string
	TempDir    string
	ProjectDir string
	AgentHome  string
	UserHome   string
	FixedNow   time.Time

	LookPathMap map[string]string
}

type Response struct {
	Stdout       string
	Stderr       string
	ErrText      string
	ExitCode     int
	Foreground   []string
	NewTerminal  []string
	WaitedIDs    []string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	var foreground, newTerminal, waited []string
	look := req.LookPathMap
	if look == nil {
		look = map[string]string{
			"agent-run": "/usr/local/bin/agent-run",
		}
	}
	now := req.FixedNow
	if now.IsZero() {
		now = time.Date(2026, 9, 2, 12, 0, 0, 0, time.Local)
	}
	nopts := &run.NewOpts{
		AgentRunHome: req.AgentHome,
		UserHomeDir:  func() (string, error) { return req.UserHome, nil },
		Getwd:        func() (string, error) { return req.ProjectDir, nil },
		Abs: func(p string) (string, error) {
			if strings.HasPrefix(p, "/") {
				return p, nil
			}
			return req.ProjectDir + "/" + p, nil
		},
		Stat: os.Stat,
		Now:  func() time.Time { return now },
		LookPath: func(file string) (string, error) {
			if p, ok := look[file]; ok {
				return p, nil
			}
			return "", os.ErrNotExist
		},
		RunForeground: func(bin string, argv []string, dir string) error {
			foreground = append(foreground, dir+"|"+bin+" "+strings.Join(argv, " "))
			return nil
		},
		RunNewTerminal: func(bin string, argv []string, dir string) error {
			newTerminal = append(newTerminal, dir+"|"+bin+" "+strings.Join(argv, " "))
			return nil
		},
		WaitProviderSession: func(home, runner, agentRunSessionID string) (string, error) {
			_ = home
			_ = runner
			waited = append(waited, agentRunSessionID)
			return "01a064d2-70ec-7162-b36b-8a50ba323569", nil
		},
		SessionExists: func(home, sessionID string) bool {
			_ = home
			_ = sessionID
			return false
		},
	}

	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:        append([]string(nil), req.Args...),
		Stdout:      &stdout,
		Stderr:      &stderr,
		GrokNewOpts: nopts,
	})
	resp := &Response{
		Stdout:      stdout.String(),
		Stderr:      stderr.String(),
		Foreground:  append([]string(nil), foreground...),
		NewTerminal: append([]string(nil), newTerminal...),
		WaitedIDs:   append([]string(nil), waited...),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
