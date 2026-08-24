# Scenario

**Feature**: two live agents appear with correct footer counts

```
agent-sess-a (busy) + agent-sess-b (idle) -> both listed;
footer 2 sessions · 1 needs attention · 1 sendable
```

## Steps

1. Inject twoAgentsResult().
2. Default list args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = []string{}
	req.LiveResult = twoAgentsResult()
	return nil
}
```
