# Irie's components

The parts every Irie client draws, each named once: what it is for, what it is made of, its states, when to use it and when not, and where each client draws it. A part here is the same part in every client, drawn the way each platform draws: SwiftUI on the Mac and iPhone, HTML and CSS in the phone's web app and on Windows, Bubble Tea in the terminal sidebar.

The rules each part follows are in [principles.md](principles.md); the colours it uses are the roles of [README](README.md) §2.

**Where to find each client's part,** in [Irie](https://github.com/irie-people/irie)'s repository. Mac: mac/Sources/Irie/. iPhone: ios/Irie.swiftpm/Sources/. Phone web: internal/remote/web/ (app.css, js/ui.js). Windows: internal/winapp/web/ (app.css, app.js). Terminal: internal/ui/.

## 1. How the parts are built

Irie's parts stack in four layers, each built only from the one under it:

1. **Roles.** The 21 colours of a palette, by what they are for: surfaces to draw on, text, the accent, the states, the icon's sun and sea (README §2). A part never names a colour, only a role, so every palette and both modes, and a person's own theme, come free.
2. **Marks.** The smallest things that mean something: a state's glyph, the sun, a count, a dot. Each carries its meaning in shape as well as colour (principles §1.7).
3. **Parts.** A session card, a group head, the task bar, a drawer: marks and text on a surface, with states of their own.
4. **Patterns.** How parts work together: loud and quiet, the working wash, the button as the mark, the disclosure ladder (§4).

One more rule sits under all four: **the core owns every rule, the clients only draw.** What state a session is in, which colour it takes, where it sits in the order, what a card's line says: the core decides, and sends it in each row. A client that works one of these out for itself will disagree with the others sooner or later.

**Densities.** The pointer clients (Mac, Windows) are dense: text around 13 points, cards of 14 (Mac) or 8 (Windows) corner radius, buttons of 16 to 30 points. The touch clients (iPhone, phone web) give each target at least 44 points, names at 16.5 points and cards a radius of 16. A part keeps its anatomy and states in both; only its sizes change.

## 2. Marks

### 2.1 State mark

The glyph for a session's state, before its name (Irie's DESIGN §5.3): `●` needs you, `!` stalled, `✕` failed or exited with an error, `◐` working, `◉` waiting on its helpers, `‖` paused or pausing, `◆` frozen, `✓` done, `○` idle, `◌` starting, `·` exited. Each state has its glyph and its words, so colour never carries the meaning alone.

- **Made of:** the glyph in its state's role (`needs`, `done`, `paused`, `frozen`, `muted`), or in the agent's colour while it works; `needs` is a dot with a 22% halo.
- **Use it** wherever a session is named: cards, menus, notices, the conversation's header.
- **Don't** animate it: marks hold still in a list (principles §5.1).
- **Mac:** `StateMark` (States.swift), 14 pt, 13 looks. **iPhone:** `StateMark` (Look.swift), 16 pt, 12 looks: it has no waiting yet. **Phone web, Windows:** `.mark` with `.r-needs`, `.r-done`… (look.js has the glyphs). **Terminal:** the glyph in its colour (card.go).
- **To converge:** a working mark is ink on the Mac and the agent's colour on the iPhone; one of them should change.

### 2.2 Sun mark

Irie's icon, a sun over two waves, in the `sun` and `waves` roles. It stands for Irie in title bars, the welcome and empty pages.

- **Mac:** `SunMark` (Theme.swift), which Make a Theme also draws in a draft's colours. **iPhone:** `SunMark` (Look.swift). **Phone web, Windows:** `.sun` SVG.

### 2.3 Count

A number with a state's mark, counting the sessions in that state: the title bar's counts, a group's "· 2 need you". The count of what needs you is red; a count of nothing is faint, never red (principles §1.5).

- **Made of:** a state mark at 10 to 14 points and tabular figures, on a capsule that fills under the pointer.
- **Use it** to say how many, and to jump to the first of them on a click.
- **Mac:** `Counts` (SidebarView.swift). **iPhone:** `NeedsPill` (Home.swift), needs only, which filters the list. **Phone web:** `.needs-pill`. **Windows:** `.count`. **Terminal:** the header line (view.go).

### 2.4 Status dot

A dot that says whether a machine is reachable: `done` connected (with a soft halo), `paused` connecting, `muted` offline, `needs` not paired.

- **iPhone:** `StatusDot` (Parts.swift). **Phone web:** `.dot`. Its "connecting" breathes: the one mark allowed to move, because it is live.

## 3. Parts

### 3.1 Session card

One session in the list: its state mark, its name and age on the first line, what it is doing in plain words on the second, and up to three of its helpers under that. The card is Irie's main part; most of the principles are about it.

- **Made of:** the `card` surface (flat at rest, principles §2.5), a state mark, the name in `ink` (or `muted` when the session is quiet), the line in `muted`, a fixed-width trailing slot for the age, and the quick actions (§3.2) in that slot under the pointer.
- **States:**
  - *quiet*, nothing running: a muted name.
  - *working*: washed in the agent's colour (§4.2), its line in `ink`.
  - *needs you* (loud): the `needs_tint` fill, a `needs` edge at half strength, the line in `needs` and bold.
  - *selected*: the `selection` fill, an edge in the agent's colour or the accent, bolder while the list has the keyboard, one gentle shadow, and a line to its terminal.
  - *hover*: the `hover` fill and a faint accent edge.
  - *loose*, work that has not landed: the line in `accent_ink`.
  - *parked*: in the Parked group, otherwise the same.
  - *unreachable* (phone): drawn at 62%.
- **Use it** for a session, anywhere sessions are listed. For a session outside Irie, the one-line row (`OutsideRow`).
- **Don't** move it when its state changes (principles §2.1), or add a badge for a state its button shows (§4.3).
- **Mac:** `SessionCard` (SidebarView.swift): two lines, radius 14, padding 8 by 12, a 66 pt trailing slot; click shows, ⌘ and ⇧ pick, double-click renames, drag reorders, right-click opens its menu; ↑↓, z, Z from the keys. **iPhone:** `SessionCard` in `CardRow` (Home.swift): radius 16, swipe for Close and Park, long press for its menu. **Phone web:** `.card` (and `.row` for the swipe). **Windows:** `.card`, radius 8. **Terminal:** card.go, one line where the pane is short.

### 3.2 Quick actions

The buttons that show on the card under the pointer: pause, boost, park, and `⋯` for the rest. A button whose state is on stays on the card, filled, in its own place, as that state's one mark (§4.3).

- **Made of:** icon buttons 16 to 18 points wide in the card's trailing slot, where the age was.
- **Use it** for the few actions people take most, one click away. Everything else goes in the `⋯` menu and the card's own menu, in the Session menu's words and order (principles §4.4).
- **Mac:** `QuickActions` (SidebarView.swift) with `PauseButton`, `PromoteButton`, `ParkButton` (Controls.swift). Only the card under the pointer builds them. **iPhone, phone web:** none; swipe actions and the long press take their place. **Windows:** in the card's menu.
- **Close mark:** under the pointer a card's state mark gives way to a `×`, as a browser's tab does (`CloseMark`, Tabs.swift).

### 3.3 Helper rows

A card's helpers and teammates, under it: `↳ Explore  Reading DESIGN.md`, three at most, then "and N more". A `×` under the pointer hides one from the card; it keeps working.

- **Mac:** `HelperRows` with `SubRow` (SidebarView.swift, Tabs.swift). **iPhone, phone web:** in the card (`.c-helper`). **Terminal:** in card.go.

### 3.4 Group head

The header of a group of cards: a chevron, the group's name, its count, and "· N need you" in red when it is folded and something in it needs you. Pinned heads the list; Parked follows the groups, folded.

- **Made of:** a chevron that turns as it folds, the name in `muted` at a small, heavy size, the count in tabular figures.
- **Use it** to fold, and to drop cards into a group. Double-click renames, right-click opens the group's menu.
- **Mac:** `GroupHeader` (SidebarView.swift), with the group's rail. **iPhone:** the Parked, Closed and machine heads (Home.swift). **Phone web:** `.divider`, `.msec-head`. **Windows:** `.group-head`. **Terminal:** list.go.

### 3.5 Task bar

Where a person types a task, at the foot of the sidebar: three lines tall, always there, with the start button at its end. Nothing appears around it or is drawn over it.

- **Made of:** the `field` surface, an accent focus ring, the start button (§3.6), and the repo it will start in.
- **Use it** for new work; ↑ and ↓ bring back earlier tasks.
- **Mac:** `TaskBar` (TaskBar.swift), radius 16. **iPhone:** `TaskBar` (TaskBar.swift), radius 23, with chips while it has the keyboard. **Terminal:** command.go, the command bar.

### 3.6 Buttons

- **Start button:** a round button in `accent_fill` with an `on_accent` arrow: starts a task, sends a message. **Mac:** in `TaskBar` and the composer. **iPhone:** `SendButton` (Parts.swift), 34 pt with its touch area 6 pt wider.
- **Soft button:** a capsule on `card` with an `edge` hairline; prominent, it takes `accent_fill`. **Mac:** `SoftButton` (ConversationView.swift). **iPhone:** `PrimaryButton`, `SecondaryButton` (Parts.swift), with a danger look in `needs`.
- **Icon button:** an icon in `muted` that turns `ink` on a `hover` fill. **Phone web, Windows:** `.icon-btn`.
- **Use** one prominent button per view at most; every other button is soft or an icon. A button says a verb and its object (principles §6.2).

### 3.7 Segmented control

A short row of choices, one picked: the List and Repo views. The picked one is a `card` thumb with a shadow on a `panel` track. When it swaps a panel, it is a tab row in this dress and takes the tab list's roles and keys (§3.13, principles §4.6); when it only sets a mode, its buttons are pressed or not (`aria-pressed`), as the gallery's pickers are.

- **Mac:** `ModeSwitch` (SidebarView.swift), at most 190 pt wide. **Windows:** `.seg`.

### 3.8 Field

A one-line text field, such as Find: the `field` surface, an `edge` hairline that turns accent with a focus ring while it has the keyboard.

- **Mac:** `FindField` (SidebarView.swift). **iPhone:** `fieldStyle` (Pairing.swift). **Windows:** `.find`.

### 3.9 Banner and notice

- **Banner:** a note across a part that stays while what it says is true: Irie is not ready yet, a machine is unreachable. Loud ones take `needs_tint`. **Mac:** `ReadyBanner` (Onboarding.swift), the composer's banner. **iPhone:** `Banner` (Parts.swift). **Phone web:** `.banner`. **Windows:** `.ready-note`.
- **Notice:** a message that comes and goes at the top right of the work, with Undo where it can: errors stay at least 8 seconds and while the pointer is on them, and every notice is kept in the Notices tab (principles §7.3). **Mac:** `NoticeToast` (Notices.swift), on `field`, edged in `needs` for an error. **iPhone:** `ToastStack` (Parts.swift), on the ink, inverted.
- **To converge:** the Mac's notice is light on light and the iPhone's inverted; pick one.

### 3.10 Details panel

What a card holds that its two lines cannot: the exact command, the repo and branch, the costs. It opens at one line and folds back to it, never showing half a line.

- **Mac:** `DetailsPanel` (SidebarView.swift): `panel`, radius 14, edged in the state's colour, a drag grip, `DetailsFold`. **iPhone:** `DetailsPanel` (SessionScreen.swift), a grid of names and values.

### 3.11 Drawer

A panel that slides in beside the work, never over it: the right sidebar's tabs (helpers, review, notices, mail, plugins). The terminal resizes once as it opens, closes or its edge is let go.

- **Mac:** `RightSidebar` (RightSidebar.swift): 340 pt wide by default, 260 to 520, its tab names turning to icons when narrow; `.smooth(0.28)`, one step under Reduce Motion.

### 3.12 Conversation

A session's turns as a person reads them:

- **Your message:** on `card` (Mac `YouTurn`), or a bubble in `accent_fill` on the right (iPhone `UserBubble`); dashed or faint until the agent has it.
- **The agent's text:** Markdown (`MarkdownText`, `MarkdownView`).
- **A tool call:** one line in plain words that opens to its exact command (`CallRow`, `ToolCall`), and a group of calls folds into one (`GroupRow`, `GroupCall`).
- **Needs card:** what the agent asks, with its options (`NeedsCard` in both apps, which differ a lot: **to converge**).
- **Note line:** a hairline, a few words, a hairline, between turns (`NoteLine`).

### 3.13 Tab row

A row of tabs, one picked, that swaps what shows under it: the panes' index tabs, the right sidebar's tabs, a code card's open files, the spaces, Settings' pages. The picked tab either joins what it shows, as a folder's tab does, or fills on the row, as a segmented control's thumb does (§3.7).

- **Made of:** labels in `muted` on the surface around the row; the picked tab and its panel in `card`, its label in `ink` and heavier; an inset `hover` fill under the pointer, its corners the tab's less the inset (principles §2.6).
- **Folder tabs:** the picked tab and its panel are one outline. Its top corners are round, and where each side meets the panel a concave curve of the same radius turns into the panel's top edge. Beside the panel's edge the curve and the panel's corner share the room (principles §2.7): a first tab flush with the edge has no curve there, and that corner of the panel is square until the row scrolls it away. The row leaves room for both after its last tab.
- **States:**
  - *picked*: joined to its panel (folder) or filled (segmented); its label heavier, at the width it always keeps, so nothing after it moves (principles §2.3).
  - *hover*: the inset fill. *Focus*: a ring drawn inside the tab, so the row's edge never clips it (principles §4.5).
  - *overflowing*: the row scrolls with no scrollbar, and each side where tabs hide gets a fade and a chevron button, named "Previous tabs" and "Next tabs" and placed outside the tab list; a chevron goes when its side has nothing left. The picked tab scrolls into view, clear of the fade (principles §2.8). Chevrons show only with a fine pointer; a finger swipes.
  - *moving*: when the pick changes, the picked box may slide to its new tab in about 0.2 s; under Reduce Motion it jumps (principles §5.3).
- **Keys:** on the web, the tab list's roles and keys (principles §4.6). On the Mac, each row's tabs have keys of their own: ⌥⌘1–9 for the right sidebar, ⌃1–9 for spaces, ⌥Tab through the panes.
- **Use it** when each tab has its own content in one place. For a few views of one thing, a segmented control (§3.7); for pages, links in a navigation row.
- **Don't** let a picked label widen the row, hide a scrollbar without the fade and chevron, or give a row the tabs' roles without their keys.
- **Mac:** the panes' index tabs (`IndexCard`, `CardTab`, `RolodexLayout`, Rolodex.swift): folder tabs drawn as one outline with tangent arcs, their widths set by the row, so nothing moves. The right sidebar's tabs (RightSidebar.swift), their names turning to icons when narrow (`ViewThatFits`). Space chips (`SpaceBar`, `SpaceChip`, Spaces.swift). A code card's files (CodePane.swift). **Phone web:** Browser and Term (`.seg`, js/conv.js). **Windows:** List and Repo (`.seg`), the spaces (`.spaces`) and Settings' pages (`.set-tabs`, js/settings.js). **Web:** `.ir-tabs` (web/components.css), with the keys and chevrons in the gallery's script.
- **To converge:**
  - The Mac's right sidebar tabs, space chips and code card files turn the picked label bold while each sizes to its text, so the labels after it move at every pick, by 1 to 5 points (measured in the system fonts at their sizes).
  - A code card's files and the space chips scroll with no scrollbar, fade or chevron, and a newly opened file's chip is added at the end, where it can land out of sight.
  - The web apps' four tab rows have the tabs' roles but not their keys; Settings' row has no name and, as written, loses the keyboard when it draws itself again.
  - A segmented control is pressed buttons in the gallery and tabs in the web apps; each should follow what it does (§3.7).

## 4. Patterns

### 4.1 One loud thing

Of everything on screen, only what needs the person is loud: red text, a red-tinted card, a red count. Everything else is quiet, so the loud thing is found at a glance (principles §1.5). A new part that wants attention takes a state's mark, never red of its own.

### 4.2 The working wash

While a session works, its card is washed in its agent's colour: the card's lightness, a little lifted, at a low chroma, fading toward the card. Its mark and helpers take the same hue. The selected working card turns its wash around and deepens toward its terminal. The steps, and their numbers, are the same in every client (tokens.json, `agents`); under Increase Contrast the wash goes.

### 4.3 The button is the mark

A toggle's own button, kept on the card and filled, shows its state: one mark, in one place, which is also how to undo it (principles §4.1).

### 4.4 The disclosure ladder

Each step shows more, and nothing is two steps away from where the person is: the card's two lines, then the details panel at one line, then the panel opened, then the conversation or the terminal. Hover shows a card's buttons; a menu holds the rest.

### 4.5 Places that hold still

A part keeps its place and size while it changes: a state change changes colours, not places; hover swaps a slot of fixed width; a rename edits in place (principles §2.3).

### 4.6 Undo, not "Are you sure?"

Act at once, show a notice with Undo, and keep it in the history (principles §4.2).

## 5. Adding a part

1. Check that no part here does the job: a new look for an old job is a variant of that part, not a new part.
2. Build it from roles and marks, never from colours (§1). If it needs a colour no role gives, that is a new role, which every palette must then give in both modes and the contrast tests check.
3. Give it a section here: what it is for, made of, states, use and don't, and each client's type or class.
4. Draw it in the gallery ([gallery.html](gallery.html)) in its states, so it can be seen in every palette and mode.
5. Parts drawn in more than one client should match: when they don't, say so under **To converge**.
