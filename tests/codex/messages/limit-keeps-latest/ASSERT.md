## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "showing last 2 of 3", "stdout")
	assertContains(t, out, "m1", "stdout")
	assertContains(t, out, "m2", "stdout")
	if strings.Contains(out, "[user] : m0") {
		t.Fatalf("m0 should be dropped:\n%s", out)
	}
}
```
