## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "inactive", "stdout")
	assertContains(t, resp.Stdout, "Path:", "stdout")
	assertContains(t, resp.Stdout, "summary.json", "stdout")
	assertContains(t, resp.Stdout, "PIDs: none", "stdout")
}
```
