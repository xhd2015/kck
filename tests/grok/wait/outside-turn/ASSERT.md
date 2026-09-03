## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "reason: turn_completed", "stdout")
	assertContains(t, resp.Stdout, "session-id: "+req.SessionID, "stdout")
}
```
