# Scenario

**Feature**: `--limit 1` shows at most one of two live agents

```
twoAgentsResult + --limit 1 -> exactly one of agent-sess-a / agent-sess-b
```

## Steps

1. Inject twoAgentsResult.
2. Args = `--limit 1`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = []string{"--limit", "1"}
	req.LiveResult = twoAgentsResult()
	return nil
}
```
