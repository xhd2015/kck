# Scenario

**Feature**: `-h` documents list flags (no send stub)

```
MainWith(["-h"]) -> stdout mentions --home, --json, --needs-confirm, --sendable,
  --no-iterm, --fast; does not advertise list --send/--session
```

## Steps

1. Args = `["-h"]`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"-h"}
	return nil
}
```
