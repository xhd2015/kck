## Expected

- Exit 0.
- Exactly two SendText calls.
- Stdout contains two `sent to session` and `cron done: until reached`.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.SendCalls) != 2 {
		t.Fatalf("SendCalls=%d, want 2; stdout:\n%s", len(resp.SendCalls), resp.Stdout)
	}
	if strings.Count(resp.Stdout, "sent to session") != 2 {
		t.Fatalf("want 2 sent lines:\n%s", resp.Stdout)
	}
	assertContains(t, resp.Stdout, "cron done: until reached", "stdout")
	assertContains(t, resp.Stdout, "next ", "stdout")
}
```
