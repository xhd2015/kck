# Scenario

**Feature**: `--limit 1` shows only the newest session

```
sess-sendable (newest), sess-attention, sess-exited + --limit 1 -> only sess-sendable
```

## Steps

1. Parent threeSessionFixture.
2. Args = `--limit 1`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--limit", "1"}
	return nil
}
```
