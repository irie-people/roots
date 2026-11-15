# Contributing

This page expands [AGENTS.md](../AGENTS.md), which applies to every agent and every contributor, and gives a recipe for each kind of change. Agree anything sizeable in an issue first.

## The rules

### Public sources only

- Build from public documentation, standards, and public open-source work whose licence allows it.
- Record every source that shapes a rule, a value or the code in [PROVENANCE.md](../PROVENANCE.md), in the same change: date, URL, and what it was used for. Third-party code also gets its origin and licence listed there.
- Never open, search or copy from private codebases or private notes. A rule known only from private notes stays out. If a local guard blocks something, stop and say what you needed; do not work around it.

### No personal data

Nothing in the repository may identify a person: no names, usernames, emails, home paths, hostnames or employers, in code, tests, docs or screenshots. Use `alice`, `/home/alice` and `example.com`. A rule from a request says when it was asked for, never who asked. The public authors of sources are credited in PROVENANCE.md.

### Commits

- Commit in UTC and sign off every commit (the DCO): `TZ=UTC git commit -s`.
- Write each subject as a [Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0/): `type(scope): description`. The type is one of `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`.
- Say when an agent wrote part of a change.

### Checks

The library uses the Go standard library only. Before a commit, all three must pass; CI runs them on every push to main and every pull request:

```sh
gofmt -l .        # prints nothing
go vet ./...
go test ./...
```

After changing `palettes/` or `tokens.json`, run `go run ./cmd/webtokens` first.

### Windows stay out of sight

Look at the gallery without opening a window on anyone's screen: render it with a headless browser (`--headless --screenshot`). Never open it in front of what someone is looking at.

## Recipe: change a colour or a number

1. Edit the colour in `palettes/NAME.json`, or the number in `tokens.json`. Colours are `#RRGGBB` in capitals.
2. Run `go run ./cmd/webtokens`, then `go test ./...`. `TestPalettesKeepTheirContrast` holds every text role to 4.5:1 and the accent to 3:1 on all eight surfaces ([principles.md](../principles.md) §1.2); `TestWebTokensAreCurrent` fails if you skipped the command.
3. If the change touches the state colours, look at principles §1.7 and §1.8: a change to one colour checks every pair it meets.
4. Add a PROVENANCE.md row if a source shaped it. Commit and push.
5. Take it into Irie, in its checkout:
  - `go get github.com/irie-people/roots@main`
  - `tools/design.sh sync` writes the values into every copy and the themes page's tables.
  - `tools/check.sh ./internal/designcheck ./internal/ui` runs the checks, including the colour-blind tests.
  - Redraw by hand any picture on `site/themes.html` that `tools/design.sh` reports, until the page is made from the library ([copies.md](copies.md)).
  - Commit the new `go.mod`, `go.sum` and the copies in one commit.

A new palette is a new file in `palettes/`, its ID in tokens.json's `palettes`, and a copy in every client that lists the palettes (Irie's Theme.swift). Roots refuses to load if the file count and the list differ.

## Recipe: add a rule

Follow principles §9:

1. Add it to the part of principles.md it belongs to: one sentence, then **Why**, **From** and **Check**. Keep numbers exact.
2. Give the source a row in PROVENANCE.md.
3. Write the check if you can, here or in the project the rule is about. Write "unchecked" if not.
4. If the rule changes a part, change its entry in components.md.

## Recipe: add a part

Follow components §5:

1. Check that no part does the job: a new look for an old job is a variant, not a new part.
2. Build it from roles and marks, never from colours. If it needs a colour no role gives, that is a new role: add it to `roles` in `roots.go` and to every palette in both modes, run the tests, and update every client's list of roles (in Irie, `ColorRole.all` and `internal/ui`).
3. Give it a section in components.md: what it is for, made of, states, use and don't, and each client's type or class.
4. Add its classes to `web/components.css` (every class starts `ir-`) and draw it in `gallery.html` in all its states.
5. If clients draw it differently, say so under **To converge**.

## Recipe: add a copy

See [copies.md](copies.md). In short: mark the blocks in the file, add a `roots.Copy` to the project's list, run `roots.Check`, and add a test that runs it.

## Review checklist

- `gofmt -l .` prints nothing, and `go vet ./...` and `go test ./...` pass.
- `web/tokens.css` and `web/tokens.js` are current.
- Every new source has a PROVENANCE.md row.
- No names, usernames, emails, home paths or hostnames.
- Commits are signed off, in UTC, with Conventional Commit subjects; an agent's part is stated.
- A new rule has Why, From and Check (or "unchecked").
- A new part is in components.md and the gallery, and uses roles only.
- A change to a value says what Irie's `tools/design.sh sync` will change.
