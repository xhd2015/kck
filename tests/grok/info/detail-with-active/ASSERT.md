## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "Session: "+req.SessionID, "stdout")
	assertContains(t, out, "kck info fixture", "stdout")
	assertContains(t, out, "Files:", "stdout")
	assertContains(t, out, "Active:", "stdout")
	assertContains(t, out, "5001", "stdout")
}
```
