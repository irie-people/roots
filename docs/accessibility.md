# Accessibility

What Irie does so that people with different eyes, hands and settings can use it, how each item is held, and where it falls short today. Contrast and colour-blind floors are checked by tests in every palette. Increase Contrast, Reduce Motion and VoiceOver are built into the Mac app; the other clients have less, and the gaps are listed at the end. Each section names the test that checks it, as [principles.md](../principles.md) does, and says so plainly where nothing checks it yet.

## Summary

| Item | Rule | Check |
|---|---|---|
| Contrast | Text 4.5:1, accent 3:1, on every surface | `TestPalettesKeepTheirContrast` (Roots); `TestPalettesMeetWCAG`, `TestContrastMatchesWCAG` (Irie, internal/ui) |
| Colour-blind safety | States apart in four visions, glyph and words on every state | `TestEveryStateHasItsGlyph`, `TestRedIsNeverTheOnlyCue`, `TestStatesStayApartForColourBlindViewers` |
| Increase Contrast | Stronger text, hairlines and agent colours | unchecked |
| Reduce Motion | Every slide, turn and fold in one step | unchecked |
| VoiceOver | One element per card, with its actions | unchecked |
| Focus | The keyboard's place always shows | unchecked |
| Notices | Wait for the person, are kept | unchecked |
| Drags | A click or a key can do the same | unchecked |
| Menus | Name what they act on | `TestTitleQuotesEachName`, `TestABatchNamesItsSessions`, `TestNoMenuKeyTakesFn` |

## Contrast

Every text role keeps 4.5:1 and `accent` 3:1 on all eight surfaces, in all four palettes and both modes (WCAG 2.2, 1.4.3 and 1.4.11). Roots' own test checks 65 pairs per set. The exact list is in [colour.md](colour.md#contrast-floors). Make a Theme applies the same floors to a theme of your own, and its terminal text on its ground at 4.5:1 ([theming.md](theming.md)). The terminal's built-in ground `#14110D` and ink `#EEE5D5` have 15.06:1.

An agent's colours are moved to their own floors: a mark 3.2:1 on its wash and on the selection, its text 4.6:1 on the four grounds it sits on, the wash such that quiet text keeps 4.6:1 ([colour.md](colour.md#the-agents-colours)).

## Colour-blind safety

Colour is never the only cue (WCAG 2.2, 1.4.1). Every state has a glyph and words of its own. The four state colours stay at least deltaE-OK 0.06 apart in normal, protan, deutan and tritan vision, and apart from `ink` and `muted` by 0.06 (0.03 for dichromats). Two known exceptions are listed in [colour.md](colour.md#colour-blind-floors). The agents' colours share one lightness, so some pairs come close; they say which session is which, never a state.

## Increase Contrast

The Mac app follows the system's Increase Contrast setting (AppKit's `accessibilityHighContrastAqua` and `accessibilityHighContrastDarkAqua` appearances). The numbers are in `Theme.swift`. A colour is moved toward another by mixing each sRGB channel, rounded.

| What | Under Increase Contrast |
|---|---|
| Text colours (`ink`, `muted`, `accent`, `accent_ink`, `on_accent`, `needs`, `done`, `paused`, `frozen`) | Dark text, with relative luminance under 0.18, moves 35% toward black. Lighter text moves 45% toward white |
| `edge` | Moves 40% toward the set's `ink` |
| Surfaces, `accent_fill`, `sun`, `waves` | Unchanged |
| A working card's wash and the selected card's deeper end | Gone: the card keeps its plain fill |
| An agent's mark | Held at 4.5:1 on the plain card and the selection (3.2:1 otherwise) |
| An agent's text | Held at 7:1 on its grounds (4.6:1 otherwise) |

This follows WCAG 2.2 1.4.6. It is unchecked by a test.

**Elsewhere.** The phone's web app has a `prefers-contrast: more` rule: a working card loses its wash, card, call and option borders take `muted`, and the small meta, age and hint text take `ink`. Windows, the iPhone app and Roots' `web/components.css` have no such rule today.

## Reduce Motion

With Reduce Motion on, every slide, turn and fold should happen in one step ([principles.md](../principles.md) §5.3). Motion only explains where something came from or went: a drawer slides in from its edge, a fold turns its chevron, nothing bounces, and nothing loops except what is live ([sizes-and-motion.md](sizes-and-motion.md)).

- **Mac.** The right sidebar's drawer (`RightSidebar.swift`) and the panes beside the terminal (`Rolodex.swift`) read `accessibilityDisplayShouldReduceMotion` and open, close and fold in one step. The shorter animations of 0.12 to 0.2 s, such as folds and hovers, have no such check in the code today.
- **Phone web.** A `prefers-reduced-motion: reduce` rule removes every animation and transition and the pressed-card scale.
- **Roots' web parts.** The same rule covers everything inside `.ir`.
- **iPhone, Windows.** No rule today.

This is unchecked by a test.

## VoiceOver and screen readers

- **Mac.** Each card is one element: its name, its state, what it needs or does, when it started ("started 3 minutes ago") and its tokens. Its actions offer what the hover buttons and the right-click menu do. A notice is announced when it appears. The sun mark is hidden from VoiceOver as decoration.
- **iPhone.** Some elements carry accessibility labels, such as "Boosted". The app has not been audited as a whole.
- **Web.** Tab rows have `role="tablist"`, `role="tab"` and `aria-selected` but not yet the arrow keys of the tabs pattern. Roots' parts use `aria-expanded` on a group head and `aria-pressed` on a segmented picker ([components.md](../components.md) §3.13, [principles.md](../principles.md) §4.6).
- **Terminal.** The sidebar's own cursor marks the place in the command bar, so a screen magnifier or reader can follow the typing.

## Focus

The keyboard's place always shows ([principles.md](../principles.md) §4.5).

- **Mac.** The selected card keeps its full selection fill. Its edge goes bolder while the list has the keyboard.
- **Web parts.** `.ir :focus-visible` draws a 2 px outline in `accent`, offset 2 px; the phone's web app uses 2.5 px offset 2 px, Windows 2 px offset 1 px. A tab's ring is drawn inside the tab so the row never clips it. The selected `.ir-card` edge goes from 1.5 px to 2.5 px with `is-keyboard`.
- **Terminal.** The selected card or row takes an accent edge while the list has the keyboard.

## Notices

A notice never covers the task bar. An error stays at least 8 seconds and as long as the pointer is on it, and every notice is kept in the Notices tab so none is lost as it goes (WCAG 2.2, 2.2.1). The Mac app does this: a confirmation goes after 3 s, an error after 8 s at least, and neither while the pointer rests on it. The other clients differ; see [sizes-and-motion.md](sizes-and-motion.md#notices).

## Drags have a click or a key

Whatever a drag does, a click or a key can do too (WCAG 2.2, 2.5.7). On the Mac:

- Reordering: Session > Move, with ⌥⌘↑ and ⌥⌘↓ to move a card up or down its group, and ⇧ to its top or bottom.
- Moving to a group or space: the card's menu.
- Dragging the panes' edge past the terminal's narrowest: View > Panes Over the Terminal.

In the terminal sidebar, `:group NAME` and `:ungroup` move a session into and out of a group. Unchecked by a test.

## Menus and keys

Every command is in the menu bar, and a menu item names its object ("Close “api-fix”"). A key with fn in it, such as Fn+Tab for the next space, is caught by the key router and named in the item's help, never shown as its shortcut, because AppKit would drop the fn. `TestNoMenuKeyTakesFn` (internal/designcheck) holds that rule.

## Target sizes

Pointer clients use compact targets (28 px for Roots' web parts, 16 to 30 points on the Mac). Touch clients give each target at least 44 points ([sizes-and-motion.md](sizes-and-motion.md#density)).

## Known gaps

- The iPhone app does not honour Increase Contrast or Reduce Motion.
- Windows has neither a `prefers-contrast` nor a `prefers-reduced-motion` rule.
- The Mac's short animations (0.12 to 0.2 s) do not check Reduce Motion.
- Notices on the phone's web app (errors 7 s) and the iPhone app (4.5 s) are shorter than the 8 s an error should stay.
- The web apps' tab rows have the tabs' roles but not their keys.
- Of the checks above, only contrast, colour-blind safety and menus have tests. Increase Contrast, Reduce Motion, VoiceOver, focus, notices and drags are unchecked.
