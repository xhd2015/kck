## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	msg := resp.ErrText + resp.Stderr
	assertContains(t, msg, "Error:", "error")
	assertContains(t, msg, "not found", "error")
	if len(resp.Focused) != 0 {
		t.Fatalf("Focused=%v", resp.Focused)
	}
}
```
