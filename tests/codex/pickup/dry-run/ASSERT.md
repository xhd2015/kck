# Assert

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "Would pickup from "+req.SessionID, "stdout")
	assertContains(t, resp.Stdout, "terminal: current", "stdout")
	assertContains(t, resp.Stdout, "extract TODOs", "stdout")
	assertContains(t, resp.Stdout, "~/.cache/kck-pickup-a-session/SKILL.md", "stdout")
	if len(resp.Opened) != 0 || len(resp.Foreground) != 0 {
		t.Fatalf("dry-run must not launch: opened=%v fg=%v", resp.Opened, resp.Foreground)
	}
}
```
