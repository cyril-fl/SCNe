// Command audit reports strings that bypass i19n and translation keys that do not resolve.
// Run it from the project root: go run .agents/skills/i19n-audit/scripts/audit.go
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// Calls whose string arguments are translated: a literal there must be a key.
var translated = map[string]bool{
	"i19n.Localize": true, "i19n.Unlocalize": true,
	"logs.Error": true, "logs.MildError": true, "logs.SoftError": true, "logs.CeaseError": true, "logs.CatalogError": true,
	"logs.SilentError": true, "logs.SilentWarning": true, "logs.SilentInfo": true, "logs.SilentSuccess": true,
}

// Calls whose Label or ID field is rendered through i19n further down.
var labelled = map[string]bool{
	"NewShortcut": true, "shortcut.NewShortcut": true,
	"navigation.NewSection": true, "navigation.NewFocusableSection": true,
}

// Calls taking text that the player may read but that cannot be a key.
var rawErrors = map[string]bool{"logs.RawError": true, "fmt.Errorf": true, "errors.New": true}

// Calls whose literals are never text for the player.
var technical = map[string]bool{"os.Getenv": true, "exec.Command": true}

var formatVerb = regexp.MustCompile(`%[-+#0 ]*\d*[a-zA-Z]`)

type finding struct{ pos, call, value string }

type report struct{ missing, hardcoded, raw []finding }

func main() {
	keys, locales := loadLocales("locales")
	printLocaleGaps(keys, locales)

	var r report
	for _, root := range []string{"app", "cmd", "config"} {
		scanDir(root, keys, &r)
	}

	printSection("Keys used in code but absent from every locale (rendered raw)", r.missing)
	printSection("Hard-coded text candidates (review: some are technical)", r.hardcoded)
	printSection("Raw errors (review: do they reach the player?)", r.raw)
}

func loadLocales(root string) (map[string][]string, []string) {
	keys := map[string][]string{}
	var locales []string

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".yaml" {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "read error:", path, err)
			return nil
		}
		messages := map[string]any{}
		if err := yaml.Unmarshal(content, &messages); err != nil {
			fmt.Fprintln(os.Stderr, "parse error:", path, err)
			return nil
		}
		locale := strings.TrimSuffix(filepath.Base(path), ".yaml")
		if !slices.Contains(locales, locale) {
			locales = append(locales, locale)
		}
		for key := range messages {
			if !slices.Contains(keys[key], locale) {
				keys[key] = append(keys[key], locale)
			}
		}
		return nil
	})

	return keys, locales
}

func printLocaleGaps(keys map[string][]string, locales []string) {
	var lines []string
	for key, found := range keys {
		if len(found) != len(locales) {
			lines = append(lines, fmt.Sprintf("  %s: only in %s", key, strings.Join(found, ", ")))
		}
	}
	slices.Sort(lines)
	fmt.Printf("== Keys missing from a locale (%d)\n%s\n\n", len(lines), strings.Join(lines, "\n"))
}

func printSection(title string, findings []finding) {
	fmt.Printf("== %s (%d)\n", title, len(findings))
	for _, f := range findings {
		fmt.Printf("  %s\t%s\t%q\n", f.pos, f.call, f.value)
	}
	fmt.Println()
}

func scanDir(root string, keys map[string][]string, r *report) {
	fset := token.NewFileSet()
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, "parse error:", path, err)
			return nil
		}
		scanFile(fset, file, keys, r)
		return nil
	})
}

func scanFile(fset *token.FileSet, file *ast.File, keys map[string][]string, r *report) {
	var stack []ast.Node
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		// Inspect sends no closing nil when the callback returns false: do not push these.
		switch n.(type) {
		case *ast.ImportSpec, *ast.Field:
			return false
		}
		stack = append(stack, n)

		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			classify(fset, lit, stack, keys, r)
		}
		return true
	})
}

func classify(fset *token.FileSet, lit *ast.BasicLit, stack []ast.Node, keys map[string][]string, r *report) {
	value, err := strconv.Unquote(lit.Value)
	if err != nil || !hasText(value) {
		return
	}

	parent := stack[len(stack)-2]
	// "stat."+name builds its key at runtime, so there is nothing to look up.
	if _, ok := parent.(*ast.BinaryExpr); ok {
		return
	}

	field := ""
	if kv, ok := parent.(*ast.KeyValueExpr); ok {
		if kv.Key == ast.Expr(lit) {
			return
		}
		if id, ok := kv.Key.(*ast.Ident); ok {
			field = id.Name
		}
	}

	call := enclosingCall(stack)
	found := finding{fset.Position(lit.Pos()).String(), call, value}

	switch {
	case technical[call]:
	case rawErrors[call]:
		r.raw = append(r.raw, found)
	case translated[call] || (labelled[call] && (field == "Label" || field == "ID")):
		if _, ok := keys[value]; !ok && !strings.Contains(value, "%") {
			r.missing = append(r.missing, found)
		}
	case looksLikeText(value):
		r.hardcoded = append(r.hardcoded, found)
	}
}

func hasText(value string) bool {
	return strings.ContainsFunc(formatVerb.ReplaceAllString(value, ""), unicode.IsLetter)
}

// looksLikeText keeps sentences and capitalised words, leaving ids such as "field.name" out.
func looksLikeText(value string) bool {
	text := strings.TrimSpace(formatVerb.ReplaceAllString(value, ""))
	if strings.Contains(text, " ") {
		return true
	}
	return unicode.IsUpper([]rune(text)[0])
}

func enclosingCall(stack []ast.Node) string {
	for i := len(stack) - 2; i >= 0; i-- {
		call, ok := stack[i].(*ast.CallExpr)
		if !ok {
			continue
		}
		switch fn := call.Fun.(type) {
		case *ast.SelectorExpr:
			if x, ok := fn.X.(*ast.Ident); ok {
				return x.Name + "." + fn.Sel.Name
			}
			return "?." + fn.Sel.Name
		case *ast.Ident:
			return fn.Name
		}
		return "?"
	}
	return "-"
}
