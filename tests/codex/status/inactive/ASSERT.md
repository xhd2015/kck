## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "inactive", "stdout")
	assertContains(t, out, "File: no", "stdout")
	assertContains(t, out, "PIDs: none", "stdout")
	assertContains(t, out, "Path:", "stdout")
}
```
