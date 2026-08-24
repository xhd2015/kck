## Expected

- Exit 0; frontmatter `name: kck/send`; body mentions `--open`.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if !strings.Contains(resp.Stdout, "name: kck/send") {
		t.Fatalf("expected name: kck/send:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "--open") {
		t.Fatalf("expected --open in send topic:\n%s", resp.Stdout)
	}
}
```
