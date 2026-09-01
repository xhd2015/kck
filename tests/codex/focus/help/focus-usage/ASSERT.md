## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	for _, want := range []string{"kck codex focus", "--index", "never resumes"} {
		assertContains(t, resp.Stdout, want, "stdout")
	}
}
```
