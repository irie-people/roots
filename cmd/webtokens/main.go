// Command webtokens writes the library's files for web pages,
// web/tokens.css and web/tokens.js, from its palettes and tokens.json. Run
// it from the library's root after changing either.
package main

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/irie-people/roots"
)

func main() {
	files := roots.WebTokens()
	for _, path := range slices.Sorted(maps.Keys(files)) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := os.WriteFile(path, files[path], 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("wrote", path)
	}
}
