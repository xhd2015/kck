## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "[user] : u0", "stdout")
	assertContains(t, out, "[thinking] : th0", "stdout")
	assertContains(t, out, "[tool] :", "stdout")
	assertContains(t, out, "[assistant] : a0", "stdout")
	ui := strings.Index(out, "[user]")
	ti := strings.Index(out, "[thinking]")
	to := strings.Index(out, "[tool]")
	ai := strings.Index(out, "[assistant]")
	if !(ui < ti && ti < to && to < ai) {
		t.Fatalf("kind order wrong:\n%s", out)
	}
}
```
