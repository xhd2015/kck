# `kck codex status`

Thin CLI over `sessions.Status` + formatters. PID liveness plus rollout path
(`pathfmt.TildeHome` in text; absolute in JSON). File always no (no
`active_sessions.json` for Codex).

L2 only — injectable `CodexHome` + `CodexLiveOpts`. No live ps/lsof.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — `run.MainWith` with `Args: ["codex","status",…]`.
- **`kck/run`** — parses flags, `Error:` prefix, kck-flavored help.
- **`sessions.Status`** — Find + optional live PIDs (FileActive always false).

## Behaviors

- Root / codex / status help document status + path.
- Missing / unknown id → `Error:`.
- Running / inactive / `--no-pid` states.
- Text includes `Path:` with `rollout-`.
- `--json` includes absolute `path`, `file_active: false`, no ANSI.

## Version

0.0.1

## Decision Tree

```text
codex/status/
├── help/
│   ├── root-lists-status/
│   ├── codex-usage/
│   └── status-usage/
├── missing-session-id/
├── unknown-session/
├── running/
├── inactive/
├── no-pid/
└── json-path/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-status/` | `kck -h` mentions `codex status`. |
| `help/codex-usage/` | `kck codex --help` lists `status`. |
| `help/status-usage/` | status help documents path + `--json` + `--no-pid`. |
| `missing-session-id/` | `Error:` usage. |
| `unknown-session/` | `Error: codex session not found`. |
| `running/` | State running; Path line; pid present; File no. |
| `inactive/` | State inactive; Path line; PIDs none; File no. |
| `no-pid/` | PIDs skipped; inactive (no file signal). |
| `json-path/` | absolute `path` has rollout; `file_active` false; no ANSI. |

## How to Run

```sh
doctest vet ./tests/codex/status
doctest test ./tests/codex/status
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/agent-pro/agent/codex/sessions"
	"github.com/xhd2015/agent-pro/pkgs/procresolve"
	"github.com/xhd2015/doctest/session"

	"kck/run"
)

type FixtureProc struct {
	PID  int
	PPID int
	Cmd  string
}

type Request struct {
	Args      []string
	CodexHome string
	TempDir   string
	SessionID string
	Procs     []FixtureProc
	OpenFiles map[int][]string
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
	files := req.OpenFiles
	if files == nil {
		files = map[int][]string{}
	}
	live := &sessions.LiveOptions{
		ListProcs: func() []procresolve.Proc {
			out := make([]procresolve.Proc, 0, len(req.Procs))
			for _, p := range req.Procs {
				out = append(out, procresolve.Proc{PID: p.PID, PPID: p.PPID, Cmd: p.Cmd})
			}
			return out
		},
		Lsof: func(pid int) []string {
			return files[pid]
		},
	}
	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:          append([]string(nil), req.Args...),
		Stdout:        &stdout,
		Stderr:        &stderr,
		CodexHome:     req.CodexHome,
		CodexLiveOpts: live,
	})
	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
