# Colour

Irie's colour system is 21 named roles per palette, in four palettes, each with a light and a dark set. Parts name roles, never colours, so every palette, both modes and a theme of your own work everywhere. Every role has a contrast floor, the states stay apart for colour-blind viewers, and the colours that mark each working agent are worked out from a few numbers in `tokens.json`. The rules and their reasons are in [principles.md](../principles.md) §1; this page says how they are carried out.

## The 21 roles

`roots.Roles()` lists them in the palette files' order. *Group* is how Make a Theme groups the role; *use* says how it is drawn and so what it must keep.

| Role | Make a Theme calls it | Group | Use |
|---|---|---|---|
| `sidebar` | Sidebar | neutral | surface |
| `ground` | Window | neutral | surface |
| `card` | Card | neutral | surface |
| `panel` | Details panel | neutral | surface |
| `field` | Field | neutral | surface |
| `edge` | Hairline | neutral | line |
| `hover` | Card under the pointer | neutral | surface |
| `selection` | Selected card | accent | surface |
| `needs_tint` | Card that needs you | state | surface |
| `ink` | Text | neutral | text |
| `muted` | Quiet text | neutral | text |
| `accent` | Accent | accent | mark |
| `accent_ink` | Accent text | accent | text |
| `accent_fill` | Start button | accent | fill |
| `on_accent` | Arrow on the start button | accent | on_fill |
| `needs` | Needs you | state | text |
| `done` | Done | state | text |
| `paused` | Paused | state | text |
| `frozen` | Frozen | state | text |
| `sun` | Sun | sun | picture |
| `waves` | Waves | accent | picture |

### How each use is drawn

| Use | Drawn as | Roles | Floor |
|---|---|---|---|
| `surface` | A ground other colours sit on | 8 | none of its own |
| `line` | A hairline between surfaces | `edge` | none: it carries no meaning alone |
| `text` | Words, or a state's glyph | 7 | 4.5:1 on every surface |
| `mark` | A mark that is not text | `accent` | 3:1 on every surface |
| `fill` | A button's face | `accent_fill` | its own text goes on it at 4.5:1 |
| `on_fill` | The arrow on that face | `on_accent` | 4.5:1 on `accent_fill` |
| `picture` | A colour of Irie's icon, drawn large | `sun`, `waves` | none |

## The palettes

`tokens.json` lists four palettes in this order, and names `rosewater` the default. A palette is one file in `palettes/`, with a `name` and a `light` and `dark` set of all 21 roles as `#RRGGBB` in capitals.

| ID | Light accent | Dark accent |
|---|---|---|
| `rosewater` (default) | `#C4577A` | `#EE8FAD` |
| `terracotta` | `#C2603A` | `#EE9469` |
| `lilac` | `#9A5CB0` | `#C79BE0` |
| `evergreen` | `#2F5D46` | `#9CC5A1` |

In Go, `roots.IDs()`, `roots.Default()` and `roots.PaletteByID("evergreen")` read them. A palette missing a role, or holding a name that is no role, fails to load.

## OKLCH and OKLab

Colours are chosen in OKLCH (lightness, chroma, hue) and mixed or faded in OKLab, the space OKLCH is the polar form of. Roots never reasons in HSL, whose "50%" lightness is far from what the eye sees: HSL's yellow and blue at that setting differ greatly in lightness. In OKLCH, colours of equal lightness and chroma look equally light and equally vivid, so six hues can be spaced evenly and none glares. The space is Björn Ottosson's; see [principles.md](../principles.md) §1.3.

Contrast itself is WCAG 2.2's ratio of relative luminance, not a perceptual one. `roots.Contrast(a, b)` gives it, from 1 to 21.

## Contrast floors

`tokens.json` holds two numbers: `contrast.text` 4.5 and `contrast.mark` 3.0. They follow WCAG 2.2 1.4.3 and 1.4.11. `TestPalettesKeepTheirContrast` in Roots checks, for each palette and both modes:

- **Text on surfaces.** Each of the 7 text roles (`ink`, `muted`, `accent_ink`, `needs`, `done`, `paused`, `frozen`) on each of the 8 surfaces (`sidebar`, `ground`, `card`, `panel`, `field`, `hover`, `selection`, `needs_tint`): 56 pairs at 4.5:1.
- **Accent on surfaces.** `accent` on each of the 8 surfaces: 8 pairs at 3:1.
- **The start button.** `on_accent` on `accent_fill`: 1 pair at 4.5:1.

That is 65 pairs per set, 520 across the four palettes and two modes. The test fails if the counts of surfaces and text roles are not 8 and 7.

Not checked: `edge`, `sun`, `waves`, and `accent_fill` against the surfaces. Colour is measured as shipped, after any transparency, on the ground it really sits on. A ratio only just over its floor is nudged further, since rounding to 8 bits can remove the margin.

Irie checks the same floors in its own tests (`TestPalettesMeetWCAG`, `TestContrastMatchesWCAG` in internal/ui). See [accessibility.md](accessibility.md) for the other contrast floors.

## Calm, not neon

Colour goes up only where meaning does. Measured on the built-in palettes in OKLCH chroma ([principles.md](../principles.md) §1.4):

| Colours | Chroma |
|---|---|
| Surfaces | 0.003 to 0.041 |
| `selection`, `needs_tint`, text | under 0.06 |
| Agents' marks | 0.085 |
| Accents | 0.04 to 0.14 |
| States | 0.05 to 0.16, `needs` the most vivid |

These are measurements of the files, not limits the tests enforce. `TestAgentColorsAreCalm` (internal/ui) holds the agents' colours to them; the rest is unchecked. Tell things apart by hue and lightness, never by raising saturation.

## Red means one thing

`needs` marks only what waits on the person: a card, a count, a group's "· 1 needs you". A count of nothing is faint, not red. Checks: `TestRedMarksOnlyTheRowsThatNeedYou`, `TestZeroCountsAreNotRed` (internal/ui).

## Colour is never alone

Every state has a glyph and words as well as a colour: `●` needs you, `!` stalled, `✕` failed, `◐` working, `◉` waiting, `‖` paused, `◆` frozen, `✓` done, `○` idle, `◌` starting, `·` exited. See [components.md](../components.md) §2.1.

### Colour-blind floors

The four state colours must also stay apart when seen by colour-blind viewers ([principles.md](../principles.md) §1.7). Distance is deltaE-OK, in normal, protan, deutan and tritan vision (simulated as Brettel, Viénot and Mollon describe):

| Pair | Floor |
|---|---|
| Any two of `needs`, `paused`, `done`, `frozen` | 0.06, in all four visions |
| Each state against `ink` and `muted` | 0.06 in normal vision, 0.03 for dichromats |
| `done` against `accent` | 0.06 |

So the states differ in lightness as well as hue, each in its own hue family. In light mode the colours moved apart go deeper, in dark mode paler, and `needs` is deeper than `paused`. Two pairs stop short of 0.06 for dichromats as far as their hues allow: Evergreen's light `done` and `muted` (0.033), and Terracotta's dark `needs` and `muted` (0.055). The checks are `TestEveryStateHasItsGlyph`, `TestRedIsNeverTheOnlyCue` and `TestStatesStayApartForColourBlindViewers` (internal/ui).

## Dark sets

A dark set is drawn for the dark, not inverted from the light one. Its grounds are not black (OKLCH lightness 0.18 to 0.21 in Irie's palettes), surfaces are a little lighter the nearer they are, and accents are lighter than in light mode and no more vivid. The contrast tests cover both sets; the rest is unchecked.

## Washes keep to one hue

A gradient goes from a colour to a quieter step of the same hue, as a working card's wash does. Irie never draws a gradient across hues, gradient text, a glow, a coloured shadow or a neon edge ([principles.md](../principles.md) §1.10). This is unchecked.

## The agents' colours

A working session's card takes one of six colours so that cards side by side stand apart. They are worked out, not stored: each client takes the same steps with the numbers under `agents` in `tokens.json`, from the palette in use. The Go version is `roots.TheAgents()`; the Mac app's is `Theme.swift`.

| Number | Value |
|---|---|
| `slots` | 6 |
| `needs_gap` | 30 degrees |
| `mark.chroma` | 0.085 |
| `mark.lightness` | 0.50 light, 0.78 dark |
| `mark.contrast` | 3.2 |
| `ink.lightness` | 0.47 light, 0.82 dark |
| `ink.contrast` | 4.6 |
| `wash.chroma` | 0.03 light, 0.04 dark |
| `wash.lift` | -0.03 light, +0.025 dark |
| `wash.contrast` | 4.6 |

### The steps

1. **Hue.** Take the OKLCH hue of the palette's `needs` colour in the set being drawn. Leave out `needs_gap` degrees either side of it, so 300 degrees remain. Cut those into `slots` bands of 50 degrees and put each slot's hue in the middle of its band. Slot *n* (from 0) is at the `needs` hue plus 55, 105, 155, 205, 255 and 305 degrees. So `needs` red is the only red, and no agent colour is near it.
2. **Wash.** Start from the card's OKLCH lightness plus `wash.lift` (a little deeper in light mode, lifted in dark mode), at `wash.chroma`, in the slot's hue. If the card's `muted` text on it falls under `wash.contrast` (4.6:1), mix the wash back toward the plain card, as little as keeps that.
3. **Mark.** Start at `mark.lightness` and `mark.chroma`, in the slot's hue. If the mark has under `mark.contrast` (3.2:1) on the wash or on `selection`, move its lightness toward black (light mode) or white (dark mode), as little as gives it that. It is the card's edge, mark and terminal frame.
4. **Ink.** Start at `ink.lightness` and the mark's chroma. Move it the same way until it has `ink.contrast` (4.6:1) on `sidebar`, `ground`, `card` and `panel`. This is the colour of the agent's text away from its card: its helpers, the details panel and the subtitle above its terminal.

The contrast targets, 3.2 and 4.6, sit a little over the floors of 3 and 4.5 so rounding to 8 bits cannot take the margin. Each search takes 14 halvings. Under Increase Contrast the wash goes, the mark is moved to 4.5:1 and the ink to 7:1 ([accessibility.md](accessibility.md)).

A selected working card ends in a deeper step: its wash taken toward its mark by at most 0.5, and less where that would drop `ink` under 4.6:1. That 0.5 is a constant in the Mac app's `Theme.swift`, not in `tokens.json`.

### Spreading neighbours

The core, not the client, gives each session its slot (DESIGN §4.1 in Irie), so every client agrees. It does so the first time it lists a session, from the top of the list down, and the row carries it as `hue` (1 to 6). A session keeps its colour wherever it moves, through a restart, and through a close and reopen. For a card the core chooses a colour that:

- sits at least 100 degrees round the wheel from the card just above, which is two steps of the six (or the two slots either side of the `needs` gap);
- is not the colour of the card above that, nor of the nearest card above it in its own repo;
- among those, is one that the fewest sessions have, and then the one as near opposite the card above as it can.

Parked cards take theirs among the parked. The checks are `TestHuesApart`, `TestAgentColorsReadOnTheTerminal` and `TestAgentColorsAreCalm` (internal/ui), `TestCardsSideBySideTakeColoursNearOpposite` and `TestColoursHoldThroughARestart` (internal/core), and `TestCopiesMatchTheLibrary` (internal/designcheck) for the numbers in each client's copy.

The six share one lightness, so some pairs come close in colour-blind vision. They mark which session is which, never a state, and every card has its name.
