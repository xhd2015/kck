# Scenario

**Feature**: busy bare grok (no agent-run in tree) → AGENT_RUN no

```
Tree grok only + Agent -> AGENT_RUN no, AGENT_SID still set
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
