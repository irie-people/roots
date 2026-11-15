# Working on Roots, Irie People's design library

These rules apply to every agent and every contributor.

## Public sources only

- Build from public documentation, standards, and public open-source work whose license allows it.
- Record every source that shapes a rule, a value or the code in [PROVENANCE.md](PROVENANCE.md) (date, URL, what it was used for) in the same change. Third-party code also gets its origin and license listed there.
- Never open, search or copy from private codebases or private notes. A rule known only from private notes stays out of this repository. If a local guard blocks something, stop and say what you needed; don't work around it.

## No personal data

Nothing in the repository should identify a person: no names, usernames, emails, home paths, hostnames or employers in code, tests, docs or screenshots. Use `alice`, `/home/alice` and `example.com`. A rule from a request says when it was asked for, never who asked. The public authors of sources are credited in PROVENANCE.md.

## Commits

Commit in UTC (`TZ=UTC git commit`) and sign off every commit (`git commit -s`, the DCO). Write each subject as a [Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0/): `type(scope): description`, with `type` one of `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`. Say when an agent wrote part of a change.

## Stack and checks

Go, with the standard library only. Before a commit, `gofmt -l .` prints nothing and `go vet ./...` and `go test ./...` pass; CI runs the same. After changing palettes/ or tokens.json, run `go run ./cmd/webtokens`. The library holds to its own principles: its palettes keep their contrast (`TestPalettesKeepTheirContrast`).

## What goes where

- A colour or a number: palettes/ or tokens.json ([README](README.md) §3). Irie then takes it with `go get github.com/irie-people/roots@main` and writes it into its copies with `tools/design.sh sync`.
- A rule: [principles.md](principles.md), with its Why, From and Check (principles §9).
- A part, or a new state of one: [components.md](components.md) and [gallery.html](gallery.html) (components §5).

## Windows stay out of sight

Look at the gallery without a window on the owner's screen: render it with a headless browser (`--headless --screenshot`), and never open it in front of what they are looking at.
