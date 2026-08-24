## Expected

- No host → Error; no SendText.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertErrorPrefix(t, resp)
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "no hosting iTerm tab", "error")
	assertNoSend(t, resp)
}
```
