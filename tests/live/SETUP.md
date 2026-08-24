# Scenario

**Feature**: kck default list is live iTerm inventory via injectable LiveCapture

```
# live path (empty Home)
Caller -> run.MainWith(Options{Args, LiveCapture, Stdout, Stderr})
  -> LiveCapture() -> itermsnapshot.Result
  -> agents-only filter -> human table | JSON

# store path (Home set) — LiveCapture not called
Caller -> MainWith(Home) -> session store + Probe (tests/list contract)
```

## Preconditions

- Package `kck/run` exposes `MainWith(Options)` with **`LiveCapture`** locked in
  root `DOCTEST.md` (Classic TDD: RED until implementer lands field + live path).
- Sealed `./tests/list` uses store path via `Options.Home` — must stay GREEN.
- Parallel-safe: no `t.Setenv` / `os.Setenv` / `os.Chdir`; live fixtures are
  pure in-memory Results.
- L2 only: no real AppleScript, ps, or iTerm.

## Steps

1. Root `Setup` zeros live inject fields; leaves set `LiveResult` / `LiveErr` /
   `Home` as needed.
2. Root `Run` wires `LiveCapture` inject and calls `run.MainWith`.
3. Leaf `Assert` checks stdout/stderr/exit only.

## Context

- Agents-only filter (default): Agent attach **or** agent-like command token.
- Capture hard error → soft `warning:` + empty list exit 0.
- Footer: `N sessions · M needs attention · K sendable`.
- D1 needs_attention = live && !sendable && !exited.

```go
import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"

	"kck/run"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Live path by default: empty Home (do not invent a temp store home).
	if req.ProbeByID == nil {
		req.ProbeByID = map[string]run.ProbeResult{}
	}
	if req.Now.IsZero() {
		req.Now = time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC)
	}
	_ = filepath.Separator // keep import useful for store mode leaves that set Home
	return nil
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
}

func assertSuccess(t *testing.T, resp *Response) {
	t.Helper()
	if resp.ExitCode != 0 || resp.ErrText != "" {
		t.Fatalf("want success exit 0, got exit=%d err=%q stderr=%q stdout=%q",
			resp.ExitCode, resp.ErrText, resp.Stderr, resp.Stdout)
	}
}

func assertContains(t *testing.T, got, want, label string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("%s must contain %q; got:\n%s", label, want, got)
	}
}

func assertNotContains(t *testing.T, got, ban, label string) {
	t.Helper()
	if strings.Contains(got, ban) {
		t.Fatalf("%s must not contain %q; got:\n%s", label, ban, got)
	}
}

func assertNoANSI(t *testing.T, s, label string) {
	t.Helper()
	if strings.Contains(s, "\x1b[") {
		t.Fatalf("%s must not contain ANSI escapes; got:\n%q", label, s)
	}
}

func assertTrailingNewline(t *testing.T, s, label string) {
	t.Helper()
	if s == "" || !strings.HasSuffix(s, "\n") {
		t.Fatalf("%s must end with trailing newline; got %q", label, s)
	}
}

func assertWarningPrefix(t *testing.T, stderr string) {
	t.Helper()
	if !strings.Contains(strings.ToLower(stderr), "warning:") {
		t.Fatalf("want warning: on stderr; got %q", stderr)
	}
}
```
