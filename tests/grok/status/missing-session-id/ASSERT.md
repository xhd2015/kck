## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	combined := resp.Stderr + "\n" + resp.ErrText
	if !strings.Contains(combined, "Error:") {
		t.Fatalf("want Error: prefix; got %q", combined)
	}
}
```
