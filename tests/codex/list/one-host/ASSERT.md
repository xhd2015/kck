## Expected

- Table contains session id and iterm cell.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, req.SessionID, "stdout")
	assertContains(t, resp.Stdout, "w=3 t=1", "stdout")
	assertContains(t, resp.Stdout, "1 sessions", "stdout")
}
```
