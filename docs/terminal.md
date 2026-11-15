# The terminal sidebar

Irie's terminal sidebar is written in Go with Bubble Tea and lipgloss. Go can import Roots, so the sidebar keeps no copy: it reads the palettes and the agents' numbers from the package ([go.md](go.md)). This page follows the code in the [Irie](https://github.com/irie-people/irie) repository, in `internal/ui`.

## Reading the palettes

`internal/ui/theme.go` turns each `roots.Palette` into a `theme` of two `palette` values, light and dark. A `palette` holds the 21 roles as `color.Color` fields (`sidebar`, `ground`, `card`, `needsTint`, `accentInk` and so on); `setOf` fills them from a `roots.Set` by role key.

| In Irie | Comes from Roots |
|---|---|
| `themes` (the sidebar's palettes, by `ui.theme` name) | `roots.Palettes()` |
| `defaultTheme` | `roots.Default()` |
| `config.Themes`, the names `ui.theme` accepts | `roots.IDs()` |
| `terminalGround` | `roots.TheTerminal().Ground` |
| `agents` (steps and numbers of the agents' colours) | `roots.TheAgents()` |

A person's own theme is read by `ParseTheme` (internal/ui/custom.go). A role the file leaves out keeps Rosewater's colour in that set. A name that is not a role, or a colour that is not `#RRGGBB`, is an error: `X is no colour of a palette` or `X = "V" is not #RRGGBB`. Once the header says which palette the Mac app draws in, the sidebar draws in it too, so Irie looks the same wherever it shows.

## Light or dark

The sidebar asks the terminal for its background and picks the matching set. Until the terminal answers it uses the dark set. If it does not answer, it asks again on later ticks, at most `askBackground` (10) times, and a new window size asks again, since a client of another terminal may have attached. `setDark` swaps `m.pal` and restyles the command bar.

The sidebar paints every cell's background itself, from its palette, so contrast holds whatever the terminal's own theme is.

## Rows, cards and glyphs

- **Cards.** A session is a rounded card: its glyph, name and age, then what it is doing in plain words, then up to three helpers (`maxHelpers`), then "and N more". It sits on the `card` role. A session that needs you takes `needs_tint` and a `needs` edge; the selected card takes `selection`, with an `accent` edge while the list has the keyboard; the card under the pointer takes `hover`.
- **One-line rows.** In a pane narrower than 24 columns (`cardMinWidth`), or a list shorter than 8 lines (`cardMinList`), `compact()` draws each session as one row: mark, name, age. The same row is used for sessions outside Irie.
- **Glyphs.** Each state keeps its own glyph, so colour is never the only cue. `look` in view.go pairs a glyph with a role:

| State | Glyph | Role |
|---|---|---|
| needs you | `●` | `needs` |
| stalled | `!` | `needs` |
| failed, or exited with an error | `✕` | `needs` |
| paused, pausing | `‖` | `paused` |
| frozen | `◆` | `frozen` |
| done | `✓` | `done` (or `needs` with `ui.red_on_done`) |
| waiting on helpers | `◉` | `ink` |
| idle | `○` | `muted` |
| starting | `◌` | `muted` |
| exited cleanly | `·` | `muted` |

A frozen session under another state also shows a `◆` in `frozen` next to the main glyph, and a boosted session shows `↑` in `accent`. Glyph meanings are in [components.md](../components.md) §2.1.

## Agent colours on the status line

While a session works, Irie's status line draws its state and the rest of its line in the session's agent colour. `internal/ui/agent.go` works it out with the steps in tokens.json (`agents`):

1. Take the palette's `needs` red and its OKLCH hue.
2. Leave out `needs_gap` degrees (30) either side of it, and spread `slots` hues (6) evenly over the rest of the wheel, `hueStep` degrees apart.
3. Give each hue the mark chroma (0.085) and start at the dark lightness (0.78).
4. Move the lightness toward white until the mark keeps 3.2:1 on the card's wash and on the selection.

`AgentColor(theme, n)` returns the result as `#rrggbb`, for slot `n` from 1 to `AgentColors`, in the dark set of the palette named `theme`: the status line is drawn on the terminal's dark ground, where it reads at 4.5:1 or more. `AgentColorIn` does the same for a palette given by its colours. The core picks which slot each session takes (internal/core/hue.go) and calls these in `agentColor` (internal/core/status.go); `HuesApart` says how far apart two slots' hues are. The Mac app, the iPhone app and the phone web app work out the same colours from the same numbers ([swift.md](swift.md), [copies.md](copies.md)).

## Tests

| Test (internal/ui) | What it checks |
|---|---|
| `TestPalettesMeetWCAG` | Text 4.5:1 and the accent 3:1 on every surface |
| `TestPalettesAreWhole` | Every palette has every role |
| `TestPalettesMatchTheMacApp` | The Mac app's Theme.swift and the sidebar hold the same colours |
| `TestStatesStayApartForColourBlindViewers` | The state colours stay apart in colour-blind simulation |
| `TestLightOrDarkFollowsTheTerminal` | The set follows the terminal's background |
| `TestEveryCellIsOnThePalette` | Draws the list, hover, parked, outside, command bar and error states at several sizes on a light and a dark terminal, and checks that every cell's background is one of the palette's surfaces |

Roots' own `TestPalettesKeepTheirContrast` checks the same floors at the source ([go.md](go.md)).
