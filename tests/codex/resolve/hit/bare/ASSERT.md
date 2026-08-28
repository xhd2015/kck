## Expected

- Exit ok; stdout is the ancestor fixture session id.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertStdoutExact(t, resp.Stdout, fixtureSessionID)
}
```
