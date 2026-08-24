## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "grok session not found", "stderr")
}
```
