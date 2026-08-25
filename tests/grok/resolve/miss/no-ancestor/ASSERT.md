## Expected

- Hard failure with `Error: no ancestor grok`.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "no ancestor grok", "stderr")
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "Error:", "stderr")
}
```
