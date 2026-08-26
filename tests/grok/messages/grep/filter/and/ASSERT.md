## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "[tool] : run_terminal_command: echo hi", "stdout")
	assertContains(t, out, "Chat history (1 message):", "stdout")
	for _, bad := range []string{"[user] :", "[thinking] :", "[assistant] :"} {
		if strings.Contains(out, bad) {
			t.Fatalf("AND grep must keep only the tool line; saw %q:\n%s", bad, out)
		}
	}
}
```
