## Expected

- Exit 0.
- Stdout contains store session id `store-sess-1`.
- Stdout does **not** contain live agent id `agent-sess-grok-1`.
- `LiveCaptureCalled` remains false (store path must not invoke LiveCapture).

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "store-sess-1", "store session")
	assertNotContains(t, resp.Stdout, "agent-sess-grok-1", "live agent must not appear")
	if req.LiveCaptureCalled == nil {
		t.Fatal("LiveCaptureCalled spy missing")
	}
	if *req.LiveCaptureCalled {
		t.Fatal("LiveCapture must not be called when Home is set")
	}
}
```
