# Scenario

**Feature**: idle agent-like pane is SENDABLE yes and STATE idle

```
Idle=true Command=mark -> SENDABLE yes STATE idle; footer 1 sendable
```

## Steps

1. Inject idleMarkCmdResult().

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.LiveResult = idleMarkCmdResult()
	return nil
}
```
