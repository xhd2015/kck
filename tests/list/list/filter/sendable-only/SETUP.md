# Scenario

**Feature**: `--sendable` lists only sendable sessions

```
three fixtures + --sendable -> only sess-sendable
```

## Steps

1. threeSessionFixture (parent).
2. Args = `--sendable`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--sendable"}
	return nil
}
```
