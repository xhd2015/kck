# `kck codex focus`

Thin CLI over `codex/sessions.RunFocus`: focus a live hosting iTerm tab only.
Never resumes or creates a window (unlike `open`).

L2 — injectable `CodexHome` + `CodexFocusOpts` (FocusFake).

## Version

0.0.1

## Decision Tree

```text
codex/focus/
├── help/focus-usage/
├── focus-exactly-one/
├── no-live-fails/
└── unknown-session/
```

## How to Run

```sh
doctest vet ./tests/codex/focus
doctest test ./tests/codex/focus
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
	Args       []string
	CodexHome  string
	TempDir    string
	ProjectDir string
	SessionID  string
	Procs      []sessions.FocusProc
	OpenFiles  map[int][]string
	ITerm      []iterm2.SessionRef
}

type Response struct {
	Stdout   string
	Stderr   string
	ErrText  string
	ExitCode int
	Focused  []string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	fake := &sessions.FocusFake{
		Procs:     req.Procs,
		OpenFiles: req.OpenFiles,
		ITerm:     req.ITerm,
	}
	var stdout, stderr bytes.Buffer
	opts := fake.Opts()
	err := run.MainWith(run.Options{
		Args:           append([]string(nil), req.Args...),
		Stdout:         &stdout,
		Stderr:         &stderr,
		CodexHome:      req.CodexHome,
		CodexFocusOpts: opts,
	})
	resp := &Response{
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Focused: append([]string(nil), fake.Focused...),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
