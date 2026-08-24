# Scenario

**Feature**: live row columns — ITERM, workspace, LIVE/SENDABLE/STATE heuristics

```
pane + Agent -> SESSION_ID RUNNER ITERM WORKSPACE
Idle=false -> busy / not sendable
Idle=true agent-like -> idle / sendable
```

## Steps

1. Live path.
2. Child injects specific pane fixtures.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Home = ""
	req.Args = []string{}
	return nil
}
```
