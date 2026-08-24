# Scenario

**Feature**: busy grok under agent-run process tree → AGENT_RUN yes

```
Tree role agent-run + grok Agent -> AGENT_RUN yes, AGENT_SID set
```

## Steps

1. Inject busyGrokUnderAgentRunResult().

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.LiveResult = busyGrokUnderAgentRunResult()
	return nil
}
```
