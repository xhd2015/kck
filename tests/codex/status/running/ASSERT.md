## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "running", "stdout")
	assertContains(t, out, "Path:", "stdout")
	assertContains(t, out, "rollout-", "stdout")
	assertContains(t, out, "5001", "stdout")
	if !strings.Contains(out, "File: no") {
		t.Fatalf("want File: no:\n%s", out)
	}
}
```
