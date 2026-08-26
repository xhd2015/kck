## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "[assistant] : a1", "stdout")
	assertContains(t, out, "Chat history (1 message):", "stdout")
	for _, bad := range []string{"[user] :", "[thinking] :", "[tool] :", "a0", "a2", "u1"} {
		if strings.Contains(out, bad) {
			t.Fatalf("single grep must not include %q:\n%s", bad, out)
		}
	}
}
```
