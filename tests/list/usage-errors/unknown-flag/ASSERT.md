## Expected

- Non-zero exit / product error.
- Error text or stderr mentions the unknown flag or is a parse error (assertable
  non-empty failure related to flag).

## Errors

- Parse / unknown flag error.

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
	if blob == "" {
		t.Fatal("expected error text for unknown flag")
	}
	// less-flags / CLI should mention the bad flag or "unknown"/"unrecognized".
	if !strings.Contains(blob, "not-a-real-kck-flag") &&
		!strings.Contains(blob, "unknown") &&
		!strings.Contains(blob, "unrecognized") &&
		!strings.Contains(blob, "flag") {
		t.Fatalf("want flag-related error; stderr=%q err=%q", resp.Stderr, resp.ErrText)
	}
}
```
