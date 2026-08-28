## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "--session-id", "help")
	assertContains(t, resp.Stdout, "--open", "help")
	assertContains(t, resp.Stdout, "--enter", "help")
	assertContains(t, resp.Stdout, "--cron", "help")
}
```
