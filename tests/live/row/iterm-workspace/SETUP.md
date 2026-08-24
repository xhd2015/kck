# Scenario

**Feature**: live row shows ITERM w/t, WORKSPACE cwd, agent SESSION_ID and RUNNER

```
WindowID=42 Tab.Index=3 Cwd=/ws/grok Agent{grok, agent-sess-grok-1}
  -> ITERM w=42 t=3, WORKSPACE /ws/grok, SESSION_ID, RUNNER
```

## Steps

1. Inject busyGrokAgentResult() (window 42, tab 3, cwd /ws/grok).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.LiveResult = busyGrokAgentResult()
	return nil
}
```
