# Copies and markers

A Go program reads Roots directly. A file that cannot, such as a Swift app's palette or a web page's CSS, keeps a copy of the values. Each copied block sits between two marker comments, and the library finds the blocks, checks them with `roots.Check` and writes them with `roots.Sync`. A project lists its copies as `[]roots.Copy` ([go.md](go.md)).

## The marker

A marker is a line with `design:` after a comment opener, in the file's own syntax:

| Syntax | Marker |
|---|---|
| `//` (Swift, Go, JavaScript) | `// design: evergreen` |
| `/* */` (CSS) | `/* design: evergreen light */` |
| `<!-- -->` (HTML) | `<!-- design: terminal termInk=ink -->` |

A block opens with `design: NAME [options]` and closes with `design: end`. A line holding a marker is read as a marker and nothing else on it counts. Blocks do not nest, and a file may hold the same block name more than once.

```css
:root {
  /* design: evergreen light page=sidebar */
  --page: #F4F4EC; --card: #FBFBF6; --needs-tint: #FBE5DC;
  /* design: end */
  --bar: #183F30;
}
```

## Block names

| First word | Meaning | Which set the colours are in |
|---|---|---|
| A palette ID and `light` or `dark` (`evergreen dark`) | One mode of a palette | The one named |
| A palette ID alone (`evergreen`) | Both modes | A `light` or `dark` key before the colours, as `light: Swatches(` or `"dark": {` |
| `terminal` | The terminal's colours | The terminal's `ground` and `ink` |

In a palette-alone block, a key applies to the colours after it on its line and the lines below, until the next key. A colour that comes before any key is a problem.

## What a block may hold

Every name that has a colour in a block must be a role of the palette, or an alias. A colour is `#RRGGBB` or `0xRRGGBB`, optionally inside one call, as `fixed(0xEEE5D5)`. The name is followed by `:` or `=`.

- **Spellings.** Kebab, camel and snake are all the same role: `--needs-tint`, `needsTint` and `"needs_tint"` are `needs_tint`. A leading `--` and quotes are ignored.
- **Aliases.** Options after the name read `name=role` and map a name of the file's own to a role: `page=sidebar` reads `--page` as `sidebar`; `termInk=ink` reads `termInk` as `ink`. Spellings are normalised on both sides.
- **Outside a block.** A colour of the file's own, such as `--bar: #183F30;`, goes outside the blocks. Blocks hold only the library's colours.
- **Case.** Hex digits are compared ignoring case. `Sync` writes the library's capitals.

## Numbers

The agents' numbers are not in blocks, since each file names them its own way. A `roots.Number` has a regular expression and the dotted paths of tokens.json for its groups:

```go
{Pattern: `markLightness = \(light: ([\d.]+), dark: ([\d.]+)\)`,
 Paths: []string{"agents.mark.lightness.light", "agents.mark.lightness.dark"}}
```

- The pattern must match **once** in the file, with one group for each path.
- Numbers are compared by worth, so `0.50` is `0.5`.
- A Swift `Double` keeps its point: if the copy writes `30.0` and the library says `30`, `Sync` writes `30.0`.

## Messages

`Problem.String` prints `path:line: what`, or `path: what` when there is no line. `Check` and `Sync` use these messages, with `%s` filled in.

| Message | What to do |
|---|---|
| `X is V here and W in the library` (`X` is `evergreen light needs_tint` or `agents.mark.chroma`) | The copy drifted. Run the sync, or change the library first if the copy is the right one |
| `X is no colour of the B block: name it in its design: line (X=ROLE), or move it out of the block` | A name that is not a role sits in a block. Add the alias, or move the line outside |
| `X comes before a light or dark key in the B block` | In a palette-alone block, put the `light`/`dark` key before it |
| `the B block has no design: end` | Close the block |
| `design: end with no block open` | Remove the extra end, or add the opener |
| `the B block has no colours` | Nothing in it matched; check the spelling and the `#` or `0x` |
| `design: N inside the B block from line L` | A block opened inside another; close the first |
| `design: N names no palette, nor terminal` | Fix the name; the IDs are in tokens.json |
| `design: "O" is not name=role` | An option is not an alias |
| `a design: line with no name` | Add the block name |
| `has the blocks [...], Copies says [...]` | The file and its `Copy.Blocks` disagree; edit either |
| `N matches of P, want 1` | A number's pattern is wrong or ambiguous; make it match exactly once |
| `tokens.json has no PATH` | The path in `Paths` is mistyped |

An unreadable file's error text is reported as is. `Sync` changes nothing in a copy that has any problem and reports them. It lists in `changed` each file it wrote.

## How Irie lists its copies

Irie's list is `Copies` in `internal/designcheck/designcheck.go`; each entry has a `Why`. The files are the Mac app's Theme.swift, the iPhone app's Look.swift, the phone web app's look.js, app.css and setup.html, and the Windows app's app.css. The Swift and JavaScript ones also give `Number` patterns for the agents' numbers.

`designcheck.Check` also reads `site/themes.html`. Its pictures may hold only a palette's colours or the mixes the pictures use, and each row of its tables must give the role's colour and its lowest contrast on five grounds. A stale picture reports, for example, `the overview has #XXXXXX, which is no palette's colour: draw it again from the library`. `Sync` writes the tables; a picture with an old colour has to be drawn again by hand, until the page is made from the library.

```sh
tools/design.sh         # says which copies are not the library's
tools/design.sh sync    # writes the library's values into them and the themes tables
```

With no problem it prints `every copy is the library's`. Otherwise it prints each problem and `tools/design.sh sync writes the library's values into the copies`, and exits 1. `TestCopiesMatchTheLibrary` in Irie runs the same check with the other tests. Beside it, `TestRolesMatchTheMacApp` checks Looks.swift's role keys, titles and groups against `roots.Roles()`.

## Adding a copy

1. Mark each block of the library's colours in the file with `design: NAME` and `design: end`, in its comment syntax.
2. Use role names in any spelling, or add aliases in the opening marker. Move the file's own colours outside.
3. Add a `roots.Copy` to your list with `Path`, `Why`, `Blocks` (every block name, repeats included) and, for the agents' numbers, `Numbers`.
4. Run the check. Fix what it says, or run `Sync` to write the library's values in.
5. In Irie, add the file to `Copies` in `internal/designcheck/designcheck.go` and run `tools/design.sh`.
