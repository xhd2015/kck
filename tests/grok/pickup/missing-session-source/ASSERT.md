# Assert

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertErrorPrefix(t, resp)
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "expected --session-id", "error")
	if len(resp.Opened) != 0 {
		t.Fatalf("Opened=%v", resp.Opened)
	}
}
```
