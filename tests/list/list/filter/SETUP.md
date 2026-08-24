# Scenario

**Feature**: list filters --needs-confirm and --sendable

```
threeSessionFixture + filter flag -> subset of rows
```

## Steps

1. Leaves set threeSessionFixture and one filter flag.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	threeSessionFixture(req)
	return nil
}
```
