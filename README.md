# Roots

Irie People's design library, named for roots reggae: the roots every Irie client grows from.

How [Irie](https://github.com/irie-people/irie) looks and behaves in every client, kept in one place: the palettes and numbers every client draws with, the rules Irie has learned, and the parts it draws them as. Irie has six clients: the Mac app, the iPhone app, the phone's web app, the Windows app, the terminal sidebar and its site. They are written in Swift, HTML and CSS, JavaScript and Go. This library keeps them one product, and any Irie People project can draw from it.

To see it, open [gallery.html](gallery.html) in a browser from a checkout; it needs no server.

| File | What it holds |
|---|---|
| [principles.md](principles.md) | The rules: colour, order, disclosure, controls, motion, words, accessibility, each with its reason, source and check |
| [components.md](components.md) | The parts: what each is for, its states, when to use it, and where each client draws it |
| [palettes/](palettes) | Irie's four palettes, one file each, in the form of a theme file of your own |
| [tokens.json](tokens.json) | Everything else every client must agree on: the palettes' order, the contrast floors, the terminal's colours, the agents' colours' numbers |
| [gallery.html](gallery.html) | Every palette, role, state and part, drawn in each palette and mode (`#palette=lilac&mode=dark` picks one) |
| [web/components.css](web/components.css) | The parts for web pages, as `ir-` classes built only from roles |
| web/tokens.css, web/tokens.js | The palettes for a web page, made by `go run ./cmd/webtokens` |
| roots.go, copies.go, generate.go | The library in Go: `import "github.com/irie-people/roots"` |

## 1. The layers

1. **Principles** say what must hold, and why.
2. **Tokens** are the values: a palette's 21 roles in light and dark, and the numbers in tokens.json.
3. **Components** are the parts, built only from roles, never from colours.
4. **Patterns** are how parts work together (components §4).

In Irie the core owns every rule and the clients only draw (components §1). The library is the same idea for the look: a value lives here once, and each client takes it from here.

## 2. Roles

A palette has 21 colours, each for one job. Parts name roles, never colours, so every palette, both modes and a person's own theme work everywhere.

| Group | Roles | Drawn as |
|---|---|---|
| Surfaces | `sidebar`, `ground` (the window), `card`, `panel`, `field`, `hover`, `selection`, `needs_tint` | grounds that text sits on |
| Hairline | `edge` | a line between surfaces |
| Text | `ink`, `muted`, `accent_ink`, `needs`, `done`, `paused`, `frozen` | words and glyphs: 4.5:1 on every surface |
| Accent | `accent` | a mark: 3:1 on every surface |
| Start button | `accent_fill`, with `on_accent` on it | the one prominent button |
| Icon | `sun`, `waves` | Irie's sun and sea |

`roots.Roles()` lists them with their names in Irie's Make a Theme and how each is drawn.

## 3. Changing a colour or a number

1. Change it in [palettes/](palettes) or [tokens.json](tokens.json).
2. Run `go run ./cmd/webtokens`, then `go test ./...`. The tests check every text colour at 4.5:1 and the accent at 3:1 on every surface of every palette, and that the web files are current.
3. Commit and push it ([AGENTS.md](AGENTS.md)).
4. In Irie, `go get github.com/irie-people/roots@main` takes it, `tools/design.sh sync` writes it into Irie's copies (§4), and `tools/check.sh ./internal/designcheck ./internal/ui` checks them, with Irie's tests for colour-blind viewers.

A new palette is a new file in palettes/ and its ID in tokens.json's `palettes`, then a copy in every client that lists the palettes (Irie's Theme.swift).

## 4. Copies

A Go program reads the library itself, as Irie's terminal sidebar, status line and config do. A file that cannot, such as a Swift app's palette or a web page's CSS, keeps a copy, and the library keeps the copies honest.

In a copy, each block of the library's colours sits between two comments, `design: evergreen light` and `design: end`, in the file's own comment syntax:

```css
:root {
  /* design: evergreen light page=sidebar */
  --page: #F4F4EC; --card: #FBFBF6; --needs-tint: #FBE5DC;
  /* design: end */
  --bar: #183F30;
}
```

The first comment names a palette and a mode, a palette alone (its block then says which set each colour is in with a `light` or `dark` key, as `light: Swatches(` or `"dark": {`), or `terminal`. Inside a block every name with a colour must be one of the palette's roles, in any spelling (`--needs-tint`, `needsTint`, `"needs_tint"`), or a name the first comment gives a role (`page=sidebar` above). A colour of the file's own goes outside the blocks. Numbers, such as the agents', are found by a pattern each copy gives.

A project lists its copies as `[]roots.Copy`, runs `roots.Check(root, copies)` in a test, so a copy changed by hand fails it, and `roots.Sync(root, copies)` from a command that writes the library's values in. Irie's list, and its check of its site's themes page, are in Irie's internal/designcheck, run by `tools/design.sh`.

```go
import "github.com/irie-people/roots"

p, _ := roots.PaletteByID("evergreen")
fmt.Println(p.Dark["needs"])                      // #FF9585
fmt.Println(roots.Contrast(p.Light["ink"], p.Light["card"]))
```

## 5. Where the library stands

Built: the palettes and tokens in one place, Irie's Go reading them, every copy in Irie checked and written by one command, the principles, the components, and the gallery.

Not yet, in the order Irie's docs/proposals/design-library.md plans them:

- Irie's web clients still draw their own parts: the phone app and the Windows app each have a card, a count and a group head in their own CSS. They should share web/components.css.
- Irie's site's themes page is drawn by hand; it should be made from the library, so a colour change needs no new drawing.
- The Swift apps keep their palettes as code; they could be written from the library at build time instead.
- There is no spacing scale or type ramp: each client sets its own sizes.
- The iPhone app does not honour Increase Contrast or Reduce Motion yet (principles §5.3, §7.1).
- Several parts differ between clients where they should not: they are marked **To converge** in components.md.
- A [Design Tokens](https://www.designtokens.org/tr/2025.10/) export of the palettes, for design tools.

## 6. Adding what we learn

The library is meant to grow with every lesson about design and colour, about how parts are used, and about how they are built from smaller ones:

- **A rule** goes in principles.md with its reason, source and check (principles §9).
- **A part,** or a new state of one, goes in components.md and the gallery (components §5).
- **A colour or a number** goes in palettes/ or tokens.json (§3).
- **A check** goes here, or with the code it checks. A rule with a check stays true.

Record each public source in [PROVENANCE.md](PROVENANCE.md).

## License

Apache License 2.0, as Irie's: [LICENSE](LICENSE).
