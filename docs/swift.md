# Swift apps

Roots has no Swift package. Irie's Mac and iPhone apps keep their palettes as Swift code, between marker comments that the Go package checks ([copies.md](copies.md)). This page says what is where today, which parts the apps draw in SwiftUI, and what is planned.

## Palettes as code

Paths below are in the [Irie](https://github.com/irie-people/irie) repository.

### Mac: mac/Sources/Irie/Theme.swift

`Palette` holds an `id`, a `name`, and two `Swatches` (light and dark). `Swatches` has one `UInt32` field per role, in camel case: `sidebar`, `ground`, `card`, `panel`, `field`, `edge`, `hover`, `selection`, `needsTint`, `ink`, `muted`, `accent`, `accentInk`, `accentFill`, `onAccent`, `needs`, `done`, `paused`, `frozen`, `sun`, `waves`.

Each of the four palettes is one block:

```swift
// design: evergreen
static let evergreen = Palette(
    id: "evergreen", name: "Evergreen",
    light: Swatches(
        sidebar: 0xF4F4EC, ground: 0xE9EBDF, card: 0xFBFBF6, /* … */),
    dark: Swatches(
        sidebar: 0x141C17, ground: 0x0E1410, card: 0x1C2620, /* … */))
// design: end
```

The `light:` and `dark:` keys tell the check which set each colour is in. A fifth block, `// design: terminal`, holds `TerminalColors.standard`: the terminal's `ground` and `ink`. `Palette.all` lists the four palettes in menu order and `Palette.standard` is Rosewater, as tokens.json says.

### iPhone: ios/Irie.swiftpm/Sources/Look.swift

The iPhone app draws in Evergreen until a core sends the Mac app's palette. It holds each set as a dictionary by theme-file name (`"needs_tint": 0xFBE5DC`) between `// design: evergreen` and `// design: end`, and the terminal colours between `// design: terminal termBack=ground termInk=ink`. A second, shorter `evergreen` block gives the seven colours the agents' colours are worked out against (`sidebar`, `ground`, `card`, `panel`, `selection`, `muted`, `needs`).

### The agents' numbers

The agents' colours are worked out in the app from tokens.json's `agents`. The numbers sit as Swift constants, and `internal/designcheck` finds each one by pattern:

| Swift (Theme.swift) | tokens.json path |
|---|---|
| `agentSlots = 6` | `agents.slots` |
| `let gap = 30.0, step` | `agents.needs_gap` |
| `markChroma`, `markLightness = (light:, dark:)` | `agents.mark.chroma`, `agents.mark.lightness` |
| `inkLightness`, `washChroma`, `washLift` | `agents.ink.lightness`, `agents.wash.chroma`, `agents.wash.lift` |
| `high ? 4.5 : 3.2`, `high ? 7 : 4.6` | `agents.mark.contrast`, `agents.ink.contrast` |
| `mix(s.card, full, t)), muted) >= 4.6` | `agents.wash.contrast` |

Look.swift has the same constants, with `static let slots = 6` and its own patterns for the contrast figures. A literal such as `30.0` keeps its point when synced.

## Roles, Make a Theme and theme files

mac/Sources/Irie/Looks.swift holds `ColorRole.all`: the 21 roles in the theme files' order, each with its `key` (`needs_tint`), a key path into `Swatches`, a `title` for Make a Theme (`Card that needs you`) and a `kind`. The kinds are `neutral`, `accent`, `state` and `sun`, which match `roots.Group`; the kinds' titles are "Grounds and text", "Accent", "States" and "Mark".

`Swatches.hexes` writes a set as `#RRGGBB` by key, and `Swatches.with(_:)` reads a theme file's roles over a set, reporting `X is no colour of a theme` and `X = V is not #RRGGBB`. A person's own theme is `themes/<id>.json` in Irie's config folder; a role it leaves out keeps the colour of the palette it started from. Make a Theme (ThemeMaker.swift) draws the window.

Role keys, titles and kinds are kept here by hand. Irie's `TestRolesMatchTheMacApp` (internal/designcheck) reads `ColorRole.all` and fails when a key, a title or a kind is not the library's from `roots.Roles()`. It does not check that `Swatches` has a field for each role. If a role is added in Roots, add it here by hand ([contributing.md](contributing.md)).

## SwiftUI parts

[components.md](../components.md) §2 and §3 name each part's Swift type. The main ones, in mac/Sources/Irie/ and ios/Irie.swiftpm/Sources/:

| Part | Mac | iPhone |
|---|---|---|
| State mark | `StateMark` (States.swift) | `StateMark` (Look.swift) |
| Sun mark | `SunMark` (Theme.swift) | `SunMark` (Look.swift) |
| Count | `Counts` (SidebarView.swift) | `NeedsPill` (Home.swift) |
| Session card | `SessionCard` (SidebarView.swift) | `SessionCard` in `CardRow` (Home.swift) |
| Quick actions | `QuickActions` (SidebarView.swift), `PauseButton`, `PromoteButton`, `ParkButton` (Controls.swift) | none |
| Group head | `GroupHeader` (SidebarView.swift) | heads in Home.swift |
| Task bar | `TaskBar` (TaskBar.swift) | `TaskBar` (TaskBar.swift) |
| Soft and start buttons | `SoftButton` (ConversationView.swift) | `PrimaryButton`, `SecondaryButton`, `SendButton` (Parts.swift) |
| Banner | `ReadyBanner` (Onboarding.swift) | `Banner` (Parts.swift) |
| Notice | `NoticeToast` (Notices.swift) | `ToastStack` (Parts.swift) |
| Details panel | `DetailsPanel` (SidebarView.swift) | `DetailsPanel` (SessionScreen.swift) |
| Tab rows | `IndexCard`, `CardTab` (Rolodex.swift), `SpaceBar` (Spaces.swift), `RightSidebar` (RightSidebar.swift) | none |

Where the two apps draw a part differently, components.md marks it **To converge**. The iPhone app does not yet honour Increase Contrast or Reduce Motion ([README](../README.md) §5).

## Planned: Swift written from the library

Irie's design-library proposal plans, as step 3 of its steps after the first, that the Mac and iPhone apps' palettes, terminal colours and agents' numbers are written whole into one Swift file each by `tools/design.sh`, rather than kept by hand between markers. It is not built. Until it is, the markers and `designcheck` keep the copies honest, and a colour change in Roots reaches the apps through `tools/design.sh sync`.
