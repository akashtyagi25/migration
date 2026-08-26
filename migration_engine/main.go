package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Symbol represents a class or function found in the AST
type Symbol struct {
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Methods []string `json:"methods"`
}

// ASTGraphResult represents the JSON output from the PHP script
type ASTGraphResult struct {
	File            string   `json:"file"`
	Dependencies    []string `json:"dependencies"`
	ExportedSymbols []Symbol `json:"exported_symbols"`
	AstNodesParsed  int      `json:"ast_nodes_parsed"`
	Error           string   `json:"error,omitempty"`
}

// BuildDependencyGraph routes the file to the correct AST parser based on its extension
func BuildDependencyGraph(targetFile string, baseDir string) (*ASTGraphResult, error) {
	ext := strings.ToLower(filepath.Ext(targetFile))
	var cmd *exec.Cmd
	absTarget, _ := filepath.Abs(targetFile)

	if ext == ".php" {
		parserScriptPath := "../parsers/php_parser/ast_grapher.php"
		fmt.Printf("[Universal Engine]: Invoking PHP AST parser on %s...\n", targetFile)
		cmd = exec.Command("php", filepath.Base(parserScriptPath), absTarget)
		cmd.Dir = filepath.Dir(parserScriptPath)
	} else if ext == ".py" {
		parserScriptPath := "../parsers/python_parser/ast_grapher.py"
		fmt.Printf("[Universal Engine]: Invoking Python AST parser on %s...\n", targetFile)
		cmd = exec.Command("python", filepath.Base(parserScriptPath), absTarget)
		cmd.Dir = filepath.Dir(parserScriptPath)
	} else {
		// Use Hybrid Regex Parser for all other languages
		return ExtractSymbolsUniversal(targetFile, baseDir)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("AST execution failed: %v, Output: %s", err, string(output))
	}

	var result ASTGraphResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse AST JSON: %v, Output: %s", err, string(output))
	}

	if result.Error != "" {
		return nil, fmt.Errorf("AST Parser Error: %s", result.Error)
	}

	return &result, nil
}

func main() {
	Execute()
}
