## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	combined := resp.ErrText + resp.Stderr
	if !strings.Contains(combined, "cannot be specified together") {
		t.Fatalf("want color conflict error:\n%s", combined)
	}
}
```
