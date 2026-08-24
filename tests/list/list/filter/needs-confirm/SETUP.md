# Scenario

**Feature**: `--needs-confirm` lists only needs_attention sessions

```
sess-sendable, sess-attention, sess-exited + --needs-confirm -> only sess-attention
```

## Steps

1. threeSessionFixture (from parent).
2. Args = `--needs-confirm`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--needs-confirm"}
	return nil
}
```
