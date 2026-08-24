# Scenario

**Feature**: pane without Agents but agent-like Command is listed

```
Idle mark Command, Agents=nil -> row for iterm-uuid-idle-mark (or mark token)
```

## Steps

1. Inject idleMarkCmdResult() (no Agents map entry).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.LiveResult = idleMarkCmdResult()
	return nil
}
```
