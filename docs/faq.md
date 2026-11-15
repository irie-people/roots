# FAQ

Short answers to questions about how Roots is built and why. They come from the repository's own files and from Irie's design-library proposal (docs/proposals/design-library.md in the [Irie](https://github.com/irie-people/irie) repository).

## Why is Roots its own repository?

It began as the `design/` folder in Irie and moved out the same day, under the irie-people organisation. As a repository of its own it is a home for every Irie People project that draws, and a lesson learned outside Irie's repository can go straight into it. Irie depends on it as a Go module: a change to a palette is a change in Roots, then `go get github.com/irie-people/roots@main` and `tools/design.sh sync` in Irie. See [contributing.md](contributing.md).

## Why the name?

It is named for roots reggae: the roots every Irie client grows from. The name was chosen from names that refer to roots reggae, leaving out artists' and labels' names and names that are already a company's.

## Why not Storybook?

A Go generator builds the library's web files from its own files, with the standard library only: `go run ./cmd/webtokens` writes `web/tokens.css` and `web/tokens.js`, and `gallery.html` is one page that reads them and opens straight from a folder with no server. Storybook would need Node and many packages. Irie is a Go project that already checks its colours with Go tests, so Roots keeps one toolchain. See [web.md](web.md) and [go.md](go.md).

## Why OKLCH?

Colours of equal OKLCH lightness and chroma look equally light and equally vivid, so a palette, and the six agent colours in particular, can be spaced evenly. HSL's lightness is not what the eye sees: its yellow and blue at "50%" are far apart in lightness. Roots picks colours in OKLCH and mixes them in OKLab ([principles.md](../principles.md) §1.3, after Björn Ottosson's Oklab and CSS Color 4). Contrast floors are still measured the WCAG way, from relative luminance, with `roots.Contrast`.

## Why not APCA?

Roots measures contrast with WCAG 2.2's ratio (4.5:1 for text, 3:1 for marks, [principles.md](../principles.md) §1.2). APCA is left out for two reasons. Its [licence](https://github.com/Myndex/apca-w3) is for web content under the W3C's guidelines, with patents pending, while Roots' colours also serve Irie's native apps and terminal themes. And the [WCAG 3 draft](https://www.w3.org/TR/wcag-3.0/) of 10 September 2026 still leaves its contrast measure "to be determined". If either changes, open an issue.

## Why markers instead of templates?

A comment around each block keeps every copy's own layout, and one command checks and writes them all, with no change to how any client is built. A Swift file stays a Swift file and a CSS file stays CSS. See [copies.md](copies.md). Writing the Swift palettes whole from the library is planned ([swift.md](swift.md)).

## Why theme files rather than the Design Tokens format?

The palettes are written as a person's own Irie theme file is: a name and a light and a dark set of colours by role. Making one's own then starts by copying one. The [Design Tokens format](https://www.designtokens.org/tr/2025.10/) was weighed for tokens.json and kept for an export, so design tools can read the palettes; that export is not built ([README](../README.md) §5).

## How does Irie take a new version?

Irie pins Roots in its `go.mod`. To take a change it runs `go get github.com/irie-people/roots@main`, then `tools/design.sh sync`, which writes the library's values into every marked copy and the themes page's tables. `tools/design.sh` alone says which copies differ, and Irie's `TestCopiesMatchTheLibrary` fails the tests when a copy was changed by hand. The steps are in [contributing.md](contributing.md) and [copies.md](copies.md).

## What does "To converge" mean?

In [components.md](../components.md), a part that differs between clients where it should not carries a **To converge** note that says how. For example, a working mark is ink on the Mac and the agent's colour on the iPhone, and one should change. A part drawn in more than one client should match; when it does not, say so under **To converge** (components §5, step 5).
