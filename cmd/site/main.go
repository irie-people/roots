// Command site builds the library's docs site from the library itself: a
// landing page; the docs, rendered from README.md, principles.md,
// components.md, docs/*.md and PROVENANCE.md; a palettes page and the Go
// API reference, made from package roots; the parts explorer; and the
// gallery. Run it from the library's root. It writes _site, which GitHub
// Pages serves (.github/workflows/pages.yml); every link in it is
// relative, so it reads the same from a folder.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"html"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/irie-people/roots"
	"github.com/irie-people/roots/internal/md"
)

const repo = "https://github.com/irie-people/roots"

var out = flag.String("out", "_site", "the folder to write the site into")

// An entry is a page in the side menu, in its order.
type entry struct {
	Group, Title string
	Src          string // the Markdown it is made from; "" for a made page
	Path         string // its place in the site
	Make         func() md.Page
	Link         bool // a page made elsewhere: the menu only links to it
}

var entries = []entry{
	{Group: "Start", Title: "Overview", Src: "README.md", Path: "overview.html"},
	{Group: "Start", Title: "Getting started", Src: "docs/getting-started.md", Path: "docs/getting-started.html"},
	{Group: "Foundations", Title: "Colour", Src: "docs/colour.md", Path: "docs/colour.html"},
	{Group: "Foundations", Title: "Palettes", Path: "palettes.html", Make: palettesPage},
	{Group: "Foundations", Title: "Accessibility", Src: "docs/accessibility.md", Path: "docs/accessibility.html"},
	{Group: "Foundations", Title: "Sizes and motion", Src: "docs/sizes-and-motion.md", Path: "docs/sizes-and-motion.html"},
	{Group: "Foundations", Title: "Theming", Src: "docs/theming.md", Path: "docs/theming.html"},
	{Group: "Rules", Title: "Principles", Src: "principles.md", Path: "principles.html"},
	{Group: "Parts", Title: "Components", Src: "components.md", Path: "components.html"},
	{Group: "Parts", Title: "Explorer", Path: "explorer/index.html", Link: true},
	{Group: "Parts", Title: "Gallery", Path: "gallery.html", Link: true},
	{Group: "Guides", Title: "Web", Src: "docs/web.md", Path: "docs/web.html"},
	{Group: "Guides", Title: "Go", Src: "docs/go.md", Path: "docs/go.html"},
	{Group: "Guides", Title: "Swift", Src: "docs/swift.md", Path: "docs/swift.html"},
	{Group: "Guides", Title: "Terminal", Src: "docs/terminal.md", Path: "docs/terminal.html"},
	{Group: "Guides", Title: "Copies", Src: "docs/copies.md", Path: "docs/copies.html"},
	{Group: "Guides", Title: "Contributing", Src: "docs/contributing.md", Path: "docs/contributing.html"},
	{Group: "Reference", Title: "Go API", Path: "api.html", Make: apiPage},
	{Group: "Reference", Title: "FAQ", Src: "docs/faq.md", Path: "docs/faq.html"},
	{Group: "Reference", Title: "Sources", Src: "PROVENANCE.md", Path: "provenance.html"},
}

// siteOf is each source's page.
var siteOf = map[string]string{}

func main() {
	flag.Parse()
	log.SetFlags(0)
	if err := build(*out); err != nil {
		log.Fatal("site: ", err)
	}
}

func build(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	// A page whose Markdown is missing leaves the menu.
	entries = slices.DeleteFunc(entries, func(e entry) bool {
		if e.Src == "" {
			return false
		}
		_, err := os.Stat(e.Src)
		return err != nil
	})
	for _, e := range entries {
		if e.Src != "" {
			siteOf[e.Src] = e.Path
		}
	}
	funcs := template.FuncMap{"list": func(s ...string) []string { return s }}
	tmpl, err := template.New("site").Funcs(funcs).ParseFiles("site/layout.html", "site/index.html", "site/explorer.html")
	if err != nil {
		return err
	}
	w := writer{dir: dir}
	for _, e := range entries {
		if e.Link {
			continue
		}
		var p md.Page
		if e.Make != nil {
			p = e.Make()
		} else {
			src, err := os.ReadFile(e.Src)
			if err != nil {
				return err
			}
			p = md.Render(string(src), md.Options{Link: linker(e.Src, e.Path)})
		}
		data := pageData{
			Title: cmpOr(p.Title, e.Title), Root: rootOf(e.Path), Path: e.Path,
			Nav: nav(e.Path), Body: template.HTML(p.HTML), TOC: toc(p.Headings),
			Palettes: roots.Palettes(),
		}
		if e.Src != "" {
			data.Source = repo + "/blob/main/" + e.Src
		}
		if err := w.page(e.Path, tmpl, "layout.html", data); err != nil {
			return err
		}
	}
	if err := w.page("index.html", tmpl, "index.html", landing()); err != nil {
		return err
	}
	if err := w.page("explorer/index.html", tmpl, "explorer.html", explorer()); err != nil {
		return err
	}
	notFound := pageData{Title: "Not here", Root: "/roots/", Nav: nav(""), Palettes: roots.Palettes(),
		Body: template.HTML(`<h1 id="not-here">Not here</h1><p>No page is at this address. The <a href="/roots/">front page</a> and the menu lead to every page.</p>`)}
	if err := w.page("404.html", tmpl, "layout.html", notFound); err != nil {
		return err
	}
	// What the pages load, copied as they are.
	for from, to := range map[string]string{
		"gallery.html": "gallery.html", "web/tokens.css": "web/tokens.css", "web/tokens.js": "web/tokens.js",
		"web/components.css": "web/components.css", "web/components.js": "web/components.js",
		"site/site.css": "assets/site.css", "site/site.js": "assets/site.js", "site/icon.svg": "assets/icon.svg",
		"site/explorer.js": "assets/explorer.js", "site/stories.js": "assets/stories.js",
	} {
		b, err := os.ReadFile(from)
		if err != nil {
			return err
		}
		if from == "gallery.html" { // its links to the docs go to their pages
			b = mdLink.ReplaceAllFunc(b, func(m []byte) []byte {
				p := mdLink.FindSubmatch(m)
				return []byte(`href="` + linker(from, to)(string(p[1])) + `"`)
			})
		}
		if err := w.file(to, b); err != nil {
			return err
		}
	}
	return w.file(".nojekyll", nil)
}

type writer struct{ dir string }

func (w writer) file(name string, b []byte) error {
	p := filepath.Join(w.dir, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o644)
}

func (w writer) page(name string, t *template.Template, which string, data any) error {
	var b bytes.Buffer
	if err := t.ExecuteTemplate(&b, which, data); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return w.file(name, b.Bytes())
}

type pageData struct {
	Title    string
	Root     string // from the page to the site's root: "" or "../"
	Path     string
	Source   string // the page's Markdown on GitHub
	Nav      template.HTML
	Body     template.HTML
	TOC      template.HTML
	Palettes []roots.Palette
}

func cmpOr(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func rootOf(p string) string { return strings.Repeat("../", strings.Count(p, "/")) }

// rel is the link from page from to the site's file to.
func rel(from, to string) string {
	r, err := filepath.Rel(filepath.Dir(from), to)
	if err != nil {
		return to
	}
	return filepath.ToSlash(r)
}

// linker writes the links of the Markdown file src, made into the page
// dst: a link to another document becomes its page, to a file the site
// copies stays a file, and to anything else goes to it on GitHub.
func linker(src, dst string) func(string) string {
	return func(target string) string {
		if strings.Contains(target, "://") || strings.HasPrefix(target, "#") || strings.HasPrefix(target, "mailto:") {
			return target
		}
		p, frag, _ := strings.Cut(target, "#")
		if p == "" {
			return target
		}
		full := path.Clean(path.Join(path.Dir(src), p))
		var url string
		switch site, ok := siteOf[full]; {
		case ok:
			url = rel(dst, site)
		case full == "gallery.html" || strings.HasPrefix(full, "web/"):
			url = rel(dst, full)
		default:
			kind := "blob"
			if fi, err := os.Stat(full); err == nil && fi.IsDir() {
				kind = "tree"
			}
			url = repo + "/" + kind + "/main/" + full
		}
		if frag != "" {
			url += "#" + frag
		}
		return url
	}
}

// nav is the side menu, the page at current marked.
func nav(current string) template.HTML {
	var b strings.Builder
	group := ""
	for _, e := range entries {
		if e.Group != group {
			if group != "" {
				b.WriteString("</ul>\n")
			}
			group = e.Group
			fmt.Fprintf(&b, "<p class=\"nav-group\">%s</p>\n<ul>\n", html.EscapeString(group))
		}
		href := rel(cmpOr(current, "index.html"), e.Path)
		if current == "404.html" || current == "" {
			href = "/roots/" + e.Path
		}
		cur := ""
		if e.Path == current {
			cur = ` aria-current="page"`
		}
		fmt.Fprintf(&b, "<li><a href=\"%s\"%s>%s</a></li>\n", html.EscapeString(href), cur, html.EscapeString(e.Title))
	}
	b.WriteString("</ul>\n")
	return template.HTML(b.String())
}

// toc is the page's own contents: its second and third level headings.
func toc(heads []md.Heading) template.HTML {
	var b strings.Builder
	n := 0
	for _, h := range heads {
		if h.Level < 2 || h.Level > 3 {
			continue
		}
		n++
		fmt.Fprintf(&b, "<li class=\"l%d\"><a href=\"#%s\">%s</a></li>\n", h.Level, h.ID, html.EscapeString(h.Text))
	}
	if n < 2 {
		return ""
	}
	return template.HTML("<p class=\"toc-title\">On this page</p>\n<ul>\n" + b.String() + "</ul>\n")
}

// palettesPage gives every palette's colours, in both modes, with the
// lowest contrast each keeps where it is drawn.
func palettesPage() md.Page {
	var b strings.Builder
	var heads []md.Heading
	b.WriteString(`<h1 id="palettes">Palettes</h1>
<p>Irie's built-in palettes, made from <a href="` + repo + `/tree/main/palettes">palettes/</a>. Each role is drawn as it is used, and text and marks give the lowest contrast they keep on the eight surfaces: text must keep 4.5:1, the accent 3:1, and the start button's arrow 4.5:1 on the button (<a href="principles.html#12-text-reads-at-451-and-marks-at-31-on-every-surface">principles §1.2</a>). The tests check every one.</p>
`)
	surfaces := roots.RolesFor(roots.Surface)
	low := func(s roots.Set, key string) float64 {
		m := 99.0
		for _, g := range surfaces {
			m = min(m, roots.Contrast(s[key], s[g.Key]))
		}
		return m
	}
	cell := func(s roots.Set, r roots.Role) string {
		c := s[r.Key]
		note := ""
		switch r.Use {
		case roots.Text, roots.Mark:
			need := roots.TextContrast()
			if r.Use == roots.Mark {
				need = roots.MarkContrast()
			}
			l := low(s, r.Key)
			class := "ok"
			if l < need {
				class = "low"
			}
			note = fmt.Sprintf(` <span class="%s">%.2f:1</span>`, class, l)
		case roots.OnFill:
			note = fmt.Sprintf(` <span class="ok">%.2f:1</span>`, roots.Contrast(c, s["accent_fill"]))
		}
		return fmt.Sprintf(`<span class="sw" style="background:%s"></span><code>%s</code>%s`, c, c, note)
	}
	for _, p := range roots.Palettes() {
		def := ""
		if p.ID == roots.Default() {
			def = " (the default)"
		}
		id := "palette-" + p.ID
		heads = append(heads, md.Heading{Level: 2, Text: p.Name, ID: id})
		fmt.Fprintf(&b, "<h2 id=\"%s\">%s%s</h2>\n", id, html.EscapeString(p.Name), def)
		b.WriteString("<div class=\"table\"><table class=\"palette\">\n<thead><tr><th>Role</th><th>Drawn as</th><th>Light</th><th>Dark</th></tr></thead>\n<tbody>\n")
		for _, r := range roots.Roles() {
			fmt.Fprintf(&b, "<tr><th scope=\"row\"><code>%s</code> <span class=\"muted\">%s</span></th><td>%s</td><td>%s</td><td>%s</td></tr>\n",
				r.Key, html.EscapeString(r.Title), strings.ReplaceAll(string(r.Use), "_", " "), cell(p.Light, r), cell(p.Dark, r))
		}
		b.WriteString("</tbody>\n</table></div>\n")
	}
	t := roots.TheTerminal()
	heads = append(heads, md.Heading{Level: 2, Text: "The terminal", ID: "terminal"})
	fmt.Fprintf(&b, `<h2 id="terminal">The terminal</h2>
<p>Every built-in palette shares the terminal's own colours: <span class="sw" style="background:%s"></span><code>%s</code> ground and <span class="sw" style="background:%s"></span><code>%s</code> ink, %.2f:1. A theme of one's own may set its own (<a href="docs/theming.html">Theming</a>).</p>
`, t.Ground, t.Ground, t.Ink, t.Ink, roots.Contrast(t.Ink, t.Ground))
	return md.Page{Title: "Palettes", HTML: b.String(), Headings: heads}
}

// apiPage is the Go package's reference, from its own doc comments.
func apiPage() md.Page {
	fset := token.NewFileSet()
	var files []*ast.File
	names, _ := filepath.Glob("*.go")
	for _, n := range names {
		if strings.HasSuffix(n, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, n, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, f)
	}
	p, err := doc.NewFromFiles(fset, files, "github.com/irie-people/roots")
	if err != nil {
		log.Fatal(err)
	}
	var b strings.Builder
	var heads []md.Heading
	decl := func(n any) string {
		var buf bytes.Buffer
		cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 4}
		if fd, ok := n.(*ast.FuncDecl); ok {
			c := *fd
			c.Body = nil
			n = &c
		}
		cfg.Fprint(&buf, fset, n)
		return "<pre><code class=\"language-go\">" + html.EscapeString(buf.String()) + "</code></pre>\n"
	}
	docHTML := func(text string) string { return string(p.HTML(text)) }
	section := func(level int, id, title string) {
		heads = append(heads, md.Heading{Level: level, Text: title, ID: id})
		fmt.Fprintf(&b, "<h%d id=\"%s\">%s</h%d>\n", level, id, html.EscapeString(title), level)
	}
	b.WriteString(`<h1 id="go-api">Go API</h1>` + "\n")
	b.WriteString(`<pre><code class="language-go">import "github.com/irie-people/roots"</code></pre>` + "\n")
	b.WriteString(docHTML(p.Doc))
	fmt.Fprintf(&b, "<p>Made from the package's own doc comments: <a href=\"%s\">the source</a>. See also <a href=\"docs/go.html\">the Go guide</a>.</p>\n", repo)
	values := func(vs []*doc.Value) {
		for _, v := range vs {
			b.WriteString(decl(v.Decl))
			b.WriteString(docHTML(v.Doc))
		}
	}
	funcs := func(fs []*doc.Func, level int, prefix string) {
		for _, f := range fs {
			name := f.Name
			if f.Recv != "" {
				name = strings.TrimPrefix(f.Recv, "*") + "." + f.Name
			}
			section(level, prefix+name, "func "+name)
			b.WriteString(decl(f.Decl))
			b.WriteString(docHTML(f.Doc))
		}
	}
	if len(p.Consts) > 0 {
		section(2, "constants", "Constants")
		values(p.Consts)
	}
	if len(p.Vars) > 0 {
		section(2, "variables", "Variables")
		values(p.Vars)
	}
	section(2, "functions", "Functions")
	funcs(p.Funcs, 3, "func-")
	section(2, "types", "Types")
	sort.Slice(p.Types, func(i, j int) bool { return p.Types[i].Name < p.Types[j].Name })
	for _, t := range p.Types {
		section(3, "type-"+t.Name, "type "+t.Name)
		b.WriteString(decl(t.Decl))
		b.WriteString(docHTML(t.Doc))
		values(t.Consts)
		values(t.Vars)
		funcs(t.Funcs, 4, "func-")
		funcs(t.Methods, 4, "method-")
	}
	return md.Page{Title: "Go API", HTML: b.String(), Headings: heads}
}

// A rule or a part, for the front page.
type item struct {
	Title, ID string
}

type group struct {
	Title string
	Items []item
}

var numbered = regexp.MustCompile(`^\d+(\.\d+)*\s+`)

// mdLink is a link to a Markdown file in a page the site copies.
var mdLink = regexp.MustCompile(`href="([^"#:]+\.md(?:#[^"]*)?)"`)

// outline is a Markdown file's sections and the headings under each.
func outline(file string, sections func(string) bool) []group {
	src, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	var out []group
	for _, h := range md.Render(string(src), md.Options{}).Headings {
		switch {
		case h.Level == 2:
			if sections(h.Text) {
				out = append(out, group{Title: numbered.ReplaceAllString(h.Text, "")})
			} else {
				out = append(out, group{}) // a section left out
			}
		case h.Level == 3 && len(out) > 0 && out[len(out)-1].Title != "":
			out[len(out)-1].Items = append(out[len(out)-1].Items, item{numbered.ReplaceAllString(h.Text, ""), h.ID})
		}
	}
	return slices.DeleteFunc(out, func(g group) bool { return g.Title == "" || len(g.Items) == 0 })
}

type landingData struct {
	pageData
	Principles []group
	Parts      []group
	Rules      int
	PartCount  int
	Roles      int
}

func landing() landingData {
	d := landingData{pageData: pageData{Title: "Roots", Root: "", Path: "index.html", Nav: nav("index.html"), Palettes: roots.Palettes()}}
	d.Principles = outline("principles.md", func(t string) bool {
		return !strings.Contains(t, "Adding a rule") && !strings.Contains(t, "never draws")
	})
	d.Parts = outline("components.md", func(t string) bool { return strings.Contains(t, "Marks") || strings.Contains(t, "Parts") })
	for _, g := range d.Principles {
		d.Rules += len(g.Items)
	}
	for _, g := range d.Parts {
		d.PartCount += len(g.Items)
	}
	d.Roles = len(roots.Roles())
	return d
}

type explorerData struct {
	pageData
	Usage template.JS
}

// explorer gives the parts explorer each part's section of components.md,
// by its title, for its Usage tab.
func explorer() explorerData {
	src, err := os.ReadFile("components.md")
	if err != nil {
		log.Fatal(err)
	}
	type section struct {
		ID   string `json:"id"`
		HTML string `json:"html"`
	}
	usage := map[string]section{}
	parts := regexp.MustCompile(`(?m)^### `).Split(string(src), -1)
	for _, sec := range parts[1:] {
		title, body, _ := strings.Cut(sec, "\n")
		if i := strings.Index(body, "\n## "); i >= 0 {
			body = body[:i]
		}
		title = strings.TrimSpace(title)
		name := numbered.ReplaceAllString(title, "")
		usage[name] = section{md.Slug(md.Plain(title)), md.Render(body, md.Options{Link: linker("components.md", "explorer/index.html")}).HTML}
	}
	b, _ := json.Marshal(usage)
	return explorerData{
		pageData: pageData{Title: "Explorer", Root: "../", Path: "explorer/index.html", Nav: nav("explorer/index.html"), Palettes: roots.Palettes()},
		Usage:    template.JS(b),
	}
}
