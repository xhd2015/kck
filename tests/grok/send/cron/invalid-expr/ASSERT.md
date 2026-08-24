## Expected

- Failure with `Error: invalid --cron:`.
- SendText not called.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertErrorPrefix(t, resp)
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "invalid --cron:", "stderr/err")
	assertNoSend(t, resp)
}
```
