## Expected

- Non-zero exit / product error about extra/unrecognized args.

## Errors

- Extra positional error.

## Exit Code

non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertFailure(t, resp)
	blob := strings.ToLower(resp.Stderr + " " + resp.ErrText)
	if !strings.Contains(blob, "stray-arg") &&
		!strings.Contains(blob, "extra") &&
		!strings.Contains(blob, "unrecognized") &&
		!strings.Contains(blob, "unexpected") &&
		!strings.Contains(blob, "positional") {
		t.Fatalf("want extra-arg error; stderr=%q err=%q", resp.Stderr, resp.ErrText)
	}
}
```
