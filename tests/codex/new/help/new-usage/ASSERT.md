# Assert

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "kck codex new", "stdout")
	assertContains(t, resp.Stdout, "$brainstorm", "stdout")
	assertContains(t, resp.Stdout, "codex-tty", "stdout")
}
```
