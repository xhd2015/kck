## Expected

- Exit 0 (schedule completes).
- `warning: send failed:` on stderr.
- Three SendText attempts (2nd failed inside SendText; SendCalls may be 2 if failure is before append — check wrapper).

Note: failure happens inside wrapped SendText after Fake would record — depend on whether orig is called. Our wrapper returns before orig, so SendCalls==2 (ticks 1 and 3).

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if !strings.Contains(resp.Stderr, "warning: send failed:") {
		t.Fatalf("want warning on stderr; got %q\nstdout:\n%s", resp.Stderr, resp.Stdout)
	}
	if len(resp.SendCalls) != 2 {
		t.Fatalf("SendCalls=%d, want 2 successful (tick 2 failed before record); stdout:\n%s",
			len(resp.SendCalls), resp.Stdout)
	}
	assertContains(t, resp.Stdout, "cron done: until reached", "stdout")
}
```
