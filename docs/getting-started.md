# Getting started

Roots is Irie People's design library, named for roots reggae. It keeps how [Irie](https://github.com/irie-people/irie) looks and behaves in one place: the palettes and numbers every client draws with, the rules Irie has learned, and the parts it draws them as. Any Irie People project can draw from it.

## What Roots is

Roots has four layers ([README](../README.md) §1):

| Layer | What it is | Where it lives |
|---|---|---|
| Principles | What must hold, and why | [principles.md](../principles.md) |
| Tokens | The values: a palette's 21 roles in light and dark, and the numbers in tokens.json | `palettes/`, `tokens.json` |
| Components | The parts, built only from roles, never from colours | [components.md](../components.md), `web/components.css` |
| Patterns | How parts work together | components.md §4 |

A value lives in Roots once. Each client takes it from there.

## Who uses it

Irie has six clients, written in four languages:

| Client | Language | How it takes Roots |
|---|---|---|
| Mac app | Swift | A marked copy in Theme.swift ([swift.md](swift.md)) |
| iPhone app | Swift | A marked copy in Look.swift ([swift.md](swift.md)) |
| Phone web app | HTML, CSS, JavaScript | Marked copies in app.css, look.js and setup.html ([copies.md](copies.md)) |
| Windows app | HTML, CSS, JavaScript | Marked copies in app.css ([copies.md](copies.md)) |
| Terminal sidebar | Go | Reads the Go package ([terminal.md](terminal.md), [go.md](go.md)) |
| Site | HTML | Its themes page is checked against the library ([copies.md](copies.md)) |

## Quick start: from Go

```sh
go get github.com/irie-people/roots@main
```

```go
package main

import (
	"fmt"

	"github.com/irie-people/roots"
)

func main() {
	p, ok := roots.PaletteByID("evergreen")
	if !ok {
		panic("no such palette")
	}
	fmt.Println(p.Name, p.Dark["needs"])                       // Evergreen #FF9585
	fmt.Printf("%.2f:1\n", roots.Contrast(p.Light["ink"], p.Light["card"]))
}
```

`PaletteByID` gives a palette with a light and a dark `Set`, each a map from role name to `#RRGGBB`. `Contrast` is the WCAG 2.2 contrast ratio, from 1 to 21. The whole package is in [go.md](go.md).

## Quick start: on a web page

Link the tokens first, then the components, and wrap the page in class `ir`:

```html
<link rel="stylesheet" href="web/tokens.css">
<link rel="stylesheet" href="web/components.css">

<body class="ir" data-palette="evergreen" data-mode="dark">
  <div class="ir-card is-working">
    <span class="ir-mark" data-state="working"></span>
    <div class="ir-card-main">
      <div class="ir-card-top">
        <span class="ir-card-name">irie-tests</span>
        <span class="ir-card-slot"><span class="ir-age">4m</span></span>
      </div>
      <div class="ir-card-line">Running the tests</div>
    </div>
  </div>
</body>
```

- Without `data-palette`, the page is in the default palette (Rosewater).
- Without `data-mode`, the mode follows the system's light or dark setting.
- `data-palette` and `data-mode` are theming hooks only. Never select on them in your own CSS or script to change layout or behaviour ([web.md](web.md)).

Open `gallery.html` from a checkout to see every part in every palette and mode. It needs no server.

## Quick start: a native app that cannot read files

A Swift app, or any program that cannot read Roots where it runs, keeps a copy of the values in its own source. Mark each block with a comment, so Roots can find it:

```swift
// design: evergreen
static let evergreen = Palette(
    id: "evergreen", name: "Evergreen",
    light: Swatches(sidebar: 0xF4F4EC, ground: 0xE9EBDF /* … every role … */),
    dark: Swatches(sidebar: 0x141C17, ground: 0x0E1410 /* … every role … */))
// design: end
```

Then list the copy in a Go test and check it:

```go
var copies = []roots.Copy{{
	Path:   "App/Theme.swift",
	Why:    "The app draws its first frame without reading a file.",
	Blocks: []string{"evergreen"},
}}

func TestCopiesMatchTheLibrary(t *testing.T) {
	for _, p := range roots.Check("..", copies) {
		t.Error(p)
	}
}
```

`roots.Sync` writes the library's values into the same files. The marker rules are in [copies.md](copies.md).

## Where everything lives

| Path | What it holds |
|---|---|
| `README.md`, `AGENTS.md`, `PROVENANCE.md` | Overview, working rules, sources |
| `principles.md`, `components.md` | The rules and the parts |
| `palettes/*.json` | One file per palette: rosewater, terracotta, lilac, evergreen |
| `tokens.json` | Palette order and default, contrast floors, terminal colours, the agents' numbers |
| `roots.go`, `copies.go`, `generate.go` | The Go package `github.com/irie-people/roots` |
| `roots_test.go` | The package's tests, including the contrast check |
| `cmd/webtokens/` | The command that writes `web/tokens.css` and `web/tokens.js` |
| `web/` | `tokens.css`, `tokens.js` (generated), `components.css` |
| `gallery.html` | Every palette, role, state and part, drawn |
| `docs/` | These pages |

Next: [web.md](web.md), [go.md](go.md), [copies.md](copies.md), [contributing.md](contributing.md), and the [FAQ](faq.md).
