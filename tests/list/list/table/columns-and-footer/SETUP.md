# Scenario

**Feature**: human list prints required columns and summary footer

```
two sessions -> header has all column names; footer N/M/K counts
```

## Steps

1. twoSessionFixture (1 needs attention, 1 sendable).
2. Default list.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	twoSessionFixture(req)
	req.Args = []string{}
	return nil
}
```
