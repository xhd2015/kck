# Scenario

**Feature**: bad CLI usage fails with Error (not list)

```
MainWith(bad args) -> non-zero; Error on stderr/err
```

## Steps

1. Leaf sets invalid argv.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	return nil
}
```
