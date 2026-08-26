## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	combined := resp.ErrText + resp.Stderr
	if !strings.Contains(combined, "--grep") {
		t.Fatalf("want --grep in error:\nerr=%q stderr=%q", resp.ErrText, resp.Stderr)
	}
	if !strings.Contains(strings.ToLower(combined), "empty") {
		t.Fatalf("want empty-pattern wording:\n%s", combined)
	}
}
```
