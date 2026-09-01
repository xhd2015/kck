# Assert

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertErrorPrefix(t, resp)
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "--here and --new-window cannot be combined", "error")
	if len(resp.Opened) != 0 || len(resp.Foreground) != 0 {
		t.Fatalf("must not launch: opened=%v fg=%v", resp.Opened, resp.Foreground)
	}
}
```
