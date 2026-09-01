# Assert

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "name: kck-pickup-a-session", "stdout")
	assertContains(t, resp.Stdout, "kck grok messages", "stdout")
	assertContains(t, resp.Stdout, "kck codex messages", "stdout")
}
```
