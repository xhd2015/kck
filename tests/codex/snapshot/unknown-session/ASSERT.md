## Expected

- Hard `Error: codex session not found`.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "codex session not found", "stderr")
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "Error:", "stderr")
	assertNoContents(t, resp)
}
```
