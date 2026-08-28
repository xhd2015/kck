## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	for _, want := range []string{"kck codex info", "Active", "--no-pid"} {
		assertContains(t, resp.Stdout, want, "stdout")
	}
}
```
