## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "codex resume", "help")
	assertContains(t, resp.Stdout, "--no-agent-run", "help")
	assertContains(t, resp.Stdout, "--tab", "help")
}
```
