# The Go package

`github.com/irie-people/roots` is the library in code. It embeds `tokens.json` and `palettes/*.json`, so a Go program reads the palettes and numbers with no files to find at run time. It uses the standard library only. Take it with `go get github.com/irie-people/roots@main`.

## Errors

The package loads its files when it starts, and panics with a message that begins `roots:` if one is wrong. Examples from the loader: `tokens.json: the default palette "x" is not in palettes`, `palettes/ has 3 files, tokens.json lists 4 palettes`, `palettes/lilac.json: dark has no ink`, and a colour that is not `#RRGGBB` in capitals. Unknown fields in either file are refused.

So a program never handles a bad library at run time. The package's own tests load it first, so a broken file fails `go test ./...` in the library before it can reach a client.

## Roles

A palette has 21 colours, each a **role** ([README](../README.md) §2).

| Identifier | What it is |
|---|---|
| `type Role` | `Key` (name in a theme file, `needs_tint`), `Title` (name in Irie's Make a Theme), `Group`, `Use` |
| `type Group` | `Neutral`, `Accent`, `State`, `Sun`: how Make a Theme groups a role |
| `type Use` | How a role is drawn, which says what contrast it keeps |
| `Roles() []Role` | All 21, in the palette files' order (a copy) |
| `RolesFor(use Use) []Role` | The roles drawn as `use` |

The `Use` values:

| Const | Meaning | Floor |
|---|---|---|
| `Surface` | A ground other colours are drawn on | None |
| `Line` | A hairline between surfaces | None |
| `Text` | Words or a state's glyph | 4.5:1 on every surface |
| `Mark` | A mark that is not text | 3:1 on every surface |
| `Fill` | A button's face | Its text keeps 4.5:1 on it |
| `OnFill` | The text or arrow on a `Fill` | 4.5:1 on the fill |
| `Picture` | A colour of Irie's icon, drawn large | None |

## Palettes and tokens

| Identifier | What it does |
|---|---|
| `type Set` | One mode of a palette: `map[string]string` of `#RRGGBB` by role key |
| `type Palette` | `ID` (`"evergreen"`), `Name` (`"Evergreen"`), `Light`, `Dark` |
| `Palette.Mode(dark bool) Set` | The dark set, or the light one |
| `Palettes() []Palette` | The built-in palettes, in the menus' order |
| `IDs() []string` | Their IDs in the same order: `ui.theme`'s values in Irie |
| `Default() string` | The palette drawn until one is chosen (`rosewater`) |
| `PaletteByID(id) (Palette, bool)` | One palette by ID |
| `TextContrast()`, `MarkContrast()` | The floors from tokens.json: 4.5 and 3.0 |
| `TheTerminal() Terminal` | `Ground` and `Ink` of the terminal, the same in every built-in palette |
| `TheAgents() Agents` | The numbers of the agents' colours |
| `Value(path) (string, bool)` | One number or colour of tokens.json by dotted path, as written there |

`Value("agents.mark.chroma")` returns `"0.085"`. Because it gives the text as written, a copy that writes `0.50` still equals `0.5` ([copies.md](copies.md)).

`Agents` mirrors the `agents` object: `Slots`, `NeedsGap`, and `Mark` (`Chroma`, `Lightness`, `Contrast`), `Ink` (`Lightness`, `Contrast`) and `Wash` (`Chroma`, `Lift`, `Contrast`). A `Pair` has `Light` and `Dark` numbers, and `Pair.Of(dark bool)` picks one. Irie's terminal sidebar works the agents' colours out from these in internal/ui/agent.go ([terminal.md](terminal.md)).

```go
a := roots.TheAgents()
fmt.Println(a.Slots, a.Mark.Lightness.Of(true)) // 6 0.78
```

## Colour functions

| Function | What it does |
|---|---|
| `RGB(hex) (r, g, b uint8, err error)` | Channels of `#RRGGBB`; also accepts no `#` or a `0x` prefix |
| `Luminance(hex) float64` | WCAG 2.2 relative luminance; `NaN` for a bad colour |
| `Contrast(a, b string) float64` | WCAG 2.2 contrast ratio, from 1 to 21 |
| `Mix(a, b string, t float64) string` | `t` of `a` over `b` in sRGB, as `#RRGGBB`, each channel rounded half to even as Python's `round()` does |

```go
p, _ := roots.PaletteByID("lilac")
for _, r := range roots.RolesFor(roots.Text) {
	c := roots.Contrast(p.Dark[r.Key], p.Dark["card"])
	fmt.Printf("%-10s %.2f:1\n", r.Key, c)
}
```

`Mix` rounds as it does so Go reproduces numbers the pictures on Irie's site were first made with.

## Copies

`Copy`, `Number`, `Problem`, `Check` and `Sync` check and write the copies of the library's values in other files. They are described in [copies.md](copies.md).

## Web files

`WebTokens() map[string][]byte` gives the contents of `web/tokens.css` and `web/tokens.js`, keyed by those paths. Run the command in the library's root after changing `palettes/` or `tokens.json`:

```sh
go run ./cmd/webtokens
```

It prints `wrote web/tokens.css` and `wrote web/tokens.js`. The CSS is described in [web.md](web.md).

## Tests

Run `go test ./...`. The tests in `roots_test.go`:

| Test | What it checks |
|---|---|
| `TestEveryPaletteIsWhole` | Every palette has every role in both modes, in the order tokens.json lists |
| `TestPalettesKeepTheirContrast` | Text roles keep 4.5:1 and the accent 3:1 on all eight surfaces, and `on_accent` keeps 4.5:1 on `accent_fill`, in every palette and mode |
| `TestWebTokensAreCurrent` | The web files equal what `WebTokens` makes |
| `TestContrastIsWCAGs` | `Contrast` gives WCAG's own example ratios |
| `TestMixRoundsHalfToEven` | `Mix`'s rounding |
| `TestSyncPutsTheLibraryBack` | `Check` finds hand changes in CSS, Swift and JavaScript copies; `Sync` puts them back and changes nothing else |
| `TestScanReadsABlock` | Name spellings, aliases, mode keys and the messages for blocks it cannot read |

CI runs `gofmt`, `go vet ./...` and `go test ./...` on every push to main and every pull request.
