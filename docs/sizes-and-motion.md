# Sizes and motion

The sizes, corner radii, fonts, shadows and motion that each Irie client draws with today. Roots does not hold these as tokens yet: each client sets its own, and this page records what the sources show so that a new client or part can match them. The one set that lives in Roots is the web parts' pair of densities in `web/components.css`. An honest account of what is missing, and the plan, is at the end.

## Density

A part keeps its anatomy and states in every client; only its sizes change ([components.md](../components.md) §1). There are two densities.

- **Pointer** (Mac, Windows): dense. Text is around 13 points, and targets run from 16 to 30 points.
- **Touch** (iPhone, phone web): each target is at least 44 points, names are 16.5 points and cards have a corner radius of 16.

In Roots' web parts the density is a switch. `data-density="touch"` on any element gives what is inside it the phone's sizes; without it, parts take the pointer clients' sizes.

| Variable | Pointer | Touch |
|---|---|---|
| `--ir-text` | 13px | 15px |
| `--ir-small` | 12px | 13px |
| `--ir-name` | 13.5px | 16.5px |
| `--ir-r` (cards, panels) | 14px | 16px |
| `--ir-r-control` (fields, controls) | 8px | 12px |
| `--ir-pad` | 8px 12px | 12px 14px |
| `--ir-target` | 28px | 44px |

The start button is 24 px round at pointer density and 34 px at touch density; the quick-action buttons are 18 px, or the whole `--ir-target` at touch density.

## Each client

| | Mac | iPhone | Phone web | Windows | Terminal |
|---|---|---|---|---|---|
| Written in | SwiftUI | SwiftUI | HTML, CSS, JavaScript | HTML, CSS, JavaScript | Go, Bubble Tea |
| Density | pointer | touch | touch | pointer | cells |
| Card radius | 14 (`Palette.radius`) | 16 | 16px (`--r`) | 8px (`--r`) | rounded box |
| Frame radius | 16 (`Palette.frameRadius`) | n/a | n/a | 12px (`--r-frame`) | n/a |
| Task bar radius | 16 | 23 | n/a | 8px | rounded box |
| Text inset | 12 (`Palette.inset`) | n/a | n/a | n/a | 1 column |
| Card name | rounded 13.5, semibold | rounded 16.5, semibold | 650 16.5px rounded | 13px base | the terminal's font |
| Base text | system | system | 17px/1.35 | 13px/1.38 | the terminal's font |

*n/a* means the sources read for this page set no such value for that client.

### Mac

- **Card.** Radius 14, frame and task bar radius 16, text inset 12, padding 8 by 12, and a 66 pt trailing slot for the age and buttons. Details panel radius 14.
- **Type.** SwiftUI's rounded design. A card's name is 13.5 semibold (medium when quiet), its detail line 12, its age 11; section heads (groups, the notes, Ready to review, Outside Irie) 11.5 heavy, tracked 0.3, in `muted`.
- **List.** Cards stand 6 apart; a head keeps 8 more above it than below; a card's helpers sit 2 under it. A group's cards step 14 in, its rail 2 wide at 9 in, under the head's chevron. The list fades into the sidebar's ground over its last 14 points.
- **Shadow.** Cards lie flat. The selected card sits on a gentle shadow: the `ink` colour at 8% opacity in light mode, black at 38% in dark; so do the segmented control's thumb and the right sidebar's picked tab, smaller. The sidebar's header and footer, the strip under the terminal, the panels and the fields have no shadow: a hairline in `edge` sets them apart. What lies over the work has its own, in black: the terminal's frame, the right sidebar as it slides in, the panes over the terminal, a notice.
- **Right sidebar.** 340 pt wide by default, 260 to 520, its tab names turning to icons when narrow.

### iPhone

- **Card.** Radius 16; the task bar and message box 23; panels 14; the needs card 20; Markdown blocks 13.
- **Type.** Names in rounded 16.5, semibold (medium when the session is quiet); group heads 14.5 to 15, bold; counts 13 with tabular figures; the state mark 18 on a card.
- **Targets.** At least 44 pt, 36 for a small button. The send button is 34 pt with a touch area 6 pt wider.

### Phone web

- **Card.** `--r` 16px; 17px/1.35 base text; names 650 16.5px.
- **Targets.** 44px for the back button, group heads and buttons.
- **Fonts.** A system stack starting `-apple-system`, a rounded stack starting `ui-rounded`, and a monospace stack starting `ui-monospace`. The colours are Evergreen's, as a copy, until the core sends the Mac app's palette.
- **Focus.** 2.5px outline in `accent`, offset 2px.

### Windows

- **Card.** `--r` 8px, `--r-frame` 12px, dialog 12px, menu 8px, notice 9px.
- **Type.** 13px/1.38 in Segoe UI Variable Text; display text in Segoe UI Variable Display; code in Cascadia Mono.
- **Targets.** Icon buttons and the find field 30px, buttons 32px, group heads 28px, menu items 30px, counts 24px high.
- **Focus.** 2px outline in `accent`, offset 1px.

### Terminal

Sizes are cells, not points, and the type is the terminal's own.

- Each card sits a column in from each side of the pane, with a column of padding.
- Where the list has fewer than 8 rows, or the pane fewer than 24 columns, each session is a one-line row.
- The header is 1 line, or 2 from 10 rows. Tabs and the boxed command bar start at 10 rows. The command bar grows from 1 to 8 lines.
- A details card replaces the status strip from 24 rows. Below 6 rows the header and command bar take one line each.
- The terminal's ground and text are `#14110D` and `#EEE5D5` in every built-in palette (`tokens.json`, `terminal`).

## Fonts in Roots' web parts

`web/components.css` sets three stacks on `.ir`. Roots does not ship font files.

| Variable | Starts with | Used for |
|---|---|---|
| `--ir-font` | `-apple-system`, `BlinkMacSystemFont`, `SF Pro Text`, `Segoe UI Variable Text` | Body text |
| `--ir-round` | `ui-rounded`, `SF Pro Rounded` | Names and group heads |
| `--ir-mono` | `ui-monospace`, `SF Mono`, `Cascadia Mono` | Code |

Line height is 1.38 at the base size. Figures that must line up, such as ages and counts, use tabular numerals.

## Shadows

Depth comes mostly from the surfaces' lightness steps ([principles.md](../principles.md) §2.5). One shadow is allowed on the selected card; no coloured shadows or glows.

| Where | Value |
|---|---|
| Roots `--ir-lift` (selected card) | `0 1px 2px` and `0 6px 18px`, each `ink` at 10% |
| Roots notice | `0 8px 24px`, `ink` at 14% |
| Phone web and Windows `--lift` | `0 10px 30px` at 16% and `0 1px 3px` at 8%, in Evergreen's ink |

## Motion

Motion only explains where a thing came from or went. Animate colour, opacity and position only; no bounce, no growing on hover, and no loop on anything that is not live ([principles.md](../principles.md) §5.1). Side panels are drawers that slide in beside the terminal, never sheets over it ([principles.md](../principles.md) §5.2).

| Motion | Mac | iPhone | Phone web | Windows | Roots web parts |
|---|---|---|---|---|---|
| Drawer or sheet | 0.28 s smooth | 0.28 s snappy (details) | 0.36 s | n/a | n/a |
| Card state or press | n/a | 0.12 s ease-out (buttons) | 0.28 s movement, 0.2 s colour | n/a | 0.15 s |
| Fold or expand | 0.15 to 0.2 s | 0.22 s | chevron 0.2 s | chevron 0.12 s | chevron 0.15 s |
| Notice in or out | n/a | spring 0.32 s in, 0.2 s out | 0.32 s in, 0.2 s out | 0.16 s out | n/a |

Easing in the web clients is `cubic-bezier(.2, .85, .25, 1)`. The Mac's panes beside the terminal use springs of 0.42 s (slide) and 0.34 s (settle). The one mark allowed to move is a live one: the phone web's "connecting" dot breathes on a 1.4 s loop, and a busy send button spins. The terminal sidebar has no motion described. Reduce Motion is in [accessibility.md](accessibility.md#reduce-motion).

### Notices

A message that comes and goes. Errors should stay at least 8 seconds, and as long as the pointer is on them ([principles.md](../principles.md) §7.3). The clients differ today:

| Client | Confirmation | With an action | Error |
|---|---|---|---|
| Mac | 3 s | 3 s | 8 s at least, held by the pointer |
| iPhone | 4.5 s | 4.5 s | 4.5 s |
| Phone web | 3.8 s | 6.5 s | 7 s |
| Windows | 4 s | 6.5 s | 8 s, held by the pointer; 2.5 s after it leaves |

On the Mac a held notice goes 1.5 s after the pointer leaves, and every notice is kept in the Notices tab.

## What is missing

There is no shared spacing scale or type ramp. Each client sets its own sizes, so the same part differs by a point or two between clients, and the half sizes (13.5, 14.5, 16.5) show that sizes were chosen by eye. The rule that sizes step on a 4-point grid and that an inner corner is the outer radius less the padding between them ([principles.md](../principles.md) §2.6) is unchecked. Radii follow the platform: 14 on the Mac, 16 on the phones, 8 on Windows.

The plan is in [README](../README.md) §5 and Irie's `docs/proposals/design-library.md`, step 4:

1. Write a spacing scale, a type ramp and corner radii per density into `tokens.json`, taken from what the clients draw now.
2. Make the clients use them: the web clients through `web/tokens.css` and `web/components.css`, the Swift apps from the library at build time.
3. Add a check, so a size changed in one place fails a test.

Until then, the values on this page are a record, not a contract. If a source here disagrees with the code, the code is right and this page is stale.
