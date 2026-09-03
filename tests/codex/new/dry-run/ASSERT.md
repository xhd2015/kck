# Assert

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "Would open new codex session", "stdout")
	assertContains(t, resp.Stdout, "runner: codex-tty", "stdout")
	assertContains(t, resp.Stdout, "$brainstorm extract TODOs", "stdout")
	sid := expectedCodexSessionID("extract TODOs")
	assertContains(t, resp.Stdout, "agent-run-session-id: "+sid, "stdout")
}
```
