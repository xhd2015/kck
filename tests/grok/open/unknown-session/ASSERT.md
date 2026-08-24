## Expected

- Error: grok session not found.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertContains(t, resp.Stderr, "Error:", "stderr")
	assertContains(t, resp.Stderr, "grok session not found", "stderr")
	if len(resp.Focused) != 0 || len(resp.Opened) != 0 {
		t.Fatalf("side effects: focused=%v opened=%v", resp.Focused, resp.Opened)
	}
}
```
