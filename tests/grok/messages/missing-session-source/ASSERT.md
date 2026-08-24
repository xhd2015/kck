## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	combined := resp.ErrText + resp.Stderr
	if !strings.Contains(combined, "Error:") {
		t.Fatalf("want Error: prefix; err=%q stderr=%q", resp.ErrText, resp.Stderr)
	}
	if !strings.Contains(combined, "session id") && !strings.Contains(combined, "--tab") {
		t.Fatalf("want session source usage; got %q", combined)
	}
}
```
