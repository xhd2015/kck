# Scenario

**Feature**: extra positional args error

```
MainWith(["stray-arg"]) -> error (no bare subcommands this slice)
```

## Steps

1. Args = single unexpected positional.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"stray-arg"}
	return nil
}
```
