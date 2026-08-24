# Scenario

**Feature**: `--needs-confirm` keeps only needs_attention live rows

```
agent-sess-a busy + agent-sess-b idle + --needs-confirm -> only agent-sess-a
```

## Steps

1. Parent injects twoAgentsResult.
2. Args = `--needs-confirm`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = []string{"--needs-confirm"}
	return nil
}
```
