## Expected

- Exit 0.
- Exactly two sends (18:28 and 18:43); quiet slots skipped.
- `cron done: until reached`.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.SendCalls) != 2 {
		t.Fatalf("SendCalls=%d, want 2 (quiet skipped); stdout:\n%s", len(resp.SendCalls), resp.Stdout)
	}
	assertContains(t, resp.Stdout, "cron done: until reached", "stdout")
}
```
