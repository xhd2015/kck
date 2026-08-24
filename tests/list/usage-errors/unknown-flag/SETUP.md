# Scenario

**Feature**: unknown flag errors

```
MainWith(["--not-a-real-kck-flag"]) -> error
```

## Steps

1. Args = unknown long flag.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--not-a-real-kck-flag"}
	return nil
}
```
