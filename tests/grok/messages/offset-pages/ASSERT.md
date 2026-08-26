## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	// total=8, offset=2, shown=2 → hi=6, lo=5
	assertContains(t, out, "Chat history (showing 5-6(2) of 8):", "stdout")
	assertContains(t, out, "[tool] : run_terminal_command: echo hi", "stdout")
	assertContains(t, out, "[assistant] : a1", "stdout")
	if strings.Contains(out, "u2") || strings.Contains(out, "a2") {
		t.Fatalf("offset page must skip newest pair:\n%s", out)
	}
}
```
