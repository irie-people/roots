// Package md turns the Markdown the library's docs are written in into
// HTML for its site. It reads a small set of Markdown, the one AGENTS.md
// asks docs to keep to: ATX headings, paragraphs, bullet and numbered lists
// (one level of nesting), pipe tables, fenced code, quotes, and inline
// code, bold, italic and links. Anything else is shown as plain text.
package md

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"unicode"
)

// A Heading is a heading of the page, for its table of contents.
type Heading struct {
	Level int
	Text  string // plain text
	ID    string // the anchor GitHub gives it
}

// Page is a rendered page.
type Page struct {
	Title    string // the first level-one heading, as plain text
	HTML     string
	Headings []Heading
}

// Options change how links are written.
type Options struct {
	// Link rewrites a link's target, as a link to another .md file becomes
	// its page; nil keeps every target.
	Link func(target string) string
}

// Render turns src into a page.
func Render(src string, o Options) Page {
	r := renderer{o: o, ids: map[string]int{}}
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")
	r.blocks(lines)
	return Page{Title: r.title, HTML: r.b.String(), Headings: r.heads}
}

type renderer struct {
	o     Options
	b     strings.Builder
	title string
	heads []Heading
	ids   map[string]int
}

var (
	heading  = regexp.MustCompile(`^(#{1,6})\s+(.*?)\s*#*\s*$`)
	fence    = regexp.MustCompile("^\\s*(```+|~~~+)\\s*([\\w+-]*)")
	bullet   = regexp.MustCompile(`^(\s*)[-*+]\s+(.*)$`)
	ordered  = regexp.MustCompile(`^(\s*)(\d+)[.)]\s+(.*)$`)
	tableSep = regexp.MustCompile(`^\s*\|?\s*:?-{2,}:?\s*(\|\s*:?-{2,}:?\s*)*\|?\s*$`)
)

func blank(l string) bool { return strings.TrimSpace(l) == "" }

// blocks renders a run of lines as blocks.
func (r *renderer) blocks(lines []string) {
	for i := 0; i < len(lines); {
		l := lines[i]
		switch {
		case blank(l):
			i++
		case fence.MatchString(l):
			i = r.code(lines, i)
		case heading.MatchString(l):
			m := heading.FindStringSubmatch(l)
			r.heading(len(m[1]), m[2])
			i++
		case strings.HasPrefix(strings.TrimSpace(l), ">"):
			j := i
			var inner []string
			for j < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[j]), ">") {
				t := strings.TrimPrefix(strings.TrimSpace(lines[j]), ">")
				inner = append(inner, strings.TrimPrefix(t, " "))
				j++
			}
			r.b.WriteString("<blockquote>\n")
			r.blocks(inner)
			r.b.WriteString("</blockquote>\n")
			i = j
		case strings.Contains(l, "|") && i+1 < len(lines) && tableSep.MatchString(lines[i+1]):
			i = r.table(lines, i)
		case bullet.MatchString(l) || ordered.MatchString(l):
			i = r.list(lines, i)
		default:
			j := i
			var para []string
			for j < len(lines) && !blank(lines[j]) && !r.starts(lines, j) {
				para = append(para, strings.TrimSpace(lines[j]))
				j++
			}
			if len(para) == 0 { // a line no block takes
				para, j = []string{strings.TrimSpace(l)}, i+1
			}
			fmt.Fprintf(&r.b, "<p>%s</p>\n", r.inline(strings.Join(para, " ")))
			i = j
		}
	}
}

// starts says whether line i begins a block other than a paragraph.
func (r *renderer) starts(lines []string, i int) bool {
	l := lines[i]
	return fence.MatchString(l) || heading.MatchString(l) || strings.HasPrefix(strings.TrimSpace(l), ">") ||
		bullet.MatchString(l) || ordered.MatchString(l) ||
		strings.Contains(l, "|") && i+1 < len(lines) && tableSep.MatchString(lines[i+1])
}

func (r *renderer) heading(level int, text string) {
	plain := Plain(text)
	id := Slug(plain)
	if n := r.ids[id]; n > 0 {
		r.ids[id]++
		id = fmt.Sprintf("%s-%d", id, n)
	} else {
		r.ids[id] = 1
	}
	if level == 1 && r.title == "" {
		r.title = plain
	}
	r.heads = append(r.heads, Heading{level, plain, id})
	fmt.Fprintf(&r.b, "<h%d id=\"%s\">%s</h%d>\n", level, id, r.inline(text), level)
}

func (r *renderer) code(lines []string, i int) int {
	m := fence.FindStringSubmatch(lines[i])
	mark, lang := m[1], m[2]
	j := i + 1
	var body []string
	for j < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[j]), mark[:3]) {
		body = append(body, lines[j])
		j++
	}
	class := ""
	if lang != "" {
		class = fmt.Sprintf(` class="language-%s"`, html.EscapeString(lang))
	}
	fmt.Fprintf(&r.b, "<pre><code%s>%s</code></pre>\n", class, html.EscapeString(strings.Join(body, "\n")))
	return j + 1
}

// cells splits a table row on pipes outside code spans.
func cells(row string) []string {
	row = strings.TrimSpace(row)
	row = strings.TrimPrefix(row, "|")
	if strings.HasSuffix(row, "|") && !strings.HasSuffix(row, `\|`) {
		row = row[:len(row)-1]
	}
	var out []string
	var cur strings.Builder
	code := false
	for i := 0; i < len(row); i++ {
		c := row[i]
		switch {
		case c == '\\' && i+1 < len(row) && row[i+1] == '|':
			cur.WriteByte('|')
			i++
		case c == '`':
			code = !code
			cur.WriteByte(c)
		case c == '|' && !code:
			out = append(out, strings.TrimSpace(cur.String()))
			cur.Reset()
		default:
			cur.WriteByte(c)
		}
	}
	return append(out, strings.TrimSpace(cur.String()))
}

func (r *renderer) table(lines []string, i int) int {
	head := cells(lines[i])
	aligns := cells(lines[i+1])
	align := func(k int) string {
		if k >= len(aligns) {
			return ""
		}
		a := aligns[k]
		switch {
		case strings.HasPrefix(a, ":") && strings.HasSuffix(a, ":"):
			return ` style="text-align:center"`
		case strings.HasSuffix(a, ":"):
			return ` style="text-align:right"`
		}
		return ""
	}
	r.b.WriteString("<div class=\"table\"><table>\n<thead><tr>")
	for k, c := range head {
		fmt.Fprintf(&r.b, "<th%s>%s</th>", align(k), r.inline(c))
	}
	r.b.WriteString("</tr></thead>\n<tbody>\n")
	j := i + 2
	for j < len(lines) && strings.Contains(lines[j], "|") && !blank(lines[j]) {
		r.b.WriteString("<tr>")
		for k, c := range cells(lines[j]) {
			if k >= len(head) {
				break
			}
			fmt.Fprintf(&r.b, "<td%s>%s</td>", align(k), r.inline(c))
		}
		r.b.WriteString("</tr>\n")
		j++
	}
	r.b.WriteString("</tbody>\n</table></div>\n")
	return j
}

// list renders a list and the lists nested in its items.
func (r *renderer) list(lines []string, i int) int {
	indent := func(l string) int { return len(l) - len(strings.TrimLeft(l, " ")) }
	base := indent(lines[i])
	tag := "ul"
	if ordered.MatchString(lines[i]) && indent(lines[i]) == base && !bullet.MatchString(lines[i]) {
		tag = "ol"
	}
	r.b.WriteString("<" + tag + ">\n")
	j := i
	for j < len(lines) {
		l := lines[j]
		if blank(l) {
			// A blank line ends the list unless an item at this level follows.
			k := j + 1
			for k < len(lines) && blank(lines[k]) {
				k++
			}
			if k < len(lines) && indent(lines[k]) >= base && (bullet.MatchString(lines[k]) || ordered.MatchString(lines[k]) || indent(lines[k]) > base) {
				j = k
				continue
			}
			break
		}
		if indent(l) < base || !(bullet.MatchString(l) || ordered.MatchString(l)) || indent(l) > base {
			break
		}
		var first string
		if m := bullet.FindStringSubmatch(l); m != nil {
			first = m[2]
		} else {
			first = ordered.FindStringSubmatch(l)[3]
		}
		item := []string{first}
		j++
		for j < len(lines) && !blank(lines[j]) && indent(lines[j]) > base {
			item = append(item, lines[j])
			j++
		}
		r.b.WriteString("<li>")
		r.item(item, base)
		r.b.WriteString("</li>\n")
	}
	r.b.WriteString("</" + tag + ">\n")
	return j
}

// item renders a list item: its text, then any list nested in it.
func (r *renderer) item(lines []string, base int) {
	var text []string
	k := 0
	for k < len(lines) {
		l := lines[k]
		if k > 0 && (bullet.MatchString(l) || ordered.MatchString(l)) {
			break
		}
		text = append(text, strings.TrimSpace(l))
		k++
	}
	r.b.WriteString(r.inline(strings.Join(text, " ")))
	if k < len(lines) {
		r.b.WriteString("\n")
		r.list(lines, k)
	}
}

var (
	codeSpan = regexp.MustCompile("(`+)(.+?)(`+)")
	link     = regexp.MustCompile(`\[([^\]]+)\]\(([^)\s]+)\)`)
	strong   = regexp.MustCompile(`\*\*([^*]+?)\*\*`)
	em       = regexp.MustCompile(`(^|[^\w*])\*([^*\s][^*]*?)\*([^\w*]|$)`)
)

// inline renders a line's inline Markdown, escaping everything else.
func (r *renderer) inline(s string) string {
	// Code spans first, so nothing inside them is read as Markdown.
	var codes []string
	s = codeSpan.ReplaceAllStringFunc(s, func(m string) string {
		p := codeSpan.FindStringSubmatch(m)
		if len(p[1]) != len(p[3]) {
			return m
		}
		codes = append(codes, "<code>"+html.EscapeString(strings.TrimSpace(p[2]))+"</code>")
		return fmt.Sprintf("\x00%d\x00", len(codes)-1)
	})
	s = html.EscapeString(s)
	s = link.ReplaceAllStringFunc(s, func(m string) string {
		p := link.FindStringSubmatch(m)
		target := html.UnescapeString(p[2])
		if r.o.Link != nil {
			target = r.o.Link(target)
		}
		rel := ""
		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			rel = ` rel="noopener"`
		}
		return fmt.Sprintf(`<a href="%s"%s>%s</a>`, html.EscapeString(target), rel, p[1])
	})
	s = strong.ReplaceAllString(s, "<strong>$1</strong>")
	s = em.ReplaceAllString(s, "$1<em>$2</em>$3")
	for i, c := range codes {
		s = strings.Replace(s, fmt.Sprintf("\x00%d\x00", i), c, 1)
	}
	return s
}

// Plain is a line's text with its Markdown taken out.
func Plain(s string) string {
	s = codeSpan.ReplaceAllString(s, "$2")
	s = link.ReplaceAllString(s, "$1")
	s = strings.NewReplacer("**", "", "*", "").Replace(s)
	return strings.TrimSpace(s)
}

// Slug is a heading's anchor as GitHub makes it: lower case, letters,
// digits, hyphens and underscores kept, spaces made hyphens.
func Slug(s string) string {
	var b strings.Builder
	for _, c := range strings.ToLower(s) {
		switch {
		case unicode.IsLetter(c) || unicode.IsDigit(c) || c == '-' || c == '_':
			b.WriteRune(c)
		case c == ' ':
			b.WriteByte('-')
		}
	}
	return b.String()
}
