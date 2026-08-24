# Scenario

**Feature**: agents-only inclusion — Agent, agent-like command, omit plain shell

```
Agents[id] -> include
Command mark (no Agents) -> include
idle zsh only -> omit
```

## Steps

1. Live path (empty Home).
2. Child injects LiveResult for the inclusion case.

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
