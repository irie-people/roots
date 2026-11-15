# The web parts

Roots gives a web page two stylesheets and one script: `web/tokens.css` holds the palettes as CSS variables, `web/components.css` holds the parts as `ir-` classes built only from those variables, and `web/tokens.js` gives a page's scripts the same values. Link `tokens.css` first.

## tokens.css

`go run ./cmd/webtokens` writes it from `palettes/` and `tokens.json`; change those, not the CSS ([go.md](go.md)). It sets one variable per role, named by turning `_` into `-`, so `needs_tint` is `--needs-tint`:

| Group | Variables |
|---|---|
| Surfaces | `--sidebar`, `--ground`, `--card`, `--panel`, `--field`, `--hover`, `--selection`, `--needs-tint` |
| Hairline | `--edge` |
| Text | `--ink`, `--muted`, `--accent-ink`, `--needs`, `--done`, `--paused`, `--frozen` |
| Accent | `--accent`, `--accent-fill`, `--on-accent` |
| Icon | `--sun`, `--waves` |

It also sets `--term-ground` and `--term-ink` on `:root`, the terminal's colours, the same in every palette. Each block sets `color-scheme` to `light` or `dark`. What each role is for is in [README](../README.md) §2.

### Choosing a palette and a mode

| You write | You get |
|---|---|
| Nothing | The default palette (Rosewater), in the system's mode |
| `data-palette="lilac"` on an element | That palette for the element and what is inside it |
| `data-mode="dark"` or `"light"` on the same element | That mode, whatever the system says |

The palette IDs are `rosewater`, `terracotta`, `lilac` and `evergreen`. For each palette the file has three kinds of rule:

1. The light set, on `[data-palette="evergreen"]`. For the default palette the selector also lists `:root`.
2. The dark set, on `[data-palette="evergreen"][data-mode="dark"]`.
3. The same dark set inside `@media (prefers-color-scheme: dark)`, on `[data-palette="evergreen"]:not([data-mode="light"])`. This is the system mode: with no `data-mode`, a dark system gets the dark set; `data-mode="light"` opts out.

```html
<section data-palette="terracotta" data-mode="dark">
  <!-- everything in here is Terracotta, dark -->
</section>
```

### Theming hooks only

`data-palette` and `data-mode` exist so the tokens can be picked. Do not use them as application state: do not select on them in your own CSS to change layout, and do not read them in a script to decide behaviour. A part must look right in every palette and mode because it uses roles, never colours ([components.md](../components.md) §1). If a part needs a different colour in one palette, it needs a new role, not a selector.

### tokens.js

It sets one global, `window.IRIE_DESIGN`, with `roles` (key, title, group, use), `palettes` (id, name, light, dark), `default`, `contrast`, `terminal` and `agents`. It is a plain script, not a module, so a page opened from a file can read it.

## components.css

Every class starts with `ir-`, so a page's own classes never meet them. Put class `ir` on an ancestor: it sets the font, text colour and the `--ir-*` variables (`--ir-r`, `--ir-pad`, `--ir-target` and others) that the parts read.

### Density

`data-density="touch"` on any element gives what is inside it the phone's sizes: 15 px text, 16 px card radius, 44 px targets. Without it, parts take the pointer clients' sizes: 13 px text, 14 px radius, 28 px targets. Put it on or inside the `ir` element.

### Parts

| Class | Part | Modifiers and states |
|---|---|---|
| `ir-mark` | State glyph | `data-state`: `needs`, `stalled`, `failed`, `error`, `working`, `waiting`, `paused`, `pausing`, `frozen`, `done`, `idle`, `starting`, `exited` |
| `ir-count` | Count of sessions in a state | `is-needs` (red), `is-zero` (faint) |
| `ir-card` | Session card | `is-quiet`, `is-working`, `is-loud` (needs you), `is-loose`, `is-selected`, `is-keyboard`, `is-hover` |
| `ir-act`, `ir-actions` | Quick actions on a card | `is-on` (stays filled) |
| `ir-group` | Group head | `aria-expanded` turns the chevron |
| `ir-taskbar`, `ir-start` | Task bar and its start button | Focus ring while it has the keyboard |
| `ir-button` | Soft button | `is-prominent`, `is-danger` |
| `ir-icon-button` | Icon button | None |
| `ir-seg` | Segmented control | `aria-pressed="true"` on the picked button |
| `ir-field` | One-line text field | Focus ring |
| `ir-banner` | Banner | `is-loud` |
| `ir-notice` | Notice | `is-error` |
| `ir-panel` | Details panel | Holds a `dl` |
| `ir-tabs` | Tab row and its panel | `aria-selected`, `is-more-start`, `is-more-end` |

What each part is for, and when not to use it, is in [components.md](../components.md) §3. This page covers how to write the HTML.

A card, from the gallery:

```html
<div class="ir-card is-loud">
  <span class="ir-mark" data-state="needs"></span>
  <div class="ir-card-main">
    <div class="ir-card-top"><span class="ir-card-name">api-fix</span><span class="ir-card-slot"><span class="ir-age">12m</span></span></div>
    <div class="ir-card-line">Asks: may it push to main?</div>
  </div>
</div>
```

The card's inner classes are `ir-card-main`, `ir-card-top`, `ir-card-name`, `ir-card-slot`, `ir-card-line`, `ir-card-helpers` and `ir-card-helper`. The slot keeps one width (66 px). Under the pointer, or with `is-hover`, an `ir-actions` inside the slot replaces the `ir-age`:

```html
<span class="ir-card-slot"><span class="ir-age">26m</span><span class="ir-actions">
  <button class="ir-act" type="button" aria-label="Pause"><svg viewBox="0 0 16 16" aria-hidden="true">…</svg></button>
</span></span>
```

A group head, a count and a banner:

```html
<button class="ir-group" type="button" aria-expanded="true"><span class="ir-group-chev">›</span><span class="ir-group-title">irie</span><span class="ir-group-count">5</span></button>
<button class="ir-count is-needs" type="button"><span class="ir-mark" data-state="needs"></span>2</button>
<div class="ir-banner is-loud">Another machine can’t be reached: its sessions wait</div>
```

Tabs, in short (the full row is in `gallery.html`):

```html
<div class="ir-tabs">
  <div class="ir-tabs-row">
    <div class="ir-tablist" role="tablist" aria-label="Panes">
      <button class="ir-tab" type="button" role="tab" aria-selected="true" tabindex="0"><span class="ir-tab-label" data-text="Terminal">Terminal</span></button>
      <button class="ir-tab" type="button" role="tab" aria-selected="false" tabindex="-1"><span class="ir-tab-label" data-text="Browser">Browser</span></button>
    </div>
  </div>
  <div class="ir-tabpanel" role="tabpanel">A shell in the session's folder.</div>
</div>
```

The stylesheet draws tabs but has no script. A page sets `is-more-start` and `is-more-end` on the tab list as it scrolls, shows the `ir-tabs-more` chevron buttons (`is-start`, `is-end`), and handles the keys of the tabs pattern ([principles.md](../principles.md) §4.6). The gallery's script is the example. `data-text` on `ir-tab-label` repeats the label so a picked, heavier label never changes the row's width.

### The agent variables

A working card is washed in its agent's colour. The page sets three variables on each card:

| Variable | Used for |
|---|---|
| `--a-mark` | The working mark and the selected card's edge |
| `--a-wash` | The working card's background, fading toward `--card` |
| `--a-ink` | The text of the card's helpers |

Without them the parts fall back to `--ink`, `--card`, `--accent` and `--muted`. The stylesheet does not compute them: the client works them out from the palette and the numbers in tokens.json (`agents`), and sets them in `oklch()`, as the gallery does with `el.style.setProperty('--a-mark', …)`. The gallery uses each number's starting point only; Irie's clients also move the mark until it keeps its contrast (components.md §4.2, `agents` in tokens.json).

## Accessibility notes

- Focus: `.ir :focus-visible` draws a 2 px outline in `--accent`, offset 2 px. Tabs draw it inside the tab so the row never clips it. The task bar and field show a focus ring in the accent.
- Reduced motion: under `prefers-reduced-motion: reduce`, every transition and animation inside `.ir` is switched off.
- Colour is never the only cue: each state mark has its own glyph.

The rules behind these are in [accessibility.md](accessibility.md) and [sizes-and-motion.md](sizes-and-motion.md); themes in [theming.md](theming.md).
