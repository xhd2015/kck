# Scenario

**Feature**: AGENT_SID column shows runner-native session id

```
Agent.SessionID=agent-sess-grok-1 -> AGENT_SID and SESSION_ID show that id
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
