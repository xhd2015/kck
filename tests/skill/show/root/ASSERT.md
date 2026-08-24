## Expected

- Exit 0.
- Frontmatter `name: kck`; retrieve example `kck skill --show`.
- No `--cursor` / `--global` install plumbing in body.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	if !strings.Contains(out, "name: kck") {
		t.Fatalf("expected name: kck in stdout:\n%s", out)
	}
	if !strings.Contains(out, "kck skill --show") {
		t.Fatalf("expected retrieve example in stdout:\n%s", out)
	}
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "topic") {
		t.Fatalf("expected multi-topic index language:\n%s", out)
	}
	if strings.Contains(out, "--cursor") || strings.Contains(out, "--global") {
		t.Fatalf("root skill must not document install plumbing flags:\n%s", out)
	}
}
```
