package md

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"# Roots *now*", `<h1 id="roots-now">Roots <em>now</em></h1>`},
		{"### 1.2 Text reads at 4.5:1", `<h3 id="12-text-reads-at-451">1.2 Text reads at 4.5:1</h3>`},
		{"A `<b>` and **bold** [x](a.md#y).", `<p>A <code>&lt;b&gt;</code> and <strong>bold</strong> <a href="a.html#y">x</a>.</p>`},
		{"- one\n- two\n  - deep\n- three", "<ul>\n<li>one</li>\n<li>two\n<ul>\n<li>deep</li>\n</ul>\n</li>\n<li>three</li>\n</ul>"},
		{"1. a\n2. b", "<ol>\n<li>a</li>\n<li>b</li>\n</ol>"},
		{"| A | `x|y` |\n|---|---:|\n| 1 | 2 |", `<th>A</th><th style="text-align:right"><code>x|y</code></th>`},
		{"```go\nif a < b {}\n```", `<pre><code class="language-go">if a &lt; b {}</code></pre>`},
		{"> **Note:** calm", "<blockquote>\n<p><strong>Note:</strong> calm</p>\n</blockquote>"},
		{"a*b*c and 2 * 3 * 4", "<p>a*b*c and 2 * 3 * 4</p>"},
		{"[out](https://example.com)", `<a href="https://example.com" rel="noopener">out</a>`},
	} {
		got := Render(c.in, Options{Link: func(s string) string { return strings.Replace(s, ".md", ".html", 1) }}).HTML
		if !strings.Contains(got, c.want) {
			t.Errorf("Render(%q) =\n%s\nwant it to hold\n%s", c.in, got, c.want)
		}
	}
}

// TestTheDocsRender: every Markdown file of the library renders with no
// Markdown left over outside code: no stray ** or ]( or a ``` line.
func TestTheDocsRender(t *testing.T) {
	files, _ := filepath.Glob("../../*.md")
	more, _ := filepath.Glob("../../docs/*.md")
	files = append(files, more...)
	if len(files) < 4 {
		t.Fatalf("found %d Markdown files", len(files))
	}
	codeBlock := regexp.MustCompile(`(?s)<pre>.*?</pre>|<code>.*?</code>`)
	for _, f := range files {
		if filepath.Base(f) == "CLAUDE.md" { // a pointer to AGENTS.md, not a page
			continue
		}
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		p := Render(string(src), Options{})
		if p.Title == "" {
			t.Errorf("%s has no title", f)
		}
		text := codeBlock.ReplaceAllString(p.HTML, "")
		for _, left := range []string{"**", "](", "```"} {
			if i := strings.Index(text, left); i >= 0 {
				lo, hi := max(0, i-60), min(len(text), i+60)
				t.Errorf("%s: %q left in the page: …%s…", filepath.Base(f), left, text[lo:hi])
			}
		}
	}
}

func TestSlugIsGitHubs(t *testing.T) {
	for in, want := range map[string]string{
		"4.6 Tabs act like tabs":         "46-tabs-act-like-tabs",
		"Roots' colours, light and dark": "roots-colours-light-and-dark",
		"3.1 Session card":               "31-session-card",
		"The `ir-` classes":              "the-ir--classes",
	} {
		if got := Slug(Plain(in)); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}
