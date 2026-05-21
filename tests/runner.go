package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/metamessage/metamessage/internal/core"
	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/jsonc"
)

type FixtureResult struct {
	Path           string
	Category       string
	Name           string
	RoundTrip      bool
	Reversibility  bool
	Error          string
	OrigJSON       string
	RePrintedJSONC string
}

func main() {
	fixturesDir := "tests/fixtures"
	if len(os.Args) > 1 {
		fixturesDir = os.Args[1]
	}

	var results []FixtureResult
	passed := 0
	failed := 0

	err := filepath.Walk(fixturesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".jsonc") {
			return nil
		}

		relPath, _ := filepath.Rel(fixturesDir, path)
		parts := strings.SplitN(relPath, string(filepath.Separator), 2)
		category := parts[0]
		name := filepath.Base(path)

		result := testFixture(path)
		result.Path = relPath
		result.Category = category
		result.Name = name
		results = append(results, result)

		if result.RoundTrip && result.Reversibility {
			passed++
			fmt.Printf("  PASS  %s\n", relPath)
		} else {
			failed++
			fmt.Printf("  FAIL  %s: %s\n", relPath, result.Error)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking fixtures: %v\n", err)
		os.Exit(1)
	}

	// Summary
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total: %d, Passed: %d, Failed: %d\n\n", len(results), passed, failed)

	// Print failures detail
	for _, r := range results {
		if !r.RoundTrip || !r.Reversibility {
			fmt.Printf("--- FAIL: %s ---\n", r.Path)
			fmt.Printf("Error: %s\n", r.Error)
			if r.OrigJSON != r.RePrintedJSONC {
				fmt.Printf("Original JSON:\n%s\n\n", r.OrigJSON)
				fmt.Printf("Re-printed JSONC:\n%s\n\n", r.RePrintedJSONC)
			}
			fmt.Println()
		}
	}

	if failed > 0 {
		os.Exit(1)
	}
}

func testFixture(path string) FixtureResult {
	result := FixtureResult{}

	data, err := os.ReadFile(path)
	if err != nil {
		result.Error = fmt.Sprintf("read error: %v", err)
		return result
	}

	input := string(data)

	// Step 1: Parse JSONC → IR
	node1, err := core.ParseFromJSONC(input)
	if err != nil {
		result.Error = fmt.Sprintf("parse error: %v", err)
		return result
	}

	// Step 2: Print IR → JSONC
	reprinted := jsonc.ToJSONC(node1)
	result.RePrintedJSONC = reprinted

	// Step 3: Parse the re-printed JSONC → IR2
	node2, err := core.ParseFromJSONC(reprinted)
	if err != nil {
		result.Error = fmt.Sprintf("re-parse error: %v", err)
		return result
	}

	// Step 4: Compare IR1 == IR2 (structural equality)
	if !nodesEqual(node1, node2) {
		result.Error = "IR mismatch after round-trip"
		return result
	}
	result.RoundTrip = true

	// Step 5: Normalize and compare JSON content (reversibility)
	origJSON := normalizeJSON(input)
	reprintedJSON := normalizeJSON(reprinted)
	result.OrigJSON = origJSON
	result.RePrintedJSONC = reprintedJSON

	if origJSON == reprintedJSON {
		result.Reversibility = true
	} else {
		// Check if they're semantically equivalent (same JSON structure, ignoring comments)
		if jsonSemanticallyEqual(origJSON, reprintedJSON) {
			result.Reversibility = true
		} else {
			result.Error = fmt.Sprintf("JSON not reversible\n  original:  %s\n  reprinted: %s", origJSON, reprintedJSON)
		}
	}

	return result
}

// normalizeJSON removes comments and whitespace for comparison
var (
	lineCommentRe  = regexp.MustCompile(`//[^\n]*`)
	blockCommentRe = regexp.MustCompile(`/\*[\s\S]*?\*/`)
	whitespaceRe   = regexp.MustCompile(`\s+`)
)

func normalizeJSON(s string) string {
	// Remove comments
	s = lineCommentRe.ReplaceAllString(s, "")
	s = blockCommentRe.ReplaceAllString(s, "")

	// Remove trailing commas before ] or }
	s = regexp.MustCompile(`,(\s*[}\]])`).ReplaceAllString(s, "$1")

	// Compact whitespace (but preserve inside strings)
	var result strings.Builder
	inString := false
	escaped := false
	for _, ch := range s {
		if escaped {
			result.WriteRune(ch)
			escaped = false
			continue
		}
		if ch == '\\' && inString {
			result.WriteRune(ch)
			escaped = true
			continue
		}
		if ch == '"' {
			inString = !inString
			result.WriteRune(ch)
			continue
		}
		if inString {
			result.WriteRune(ch)
			continue
		}
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			continue
		}
		result.WriteRune(ch)
	}
	return result.String()
}

// jsonSemanticallyEqual compares two JSON strings semantically
func jsonSemanticallyEqual(a, b string) bool {
	var va, vb interface{}
	if err := json.Unmarshal([]byte(a), &va); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &vb); err != nil {
		return false
	}
	return reflect.DeepEqual(va, vb)
}

// nodesEqual recursively compares two IR nodes for structural equality
func nodesEqual(a, b ir.Node) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.GetType() != b.GetType() {
		return false
	}

	switch a.GetType() {
	case ir.NodeTypeValue:
		va := a.(*ir.Value)
		vb, ok := b.(*ir.Value)
		if !ok {
			return false
		}
		if va.Text != vb.Text || va.Data != vb.Data {
			return false
		}
		return tagsEqual(va.Tag, vb.Tag)

	case ir.NodeTypeObject:
		oa := a.(*ir.Object)
		ob, ok := b.(*ir.Object)
		if !ok {
			return false
		}
		if len(oa.Fields) != len(ob.Fields) {
			return false
		}
		for i := range oa.Fields {
			if oa.Fields[i].Key != ob.Fields[i].Key {
				return false
			}
			if !nodesEqual(oa.Fields[i].Value, ob.Fields[i].Value) {
				return false
			}
		}
		return tagsEqual(oa.Tag, ob.Tag)

	case ir.NodeTypeDoc:
		da := a.(*ir.Doc)
		db, ok := b.(*ir.Doc)
		if !ok {
			return false
		}
		if len(da.Fields) != len(db.Fields) {
			return false
		}
		for i := range da.Fields {
			if da.Fields[i].Key != db.Fields[i].Key {
				return false
			}
			if !nodesEqual(da.Fields[i].Value, db.Fields[i].Value) {
				return false
			}
		}
		return tagsEqual(da.Tag, db.Tag)

	case ir.NodeTypeArray:
		aa := a.(*ir.Array)
		ab, ok := b.(*ir.Array)
		if !ok {
			return false
		}
		if len(aa.Items) != len(ab.Items) {
			return false
		}
		for i := range aa.Items {
			if !nodesEqual(aa.Items[i], ab.Items[i]) {
				return false
			}
		}
		return tagsEqual(aa.Tag, ab.Tag)
	}

	return false
}

func tagsEqual(a, b *ir.Tag) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.ToString() == b.ToString()
}