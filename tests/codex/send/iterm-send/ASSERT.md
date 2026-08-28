## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "sent to session "+req.SessionID, "stdout")
	// Codex iTerm submit: type without newline, then two bare Enter writes.
	if len(resp.SendCalls) != 3 {
		t.Fatalf("SendCalls=%v", resp.SendCalls)
	}
	if !strings.Contains(resp.SendCalls[0].Text, "hello") {
		t.Fatalf("sent text=%q", resp.SendCalls[0].Text)
	}
	if !resp.SendCalls[0].Opts.NoSubmit {
		t.Fatalf("first write must stage (NoSubmit): %+v", resp.SendCalls[0].Opts)
	}
	for i, c := range resp.SendCalls[1:] {
		if c.Text != "" || c.Opts.NoSubmit {
			t.Fatalf("enter[%d]=%+v", i, c)
		}
	}
}
```
