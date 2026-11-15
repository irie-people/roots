# Theming

A theme of your own is one small JSON file in the same shape as Roots' built-in palettes. This page gives the file format as Irie reads it, how Make a Theme builds one from three picks, how the built-in palettes in `palettes/` relate to it, how web pages choose a palette, and how to check a theme's contrast. The roles themselves are in [colour.md](colour.md).

## The theme file

Irie reads `$XDG_CONFIG_HOME/irie/themes/<id>.json`, or `~/.config/irie/themes/<id>.json` when `XDG_CONFIG_HOME` is not an absolute path. The file's name without `.json` is the theme's ID.

```json
{
  "name": "Harbour",
  "based_on": "evergreen",
  "light": {
    "accent": "#2F5D46",
    "sun": "#D9A93F"
  },
  "dark": {
    "accent": "#9CC5A1"
  },
  "terminal": {
    "ground": "#14110D",
    "ink": "#EEE5D5"
  }
}
```

| Key | Holds | If left out |
|---|---|---|
| `name` | The name the menus show | The file's ID |
| `based_on` | The ID of a built-in palette | `rosewater` (also used when the ID is not a built-in) |
| `light`, `dark` | Role names mapped to `#RRGGBB` | The base palette's whole set |
| `terminal` | `ground` and `ink`, the terminal's own colours, the same in both sets | The built-in `#14110D` and `#EEE5D5`; one given and one not takes the built-in for the other |

Every key is optional, so a hand-made file can give only a name and a colour or two.

### Roles

A set may use any of the 21 role names in [colour.md](colour.md#the-21-roles): `sidebar`, `ground`, `card`, `panel`, `field`, `edge`, `hover`, `selection`, `needs_tint`, `ink`, `muted`, `accent`, `accent_ink`, `accent_fill`, `on_accent`, `needs`, `done`, `paused`, `frozen`, `sun` and `waves`. A colour is `#RRGGBB`; the leading `#` may be left off, and either case is read. A role a file leaves out keeps the colour of the palette named by `based_on`, in that mode. Make a Theme writes all 21, light then dark, in the order above, then `terminal`.

### Validation

Irie lists every file in the folder, oldest first, and reports what is wrong with a file by its name. A bad entry does not stop the file's other colours.

| Problem | What happens |
|---|---|
| The ID is a built-in palette's (`rosewater`, `terracotta`, `lilac`, `evergreen`) | The file is left out |
| The file is not a theme's JSON | The file is left out |
| A key in `light` or `dark` is no role | That entry is ignored and named as a problem |
| A colour is not `#RRGGBB` | That entry is ignored and named as a problem |
| A `terminal` colour is not `#RRGGBB` | That colour is ignored and named as a problem |

Irie does not refuse a theme for poor contrast. It reads the folder as the app starts, when a file comes, goes or is renamed, and whenever the app becomes active, so a file an agent writes shows in the menus at once. Yours are listed after the built-in palettes. The Mac app tells the core the theme in use, with both its sets, and the core passes it to the terminal sidebar, the terminal slot and the phone app, which then draw in it.

## Make a Theme

Open View > Theme > Make a Theme, or Settings > Appearance. It starts from the palette in use, or any other you pick, and takes three colours:

1. **Accent.** The colour of the accent roles.
2. **Tint.** The colour the grounds, cards and text lean toward.
3. **Sun.** The sun of Irie's mark.

Irie works out the rest of both sets in OKLCH. For each role:

- **Grounds, hairline and text** (`sidebar`, `ground`, `card`, `panel`, `field`, `edge`, `hover`, `ink`, `muted`). The hue turns round the wheel by as much as the tint turned from where it began, and the chroma grows or shrinks by as much, up to 4 times, but never past a cap: 0.05, or 0.09 for `ink` and `muted`. Lightness is kept.
- **Accent roles** (`selection`, `accent`, `accent_ink`, `accent_fill`, `on_accent`, `waves`). The same, following the accent pick, with caps of 0.17 and 0.08 for `selection`. In light mode `accent` is the pick itself.
- **`sun`.** The pick.
- **States** (`needs`, `needs_tint`, `done`, `paused`, `frozen`). Unchanged, so they keep their meaning.

A colour that is nearly grey (chroma under 0.006) takes the pick's hue directly. Because each colour keeps the lightness the base palette gave it, text reads as well as before.

### Moved to the floors

After the picks are applied, any text or accent colour that would miss its contrast moves in lightness, deeper on light grounds and lighter on dark ones, keeping its hue and chroma, only as far as gives it 4.5:1 for text and 3:1 for the accent on every ground. It is searched to a margin of 0.1 over the floor. A colour you set by hand is kept until its arrow button hands it back.

**Every colour** shows all 21 of a set, light or dark, to pick or type. A line below says that every colour reads, or which do not and on which ground, with a **Make Them Readable** button that moves them the same way. The terminal's colours turn with the tint at a chroma of at most 0.03 and the same lightness, so it stays dark, and its text is held to 4.5:1 on its ground. **Save and Use** writes the file and makes it the theme in use.

## The built-in palettes are themes

Each file in `palettes/` is a theme file with every role and no `based_on` or `terminal`: its `name`, a `light` set and a `dark` set of all 21 roles. The terminal's colours are in `tokens.json` for all four. Roots is stricter than Irie's reader:

- Unknown keys are refused.
- All 21 roles must be present in both sets.
- A colour must be `#RRGGBB` in capitals.
- A palette needs a name.

So a theme you like becomes a built-in by copying it to `palettes/<id>.json` with every role written out, removing `based_on` and `terminal`, adding the ID to `palettes` in `tokens.json`, then running `go run ./cmd/webtokens` and `go test ./...`. The contrast test then holds it to the floors. Each client that lists the palettes needs a copy too ([README](../README.md) §3). See [copies.md](copies.md).

## Web pages

Link `web/tokens.css`, which is made by `go run ./cmd/webtokens` from `palettes/` and `tokens.json`. Each role is a custom property with underscores turned to dashes: `--needs-tint`, `--accent-ink`, `--on-accent`. The terminal's colours are `--term-ground` and `--term-ink`.

- With no attributes, a page draws in the default palette (`rosewater`), in the system's mode (`prefers-color-scheme`).
- `data-palette="evergreen"` on an element draws that element and what is inside it in that palette.
- `data-mode="light"` or `data-mode="dark"` on the same element picks its mode instead of the system's.

```html
<section data-palette="lilac" data-mode="dark">...</section>
```

**These attributes are theming hooks only.** They say how an element is drawn. A page keeps its own state in its own attributes, as the gallery does with `data-pick-*`, so that nothing depends on them for anything but colour. `web/tokens.js` holds the same data for scripts as `window.IRIE_DESIGN`: the roles, the palettes, the default, the contrast floors, the terminal and the agents' numbers. `web/components.css` builds the parts from these variables alone ([web.md](web.md)). Roots makes CSS only for its four palettes; a theme of your own is drawn by setting the same variables on an element.

## Checking a theme's contrast

Make a Theme checks live. To check a file from outside it, apply the floors of [colour.md](colour.md#contrast-floors) with the Go library: 8 surfaces by 7 text roles at 4.5:1, the accent at 3:1, and `on_accent` on `accent_fill` at 4.5:1. This program fills missing roles from `based_on`, as Irie does:

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/irie-people/roots"
)

func main() {
	raw, _ := os.ReadFile(os.Args[1])
	var t struct {
		BasedOn     string `json:"based_on"`
		Light, Dark map[string]string
	}
	json.Unmarshal(raw, &t)
	base, ok := roots.PaletteByID(t.BasedOn)
	if !ok {
		base, _ = roots.PaletteByID(roots.Default())
	}
	for _, dark := range []bool{false, true} {
		set := roots.Set{}
		for k, v := range base.Mode(dark) {
			set[k] = v
		}
		given := map[bool]map[string]string{false: t.Light, true: t.Dark}[dark]
		for k, v := range given {
			set[k] = v
		}
		check := func(fg, bg string, want float64) {
			if c := roots.Contrast(set[fg], set[bg]); c < want {
				fmt.Printf("dark=%v %s on %s is %.2f:1, want %.1f:1\n", dark, fg, bg, c, want)
			}
		}
		for _, g := range roots.RolesFor(roots.Surface) {
			for _, x := range roots.RolesFor(roots.Text) {
				check(x.Key, g.Key, roots.TextContrast())
			}
			check("accent", g.Key, roots.MarkContrast())
		}
		check("on_accent", "accent_fill", roots.TextContrast())
	}
}
```

Make a Theme checks the 8 surfaces by `ink`, `muted`, `needs`, `done`, `paused`, `frozen` and `accent_ink`, the accent, and the terminal's text on its ground, but not `on_accent` on `accent_fill`, which the program above and Roots' own test do check. The colour-blind floors for the four states ([colour.md](colour.md#colour-blind-floors)) are checked by Irie's tests for the built-in palettes only; a theme that changes the state colours should be judged against them by eye or by its own test.
