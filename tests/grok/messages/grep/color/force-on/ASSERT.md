## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	if !strings.Contains(out, "\x1b[1m\x1b[31m") {
		t.Fatalf("want bold-red match highlight:\n%q", out)
	}
	if !strings.Contains(out, "\x1b[2m") {
		t.Fatalf("want dim timestamp:\n%q", out)
	}
	assertContains(t, out, "a1", "stdout")
}
```
