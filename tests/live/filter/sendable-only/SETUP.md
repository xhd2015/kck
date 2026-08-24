# Scenario

**Feature**: `--sendable` keeps only sendable (idle) live rows

```
agent-sess-a busy + agent-sess-b idle + --sendable -> only agent-sess-b
```

## Steps

1. Parent injects twoAgentsResult.
2. Args = `--sendable`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = []string{"--sendable"}
	return nil
}
```
