# Scenario

**Feature**: list orders sessions by updated_at newest first

```
s-old (t=Aug1) + s-new (t=Aug2) -> first data row is s-new
```

## Steps

1. Use twoSessionFixture.
2. Default list args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	twoSessionFixture(req)
	req.Args = []string{}
	return nil
}
```
