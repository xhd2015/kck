## Expected

- Non-zero exit; stderr mentions `kck skill:` and unknown/missing topic.

## Exit Code

- 1

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	blob := resp.Stderr + " " + resp.ErrText
	if !strings.Contains(blob, "kck skill:") {
		t.Fatalf("stderr/err missing kck skill: prefix: stderr=%q err=%q", resp.Stderr, resp.ErrText)
	}
	lower := strings.ToLower(blob)
	if !strings.Contains(lower, "does not exist") && !strings.Contains(lower, "not found") && !strings.Contains(lower, "unknown") && !strings.Contains(lower, "read skill") {
		t.Fatalf("stderr should mention missing topic: %q", blob)
	}
}
```
