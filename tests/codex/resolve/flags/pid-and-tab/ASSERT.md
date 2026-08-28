## Expected

- Mutual exclusion error with `Error:` prefix.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "--pid and --tab/--tab-index cannot be specified together", "stderr")
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "Error:", "stderr")
}
```
