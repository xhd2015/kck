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
	assertContains(t, out, "hit-1", "stdout")
	assertContains(t, out, "hit-2", "stdout")
	if strings.Contains(out, "hit-0") {
		t.Fatalf("limit after grep must drop oldest hit:\n%s", out)
	}
	if strings.Contains(out, "miss") {
		t.Fatalf("grep must drop non-hits:\n%s", out)
	}
}
```
