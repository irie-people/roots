# Provenance

The design library is written from public documentation, standards and public open-source work, and from requests recorded in Irie. Every source that shaped a rule, a value or the code is listed here: what it is, when it was read, and what it was used for.

## Sources

| Date | Source | Used for |
|---|---|---|
| 2026-11-15 | Irie, https://github.com/irie-people/irie (Apache-2.0; this library began as its design/ folder, commit 9245b1f): its palettes (internal/ui, mac/Sources/Irie/Theme.swift, Looks.swift's role names), docs/DESIGN.md, docs/MAC.md, the requests its PROVENANCE.md records, and its tests | The palettes, the agents' numbers, the roles and their names, the rules and their checks (principles.md), the parts and where each client draws them (components.md) |
| 2026-11-15 | https://www.w3.org/TR/WCAG22/ (relative luminance and contrast ratio; 1.4.1, 1.4.3, 1.4.6, 1.4.11, 2.2.1, 2.3.3, 2.4.7, 2.5.7) | `Luminance` and `Contrast`, the contrast floors in tokens.json, `TestPalettesKeepTheirContrast`, the gallery's figures; principles §1.2, §1.7, §4.3, §4.5, §5.3, §7.1, §7.3 |
| 2026-11-15 | https://bottosson.github.io/posts/oklab/ (Björn Ottosson, "A perceptual color space for image processing": the Oklab matrices and OKLCH) | The agents' colours in OKLCH (tokens.json, `agents`); the gallery's hues; principles §1.3, §1.6 |
| 2026-11-15 | https://drafts.csswg.org/css-color-4/ (`oklch()`, `color-mix()`, deltaE-OK) | web/components.css and the gallery; principles §1.7, §1.8 |
| 2026-11-15 | https://drafts.csswg.org/mediaqueries-5/ (`prefers-color-scheme`, `prefers-reduced-motion`) | web/tokens.css's modes; web/components.css under Reduce Motion |
| 2026-11-15 | https://pkg.go.dev/embed ; https://pkg.go.dev/encoding/json (`Decoder.DisallowUnknownFields`, `Decoder.UseNumber`, `Number`) | Embedding and reading palettes/ and tokens.json in design.go |
| 2026-11-15 | https://github.com/DaltonLens/libDaltonLens (Brettel, Viénot and Mollon 1997) ; https://medlineplus.gov/genetics/condition/color-vision-deficiency/ | principles §1.7 |
| 2026-11-15 | https://www.science.org/doi/10.1126/science.3283936 (Livingstone and Hubel 1988) | principles §1.1 |
| 2026-11-15 | https://m2.material.io/design/color/dark-theme.html | principles §1.9 |
| 2026-11-15 | https://developer.apple.com/design/human-interface-guidelines/menus | principles §4.4 |
| 2026-11-15 | https://github.com/heydamir/unslop (Damir Sokolovsky, CC BY 4.0) | principles §1.10, §2.6, §5.1, §6.2 and §8: what not to draw, spacing rhythm and calm words, in this library's own words; nothing copied |
| 2026-11-15 | https://www.designtokens.org/tr/2025.10/ | Weighed for tokens.json's form; kept for an export (README §5) |
