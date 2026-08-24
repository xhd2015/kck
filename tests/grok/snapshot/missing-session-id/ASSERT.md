## Expected

- Usage error with Error: prefix; no Contents.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	combined := resp.Stderr + "\n" + resp.ErrText
	if !strings.Contains(combined, "Error:") {
		t.Fatalf("want Error: prefix; stderr=%q err=%q", resp.Stderr, resp.ErrText)
	}
	assertNoContents(t, resp)
}
```
