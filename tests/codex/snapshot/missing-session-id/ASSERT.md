## Expected

- Usage error with `Error:`; Contents not called.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "Error:", "stderr")
	assertNoContents(t, resp)
}
```
