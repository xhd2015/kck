# Scenario

**Feature**: mode route — store when Home set, live when Home empty

```
Home non-empty -> store list (LiveCapture not called)
Home empty -> live list (LiveCapture called)
```

## Steps

1. Grouping for mode-route leaves.
2. Child sets Home and/or LiveCapture spy.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	_ = req
	return nil
}
```
