package roots

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// TestEveryPaletteIsWhole: each palette has every role in both modes, and
// the palettes are the ones tokens.json lists, in its order.
func TestEveryPaletteIsWhole(t *testing.T) {
	var ids []string
	for _, p := range Palettes() {
		ids = append(ids, p.ID)
		for _, dark := range []bool{false, true} {
			if got := len(p.Mode(dark)); got != len(roles) {
				t.Errorf("%s: %d colours in a mode, want %d", p.ID, got, len(roles))
			}
		}
	}
	if !slices.Equal(ids, IDs()) {
		t.Errorf("palettes %v, tokens.json lists %v", ids, IDs())
	}
	if _, ok := PaletteByID(Default()); !ok {
		t.Errorf("no default palette %q", Default())
	}
}

// TestPalettesKeepTheirContrast: in every palette and mode, each colour
// drawn as text keeps the text floor and the accent the mark floor on
// every surface, and the start button's arrow keeps the text floor on the
// button (principles §1.2).
func TestPalettesKeepTheirContrast(t *testing.T) {
	surfaces, texts := RolesFor(Surface), RolesFor(Text)
	if len(surfaces) != 8 || len(texts) != 7 {
		t.Fatalf("%d surfaces and %d text roles, want 8 and 7", len(surfaces), len(texts))
	}
	for _, p := range Palettes() {
		for _, dark := range []bool{false, true} {
			s := p.Mode(dark)
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			for _, g := range surfaces {
				for _, x := range texts {
					if c := Contrast(s[x.Key], s[g.Key]); c < TextContrast() {
						t.Errorf("%s %s: %s on %s is %.2f:1, want %.1f:1", p.ID, mode, x.Key, g.Key, c, TextContrast())
					}
				}
				if c := Contrast(s["accent"], s[g.Key]); c < MarkContrast() {
					t.Errorf("%s %s: accent on %s is %.2f:1, want %.1f:1", p.ID, mode, g.Key, c, MarkContrast())
				}
			}
			if c := Contrast(s["on_accent"], s["accent_fill"]); c < TextContrast() {
				t.Errorf("%s %s: on_accent on accent_fill is %.2f:1, want %.1f:1", p.ID, mode, c, TextContrast())
			}
		}
	}
}

// TestWebTokensAreCurrent: web/tokens.css and web/tokens.js are what the
// library makes now; `go run ./cmd/webtokens` writes them.
func TestWebTokensAreCurrent(t *testing.T) {
	for path, want := range WebTokens() {
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Errorf("%s is not what the library makes: go run ./cmd/webtokens", path)
		}
	}
}

// TestContrastIsWCAGs: the contrast ratio of the colours WCAG's own
// examples use.
func TestContrastIsWCAGs(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want float64
	}{{"#FFFFFF", "#000000", 21}, {"#777777", "#FFFFFF", 4.48}, {"#595959", "#FFFFFF", 7.0}} {
		if got := Contrast(c.a, c.b); got < c.want-0.01 || got > c.want+0.01 {
			t.Errorf("Contrast(%s, %s) = %.3f, want %.2f", c.a, c.b, got, c.want)
		}
	}
}

// TestMixRoundsHalfToEven: a channel exactly between two values goes to
// the even one, as Python's round() does.
func TestMixRoundsHalfToEven(t *testing.T) {
	for _, c := range []struct {
		a, b string
		t    float64
		want string
	}{{"#000000", "#FFFFFF", 0.5, "#808080"}, {"#FFFFFF", "#000000", 0.25, "#404040"}, {"#123456", "#123456", 0.3, "#123456"}} {
		if got := Mix(c.a, c.b, c.t); got != c.want {
			t.Errorf("Mix(%s, %s, %v) = %s, want %s", c.a, c.b, c.t, got, c.want)
		}
	}
}

// project is a made-up project whose files keep copies in each syntax
// the library reads: CSS, Swift with a terminal block and a number, and
// JavaScript with numbers in a list. Its values are the library's.
func project(t *testing.T) (string, []Copy, map[string]string) {
	t.Helper()
	ev, _ := PaletteByID("evergreen")
	hex := func(s Set, role string) string { return strings.TrimPrefix(s[role], "#") }
	files := map[string]string{
		"web/app.css": strings.Join([]string{
			":root {",
			"  /* design: evergreen light page=sidebar */",
			"  --page: " + ev.Light["sidebar"] + "; --needs-tint: " + ev.Light["needs_tint"] + ";",
			"  /* design: end */",
			"  --bar: #183F30;",
			"}",
			"",
		}, "\n"),
		"app/Theme.swift": strings.Join([]string{
			"// design: evergreen",
			"static let evergreen = Palette(",
			"    light: Swatches(needs: 0x" + hex(ev.Light, "needs") + ", accentInk: 0x" + hex(ev.Light, "accent_ink") + "),",
			"    dark: Swatches(needs: 0x" + hex(ev.Dark, "needs") + ", accentInk: 0x" + hex(ev.Dark, "accent_ink") + "))",
			"// design: end",
			"// design: terminal termInk=ink",
			"static let termInk = fixed(0x" + strings.TrimPrefix(TheTerminal().Ink, "#") + ")",
			"// design: end",
			"let gap = " + strings.TrimSuffix(mustValue(t, "agents.needs_gap"), ".0") + ".0, step = (360 - 2 * gap) / 6",
			"",
		}, "\n"),
		"js/look.js": strings.Join([]string{
			"// design: evergreen",
			"export const PALETTE = { light: { needsTint: 0x" + hex(ev.Light, "needs_tint") + " }, dark: { needsTint: 0x" + hex(ev.Dark, "needs_tint") + " } };",
			"// design: end",
			"const MARK_L = [" + mustValue(t, "agents.mark.lightness.light") + ", " + mustValue(t, "agents.mark.lightness.dark") + "];",
			"",
		}, "\n"),
	}
	copies := []Copy{
		{Path: "web/app.css", Blocks: []string{"evergreen light"}},
		{Path: "app/Theme.swift", Blocks: []string{"evergreen", "terminal"}, Numbers: []Number{
			{`let gap = ([\d.]+), step`, []string{"agents.needs_gap"}},
		}},
		{Path: "js/look.js", Blocks: []string{"evergreen"}, Numbers: []Number{
			{`MARK_L = \[([\d.]+), ([\d.]+)\]`, []string{"agents.mark.lightness.light", "agents.mark.lightness.dark"}},
		}},
	}
	dir := t.TempDir()
	for path, body := range files {
		write(t, filepath.Join(dir, path), body)
	}
	return dir, copies, files
}

func mustValue(t *testing.T, path string) string {
	t.Helper()
	v, ok := Value(path)
	if !ok {
		t.Fatalf("tokens.json has no %s", path)
	}
	return v
}

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestSyncPutsTheLibraryBack: copies in every syntax read as the
// library's; changed by hand, Check finds each change, and Sync puts the
// library's values back and changes nothing else, keeping a Swift
// number's decimal point.
func TestSyncPutsTheLibraryBack(t *testing.T) {
	dir, copies, files := project(t)
	if probs := Check(dir, copies); len(probs) > 0 {
		t.Fatalf("Check of copies that are the library's: %v", probs)
	}
	ev, _ := PaletteByID("evergreen")
	changed := map[string]string{
		"web/app.css":     strings.Replace(files["web/app.css"], ev.Light["needs_tint"], "#ABCDEF", 1),
		"app/Theme.swift": strings.Replace(files["app/Theme.swift"], "let gap = ", "let gap = 1", 1),
		"js/look.js":      strings.Replace(files["js/look.js"], "MARK_L = [", "MARK_L = [1", 1),
	}
	for path, body := range changed {
		if body == files[path] {
			t.Fatalf("%s did not change", path)
		}
		write(t, filepath.Join(dir, path), body)
	}
	if probs := Check(dir, copies); len(probs) != 3 {
		t.Fatalf("Check found %v, want the three changes", probs)
	}
	wrote, probs := Sync(dir, copies)
	if len(probs) > 0 || len(wrote) != 3 {
		t.Fatalf("Sync wrote %v with problems %v", wrote, probs)
	}
	for path, want := range files {
		got, err := os.ReadFile(filepath.Join(dir, path))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Errorf("%s after Sync:\n%s\nwant:\n%s", path, got, want)
		}
	}
}

// TestScanReadsABlock: a block's names in each spelling, its own names,
// its mode keys, and what it cannot read.
func TestScanReadsABlock(t *testing.T) {
	ev, _ := PaletteByID("evergreen")
	src := strings.Join([]string{
		"/* design: evergreen light page=sidebar */",
		"--page: " + ev.Light["sidebar"] + "; --needs-tint: #000000;",
		"/* design: end */",
		"// design: evergreen",
		"light: Swatches(needsTint: 0x" + ev.Light["needs_tint"][1:] + ",",
		"dark: [\"accent_ink\": 0x" + ev.Dark["accent_ink"][1:] + "]",
		"// design: end",
		"<!-- design: terminal termInk=ink -->",
		"termInk = fixed(0x" + TheTerminal().Ink[1:] + ")",
		"<!-- design: end -->",
		"",
	}, "\n")
	c := Copy{Path: "x", Blocks: []string{"evergreen light", "evergreen", "terminal"}}
	edits, probs := scan(c, []byte(src))
	if len(probs) > 0 {
		t.Fatal(probs)
	}
	var wrong []string
	for _, e := range edits {
		if e.differs() {
			wrong = append(wrong, e.what)
		}
	}
	if len(edits) != 5 || !slices.Equal(wrong, []string{"evergreen light needs_tint"}) {
		t.Errorf("%d values, %v not the library's; want 5, [evergreen light needs_tint]", len(edits), wrong)
	}

	for _, bad := range []struct{ src, want string }{
		{"/* design: evergreen light */\n--bar: #183F30;\n/* design: end */\n", "--bar is no colour"},
		{"/* design: evergreen light */\n--ink: #1E2B22;\n", "no design: end"},
		{"/* design: evergreen */\n--ink: #1E2B22;\n/* design: end */\n", "before a light or dark key"},
		{"/* design: solarized light */\n/* design: end */\n", "names no palette"},
		{"/* design: evergreen light */\n/* design: end */\n", "has no colours"},
	} {
		_, probs := scan(Copy{Path: "x", Blocks: []string{"evergreen light"}}, []byte(bad.src))
		if !slices.ContainsFunc(probs, func(p Problem) bool { return strings.Contains(p.What, bad.want) }) {
			t.Errorf("%q: problems %v, want one saying %q", bad.src, probs, bad.want)
		}
	}
}
