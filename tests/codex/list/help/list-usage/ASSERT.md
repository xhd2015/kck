## Expected

- list help documents `--json` and `--limit`.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "--json", "help")
	assertContains(t, resp.Stdout, "--limit", "help")
}
```
