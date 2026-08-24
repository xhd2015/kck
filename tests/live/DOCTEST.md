# kck live list — default iTerm inventory (P4)

Classic TDD for **default `kck` live list**: when no store home is set,
capture live iTerm inventory via `agent-pro/pkgs/itermsnapshot` and print
agent rows (kool-aligned, **no kool import**). L2 in-process only —
inject `LiveCapture`. Sealed `./tests/list/` remains the store path and must
stay GREEN without modification.

```text
kck [OPTIONS]
  (default, no --home / Options.Home)  live list via itermsnapshot
  --home PATH                          store list (existing tests/list)
```

# DSN (Domain Specific Notion)

**Participants**

- **Caller** — user or harness; argv without program name; injects
  `Home`, writers, and **`LiveCapture`** via `run.MainWith`.
- **`kck/run`** — mode route: non-empty home → store list (Probe + ListITerm);
  empty home → **live list** (LiveCapture → rows).
- **LiveCapture** — returns `*itermsnapshot.Result`, warnings, error.
  Production default: `itermsnapshot.Capture`. L2 injects fixture Result.
- **Agents-only filter** — default live list shows panes that have an
  attached Agent **or** an agent-like command (`grok` / `codex` / `mark` /
  `agent-run`); plain idle bash/zsh is omitted.
- **Row mapper** — maps pane + optional Agent → list columns (SESSION_ID,
  RUNNER, LIVE, SENDABLE, STATE, REASON, ITERM, UPDATED, WORKSPACE).

**Behaviors**

- **Default = live** when `Options.Home` empty and args have no `--home`.
- **Store** when `Options.Home` or `--home PATH` is non-empty (existing
  list tree / Probe inject unchanged).
- Live human table reuses store column names and footer wording
  `N sessions · M needs attention · K sendable`.
- Capture **hard error** → soft-fail: `warning:` on stderr, empty list,
  exit 0 (`MainWith` nil). Capture **warnings** printed as `warning:` lines.
- `--json` live shape matches store JSON envelope (`sessions` + `summary`).
- `--needs-confirm` / `--sendable` / `--limit` apply to live rows after
  agents-only inclusion.

### Public API (Classic TDD — locked for implementer)

Extends existing `kck/run` Options. Store types (`Probe`, `ListITerm`, …)
unchanged. New field only:

```go
// Package: kck/run

import "github.com/xhd2015/agent-pro/pkgs/itermsnapshot"

// Options additions (existing fields retained):
type Options struct {
    Args      []string
    Home      string
    Stdout    io.Writer
    Stderr    io.Writer
    Probe     ProbeFunc
    ListITerm ListITermFunc
    Now       time.Time

    // LiveCapture injects itermsnapshot result for L2 live path.
    // nil → production: itermsnapshot.Capture(CaptureOpts{}).
    // Invoked only when resolved home is empty (live mode).
    // Signature: (result, warnings, error). Hard error is soft-failed by kck.
    LiveCapture func() (*itermsnapshot.Result, []string, error)
}

// Main / MainWith signatures unchanged.
func Main(args []string) error
func MainWith(opts Options) error
```

**Mode route (locked):**

| Condition | Path |
|-----------|------|
| `Options.Home` non-empty **or** `--home PATH` in Args | **store** list (do not call LiveCapture) |
| else | **live** list via LiveCapture |

No `--store` flag this phase (optional later).

**Agents-only inclusion (locked — default live filter):**

Include a pane when **any** of:

1. `Result.Agents[session.ID]` non-nil with non-empty `Kind` or `SessionID`
2. Agent-like command token appears (case-insensitive) in session `Name`,
   `Command`, `CommandLine`, or any `Processes[].Command` as a path basename
   or whitespace-delimited token among: **`grok`**, **`codex`**, **`mark`**,
   **`agent-run`**

Exclude plain shell panes (e.g. idle `zsh`/`bash` with no Agent and no
agent-like token).

**Live row mapping (locked):**

| Column | Source |
|--------|--------|
| SESSION_ID | Agent.SessionID if non-empty; else iTerm session `ID` (full) |
| RUNNER | Agent.Kind if non-empty; else matched agent-like token; else `-` |
| LIVE | `yes` (inventory pane is live) |
| SENDABLE | `yes` iff `Idle != nil && *Idle`; else `no` |
| STATE | `idle` / `busy` / `unknown` from Idle pointer |
| REASON | `-` when empty (no special live reason this phase) |
| ITERM | `w=<id> t=<tabIndex>`; `<id>` = WindowID if non-zero else window Index; tabIndex = Tab.Index (1-based). `--no-iterm` forces `-` |
| UPDATED | `Duration` string if set; else `-` |
| WORKSPACE | `Cwd` if set; else `-` |

**needs_attention (same D1):** LIVE && !SENDABLE && !Exited. Live panes are
never Exited → busy / unknown Idle ⇒ needs attention; idle sendable ⇒ not.

**Capture soft-fail (locked):**

| Outcome | Stderr | Stdout | Exit |
|---------|--------|--------|------|
| error from LiveCapture | `warning: live capture failed: …` (msg may follow colon) | header + `0 sessions · 0 needs attention · 0 sendable` | 0 |
| warnings only | each warning as `warning: <text>` | rows as usual | 0 |
| empty Result / nil Snapshot / zero windows | (none required) | 0-count footer | 0 |

**Stable messages / formats:** same LIVE/SENDABLE yes/no, ITERM `w=… t=…`,
footer wording as store path.

## Version

0.0.2

## Decision Tree

```
tests/live/
├── DOCTEST.md
├── SETUP.md
├── mode/
│   └── with-home-store/          Home set → store path; LiveCapture not called
├── capture/
│   ├── empty/                    empty Result → 0 footer
│   └── fail-soft/                LiveCapture err → warning + 0 rows exit 0
├── include/                      agents-only policy
│   ├── busy-agent-shown/         Agents hit → row present
│   ├── agent-like-cmd-shown/     no Agents; command mark → row
│   └── plain-bash-omitted/       idle zsh, no agent → not listed
├── row/
│   ├── iterm-workspace/          ITERM w/t, WORKSPACE cwd, SESSION_ID, RUNNER
│   ├── busy-state/               busy → SENDABLE no, STATE busy, needs attention
│   └── idle-sendable/            idle agent-like → SENDABLE yes, STATE idle
├── multi/
│   └── two-agents-footer/        two included rows + footer counts
├── filter/
│   ├── needs-confirm/            --needs-confirm only attention
│   └── sendable-only/            --sendable only sendable
├── limit/
│   └── max-rows/                 --limit 1 of 2
└── json/
    └── live-shape/               --json sessions+summary, no ANSI
```

Parameter ranking (most → least significant):

1. **Mode** — live vs store (home present)
2. **Capture outcome** — empty / fail-soft / success inventory
3. **Inclusion** — agents-only (agent / agent-like / omit plain)
4. **Row fields** — iterm+cwd, busy vs idle heuristics
5. **Multiplicity / modifiers** — multi footer, needs-confirm, sendable, limit
6. **Format** — human vs json

## Test Index

| # | Leaf | Description | Expect |
|---|------|-------------|--------|
| 1 | `mode/with-home-store` | Home set → store list; LiveCapture not invoked | RED |
| 2 | `capture/empty` | empty live Result → 0 footer; exit 0 | RED |
| 3 | `capture/fail-soft` | LiveCapture error → `warning:` + 0 rows; exit 0 | RED |
| 4 | `include/busy-agent-shown` | busy + Agents[id] → session id in stdout | RED |
| 5 | `include/agent-like-cmd-shown` | no Agents; Command=mark → row shown | RED |
| 6 | `include/plain-bash-omitted` | idle zsh only → not listed (0 or no shell id) | RED |
| 7 | `row/iterm-workspace` | ITERM `w=42 t=3`, WORKSPACE cwd, RUNNER grok | RED |
| 8 | `row/busy-state` | busy agent → LIVE yes SENDABLE no STATE busy; footer attention | RED |
| 9 | `row/idle-sendable` | idle agent-like → SENDABLE yes STATE idle | RED |
| 10 | `multi/two-agents-footer` | two agents → both ids + footer | RED |
| 11 | `filter/needs-confirm` | only busy/attention row | RED |
| 12 | `filter/sendable-only` | only idle sendable row | RED |
| 13 | `limit/max-rows` | `--limit 1` shows one of two | RED |
| 14 | `json/live-shape` | valid JSON sessions+summary; no ANSI | RED |

## How to Run

```sh
# from kck module root
doctest vet ./tests/live
doctest test ./tests/live
doctest test -v ./tests/live/include/busy-agent-shown

# sealed store path must stay green
doctest test ./tests/list
```

Expect **RED** until implementer adds `Options.LiveCapture`, live mode route,
agents-only filter, and live row mapping. `./tests/list` must remain GREEN.

```go
import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/agent-pro/pkgs/itermsnapshot"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2/snapshot"

	"kck/run"
)

// SessionSeed seeds store meta for mode/with-home-store only.
type SessionSeed struct {
	SessionID       string
	Runner          string
	RunnerSessionID string
	Workspace       string
	UpdatedAt       string
	Live            bool
	Sendable        bool
	State           string
	TTY             string
}

// Request drives one MainWith invocation (live or store route).
type Request struct {
	Args []string

	// Home non-empty → store path (LiveCapture must not be required).
	Home string

	// Live inject (live path only when Home empty and no --home in Args).
	LiveResult   *itermsnapshot.Result
	LiveWarnings []string
	LiveErr      string // non-empty → LiveCapture returns error with this text

	// LiveCaptureCalled, when non-nil, is set true if LiveCapture runs.
	LiveCaptureCalled *bool

	// Store path extras (mode/with-home-store).
	Sessions  []SessionSeed
	ProbeByID map[string]run.ProbeResult
	SkipITerm bool

	Now time.Time
}

// Response is the observable MainWith outcome.
type Response struct {
	Stdout   string
	Stderr   string
	ErrText  string
	ExitCode int // 0 if err==nil else 1
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d

	if req.Home != "" {
		if err := seedSessions(req.Home, req.Sessions); err != nil {
			return nil, fmt.Errorf("seed sessions: %w", err)
		}
	}

	var stdout, stderr bytes.Buffer
	opts := run.Options{
		Args:   append([]string(nil), req.Args...),
		Home:   req.Home,
		Stdout: &stdout,
		Stderr: &stderr,
		Now:    req.Now,
	}

	// LiveCapture inject: always wire when we need spy and/or live path data.
	// Product must only invoke it on live path (empty resolved home).
	needLiveInject := req.Home == "" || req.LiveCaptureCalled != nil || req.LiveErr != "" || req.LiveResult != nil
	if needLiveInject {
		opts.LiveCapture = func() (*itermsnapshot.Result, []string, error) {
			if req.LiveCaptureCalled != nil {
				*req.LiveCaptureCalled = true
			}
			if req.LiveErr != "" {
				return nil, nil, fmt.Errorf("%s", req.LiveErr)
			}
			return req.LiveResult, append([]string(nil), req.LiveWarnings...), nil
		}
	}

	if req.Home != "" {
		probeMap := mergeProbeMap(req)
		opts.Probe = func(sessionID string, meta run.SessionMetaView) (run.ProbeResult, error) {
			if p, ok := probeMap[sessionID]; ok {
				return p, nil
			}
			return run.ProbeResult{State: "unknown"}, nil
		}
		if !req.SkipITerm {
			opts.ListITerm = func() ([]run.ITermSession, error) {
				return nil, nil
			}
		}
	}

	err := run.MainWith(opts)
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

func seedSessions(home string, seeds []SessionSeed) error {
	for _, s := range seeds {
		id := strings.TrimSpace(s.SessionID)
		if id == "" {
			return fmt.Errorf("session seed missing SessionID")
		}
		dir := filepath.Join(home, "sessions", id)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		runner := s.Runner
		if runner == "" {
			runner = "grok"
		}
		updated := s.UpdatedAt
		if updated == "" {
			updated = "2026-08-01T12:00:00Z"
		}
		meta := map[string]any{
			"runner":     runner,
			"session_id": id,
			"status":     "running",
			"workspace":  s.Workspace,
			"updated_at": updated,
			"created_at": updated,
		}
		if sid := strings.TrimSpace(s.RunnerSessionID); sid != "" {
			meta["runner_session_id"] = sid
		}
		data, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, "meta.json"), data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func mergeProbeMap(req *Request) map[string]run.ProbeResult {
	out := make(map[string]run.ProbeResult)
	for k, v := range req.ProbeByID {
		out[k] = v
	}
	for _, s := range req.Sessions {
		if _, ok := out[s.SessionID]; ok {
			continue
		}
		if s.Live || s.Sendable || s.State != "" || s.TTY != "" {
			out[s.SessionID] = run.ProbeResult{
				Live:     s.Live,
				Sendable: s.Sendable,
				State:    s.State,
				TTY:      s.TTY,
			}
		}
	}
	return out
}

// --- live fixture helpers ---

func boolPtr(v bool) *bool       { return &v }
func intPtr(v int) *int          { return &v }
func strPtr(v string) *string    { return &v }
func uint64Val(v uint64) uint64  { return v }

// liveResult wraps Snapshot + Agents into itermsnapshot.Result.
func liveResult(snap *snapshot.Snapshot, agents map[string]*itermsnapshot.SessionAgent) *itermsnapshot.Result {
	return &itermsnapshot.Result{Snapshot: snap, Agents: agents}
}

// onePaneSnap builds a one-window / one-tab / one-session Snapshot.
// windowID 0 → product falls back to window Index for ITERM.
func onePaneSnap(sessID, name, tty string, windowID uint64, winIndex, tabIndex int, idle *bool, cwd *string, command *string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		CapturedAt: "2026-08-06T12:00:00Z",
		Host:       "testhost",
		Source:     "iterm2",
		Summary:    snapshot.SnapshotSummary{Windows: 1, Tabs: 1, Sessions: 1},
		Windows: []snapshot.SnapshotWindow{
			{
				Index:    winIndex,
				Name:     "W",
				WindowID: windowID,
				Tabs: []snapshot.SnapshotTab{
					{
						Index: tabIndex,
						Name:  "T",
						Sessions: []snapshot.SnapshotSession{
							{
								Index:   1,
								ID:      sessID,
								Name:    name,
								TTY:     tty,
								Profile: "Default",
								Idle:    idle,
								Cwd:     cwd,
								Command: command,
								PID:     intPtr(9001),
							},
						},
					},
				},
			},
		},
	}
}

// busyGrokAgentResult: busy pane with Agents attach (primary happy path).
// Tree has grok only (no agent-run) → AGENT_RUN=no when mapped.
func busyGrokAgentResult() *itermsnapshot.Result {
	const itermSess = "iterm-uuid-busy-grok"
	const agentSess = "agent-sess-grok-1"
	snap := onePaneSnap(itermSess, "grok-work", "/dev/ttys011", 42, 1, 3, boolPtr(false), strPtr("/ws/grok"), strPtr("grok"))
	return liveResult(snap, map[string]*itermsnapshot.SessionAgent{
		itermSess: {
			Kind:      "grok",
			SessionID: agentSess,
			Title:     "Grok Work",
			Tree: []itermsnapshot.AgentTreeNode{
				{PID: 9001, PPID: 1, Role: "input", Cmd: "zsh"},
				{PID: 9002, PPID: 9001, Role: "grok", Cmd: "grok"},
			},
		},
	})
}

// busyGrokUnderAgentRunResult: tree includes agent-run → AGENT_RUN=yes.
func busyGrokUnderAgentRunResult() *itermsnapshot.Result {
	const itermSess = "iterm-uuid-ar-grok"
	const agentSess = "agent-sess-ar-grok"
	snap := onePaneSnap(itermSess, "grok-ar", "/dev/ttys012", 50, 1, 1, boolPtr(false), strPtr("/ws/ar"), strPtr("grok"))
	return liveResult(snap, map[string]*itermsnapshot.SessionAgent{
		itermSess: {
			Kind:      "grok",
			SessionID: agentSess,
			Tree: []itermsnapshot.AgentTreeNode{
				{PID: 8001, PPID: 1, Role: "input", Cmd: "zsh"},
				{PID: 8002, PPID: 8001, Role: "agent-run", Cmd: "agent-run run"},
				{PID: 8003, PPID: 8002, Role: "grok", Cmd: "grok"},
			},
		},
	})
}

// idleMarkCmdResult: idle pane, no Agents, Command=mark (agent-like inclusion).
func idleMarkCmdResult() *itermsnapshot.Result {
	const itermSess = "iterm-uuid-idle-mark"
	snap := onePaneSnap(itermSess, "mark-pane", "/dev/ttys020", 7, 1, 1, boolPtr(true), strPtr("/ws/mark"), strPtr("mark"))
	return liveResult(snap, nil)
}

// idleZshOnlyResult: plain idle shell — must be omitted by agents-only.
func idleZshOnlyResult() *itermsnapshot.Result {
	const itermSess = "iterm-uuid-plain-zsh"
	snap := onePaneSnap(itermSess, "shell", "/dev/ttys030", 9, 1, 1, boolPtr(true), strPtr("/tmp"), strPtr("zsh"))
	return liveResult(snap, nil)
}

// twoAgentsResult: two busy agent panes for multi/limit/filter leaves.
func twoAgentsResult() *itermsnapshot.Result {
	const (
		idA = "iterm-uuid-agent-a"
		idB = "iterm-uuid-agent-b"
		sa  = "agent-sess-a"
		sb  = "agent-sess-b"
	)
	snap := &snapshot.Snapshot{
		CapturedAt: "2026-08-06T12:00:00Z",
		Host:       "testhost",
		Source:     "iterm2",
		Summary:    snapshot.SnapshotSummary{Windows: 1, Tabs: 1, Sessions: 2},
		Windows: []snapshot.SnapshotWindow{
			{
				Index:    1,
				WindowID: 10,
				Tabs: []snapshot.SnapshotTab{
					{
						Index: 2,
						Sessions: []snapshot.SnapshotSession{
							{
								Index:   1,
								ID:      idA,
								Name:    "a",
								TTY:     "/dev/ttys001",
								Idle:    boolPtr(false), // busy → not sendable
								Cwd:     strPtr("/wa"),
								Command: strPtr("grok"),
								PID:     intPtr(1001),
							},
							{
								Index:   2,
								ID:      idB,
								Name:    "b",
								TTY:     "/dev/ttys002",
								Idle:    boolPtr(true), // idle → sendable (agent still attached for inclusion)
								Cwd:     strPtr("/wb"),
								Command: strPtr("codex"),
								PID:     intPtr(1002),
							},
						},
					},
				},
			},
		},
	}
	// Note: product agents-only includes via Agents map even if Idle=true
	// (attach is busy-only in itermsnapshot production; L2 injects Agents freely).
	return liveResult(snap, map[string]*itermsnapshot.SessionAgent{
		idA: {Kind: "grok", SessionID: sa},
		idB: {Kind: "codex", SessionID: sb},
	})
}
```
