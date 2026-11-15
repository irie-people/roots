package main

import (
	"html"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	idAttr  = regexp.MustCompile(`\bid="([^"]+)"`)
	refAttr = regexp.MustCompile(`\b(?:href|src)="([^"]+)"`)
)

// TestEveryLinkLeadsSomewhere: the site builds, and every link in it to
// the site's own pages and files, and to a heading in one, is there.
func TestEveryLinkLeadsSomewhere(t *testing.T) {
	out := t.TempDir()
	t.Chdir(filepath.Join("..", ".."))
	if err := build(out); err != nil {
		t.Fatal(err)
	}
	pages := map[string]string{}
	ids := map[string]map[string]bool{}
	err := filepath.WalkDir(out, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".html") {
			return err
		}
		b, err := os.ReadFile(p)
		rel, _ := filepath.Rel(out, p)
		rel = filepath.ToSlash(rel)
		pages[rel] = string(b)
		ids[rel] = map[string]bool{}
		for _, m := range idAttr.FindAllStringSubmatch(string(b), -1) {
			ids[rel][m[1]] = true
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"index.html", "explorer/index.html", "principles.html", "components.html", "palettes.html", "api.html", "gallery.html"} {
		if _, ok := pages[want]; !ok {
			t.Errorf("the site has no %s", want)
		}
	}
	for page, src := range pages {
		if page == "404.html" { // its links start at the site's root on GitHub Pages
			continue
		}
		for _, m := range refAttr.FindAllStringSubmatch(src, -1) {
			ref := html.UnescapeString(m[1])
			if strings.Contains(ref, "://") || strings.HasPrefix(ref, "mailto:") || strings.HasPrefix(ref, "data:") || strings.Contains(ref, "${") {
				continue
			}
			u, err := url.Parse(ref)
			if err != nil {
				t.Errorf("%s: %q: %v", page, ref, err)
				continue
			}
			target := page
			if u.Path != "" {
				target = path.Clean(path.Join(path.Dir(page), u.Path))
			}
			if _, err := os.Stat(filepath.Join(out, filepath.FromSlash(target))); err != nil {
				t.Errorf("%s: %s leads nowhere", page, ref)
				continue
			}
			if u.Fragment != "" && strings.HasSuffix(target, ".html") && !ids[target][u.Fragment] {
				t.Errorf("%s: %s: %s has no #%s", page, ref, target, u.Fragment)
			}
		}
	}
}
