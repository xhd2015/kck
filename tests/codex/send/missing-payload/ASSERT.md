## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	combined := resp.ErrText + resp.Stderr
	assertContains(t, combined, "Error:", "err")
	assertContains(t, combined, "missing text or key", "err")
}
```
