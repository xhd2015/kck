# Scenario

**Feature**: live capture outcomes — empty inventory and soft-fail errors

```
LiveCapture empty Result -> 0 footer
LiveCapture error -> warning + 0 rows exit 0
```

## Steps

1. Home empty (live path).
2. Child sets LiveResult or LiveErr.

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
