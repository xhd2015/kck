## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	if !strings.Contains(resp.Stderr+"\n"+resp.ErrText, "Error:") {
		t.Fatalf("want Error:; got %q / %q", resp.Stderr, resp.ErrText)
	}
}
```
