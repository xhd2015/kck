# Scenario

**Feature**: default mode lists sessions (no --send)

```
MainWith([list flags...], Home, Probe, ListITerm) -> table | JSON | filtered rows
```

## Steps

1. Leaves under this node do not pass `--send`.
2. Home is isolated; probe/iTerm injected unless leaf opts out.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Default list: no args unless leaf adds flags.
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```
