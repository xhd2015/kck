# Scenario

**Feature**: busy pane with Agents attach is listed

```
busy grok Agents[iterm] -> SESSION_ID agent-sess-grok-1 in stdout
```

## Steps

1. Inject busyGrokAgentResult().

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.LiveResult = busyGrokAgentResult()
	return nil
}
```
