## Expected

- No host → Error: no hosting iTerm tab; Contents not called.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "no hosting iTerm tab", "stderr")
	assertNoContents(t, resp)
}
```
