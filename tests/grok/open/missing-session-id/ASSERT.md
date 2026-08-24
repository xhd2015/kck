## Expected

- Usage error with Error: prefix; no side effects.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertContains(t, resp.Stderr, "Error:", "stderr")
	if !strings.Contains(resp.Stderr, "session id") && !strings.Contains(resp.ErrText, "session id") {
		t.Fatalf("want session id usage error; stderr=%q err=%q", resp.Stderr, resp.ErrText)
	}
	if len(resp.Focused) != 0 || len(resp.Opened) != 0 {
		t.Fatalf("side effects: focused=%v opened=%v", resp.Focused, resp.Opened)
	}
}
```
