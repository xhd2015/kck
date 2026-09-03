# Assert

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "kck grok new", "stdout")
	assertContains(t, resp.Stdout, "--here", "stdout")
	assertContains(t, resp.Stdout, "--no-new-terminal", "stdout")
	assertContains(t, resp.Stdout, "--new-terminal", "stdout")
	assertContains(t, resp.Stdout, "--submit", "stdout")
	assertNotContains(t, resp.Stdout, "--new-window", "stdout")
	assertNotContains(t, resp.Stdout, "-n,", "stdout")
	assertNotContains(t, resp.Stdout, "--no-agent-run", "stdout")
}
```
