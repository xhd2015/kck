# Scenario

**Feature**: LiveCapture hard error soft-fails with warning and empty list

```
LiveCapture err "iTerm not running" -> warning: on stderr, 0 footer, exit 0
```

## Steps

1. Set LiveErr to a stable capture failure message.
2. Leave LiveResult nil.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.LiveErr = "iTerm not running"
	req.LiveResult = nil
	return nil
}
```
