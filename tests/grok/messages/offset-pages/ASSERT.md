## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "Chat history (showing 2 of", "stdout")
	assertContains(t, out, "[tool] :", "stdout")
	assertContains(t, out, "[assistant] : a1", "stdout")
	if strings.Contains(out, "u2") || strings.Contains(out, "a2") {
		t.Fatalf("offset page must skip newest pair:\n%s", out)
	}
}
```
