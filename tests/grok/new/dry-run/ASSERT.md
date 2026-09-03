# Assert

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "Would open new grok session", "stdout")
	assertContains(t, resp.Stdout, "terminal: new", "stdout")
	assertContains(t, resp.Stdout, "runner: grok-tty", "stdout")
	assertContains(t, resp.Stdout, "/brainstorm fix flaky auth", "stdout")
	assertContains(t, resp.Stdout, "--new-terminal", "stdout")
	assertContains(t, resp.Stdout, "--no-submit", "stdout")
	sid := expectedSessionID("fix flaky auth")
	assertContains(t, resp.Stdout, "agent-run-session-id: "+sid, "stdout")
	if len(resp.Foreground) != 0 || len(resp.NewTerminal) != 0 {
		t.Fatalf("dry-run must not launch: fg=%v nt=%v", resp.Foreground, resp.NewTerminal)
	}
}
```
