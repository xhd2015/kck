## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "Would open new iTerm2 window", "stdout")
	assertContains(t, resp.Stdout, "resume", "stdout")
	if len(resp.Opened) != 0 {
		t.Fatalf("dry-run must not open: %v", resp.Opened)
	}
}
```
