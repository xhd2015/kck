# Scenario

**Feature**: `--limit N` caps listed rows

```
3 sessions + --limit 1 -> one data row (newest)
```

## Steps

1. Leaves use threeSessionFixture + limit flag.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	threeSessionFixture(req)
	return nil
}
```
