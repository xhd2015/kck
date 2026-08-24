## Expected

- Exit 0.
- Stdout starts with `cron every-1h` and `next[0]` lines.
- Contains would-send wording.
- SendText not called.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "cron every-1h", "stdout")
	assertContains(t, resp.Stdout, "next[0]", "stdout")
	assertContains(t, resp.Stdout, "Would send", "stdout")
	assertNoSend(t, resp)
}
```
