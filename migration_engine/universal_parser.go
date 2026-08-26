package main

import (
	"fmt"
	"os"
	"regexp"
)

// ExtractSymbolsUniversal uses Regex heuristics to build a pseudo-AST for unsupported languages.
func ExtractSymbolsUniversal(targetFile string, baseDir string) (*ASTGraphResult, error) {
	fmt.Printf("[Universal Engine]: Invoking Regex Fallback Parser on %s...\n", targetFile)

	contentBytes, err := os.ReadFile(targetFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", targetFile, err)
	}
	content := string(contentBytes)

	// Broad regex for imports (e.g. import "os", using System, require 'json', include <stdio.h>)
	importRegex := regexp.MustCompile(`(?i)(?:import|using|require|include|from)\s+['"<]?([a-zA-Z0-9_\.\/\\]+)['">]?`)
	
	// Broad regex for classes/structs (e.g. class User, struct Point, interface DB)
	classRegex := regexp.MustCompile(`(?i)(?:class|struct|interface|trait|enum)\s+([a-zA-Z0-9_]+)`)
	
	// Broad regex for functions (e.g. func main(), def execute(), public void doSomething())
	funcRegex := regexp.MustCompile(`(?i)(?:func|def|function|void|int|string|bool)\s+([a-zA-Z0-9_]+)\s*\(`)

	var dependencies []string
	depMap := make(map[string]bool)
	
	importMatches := importRegex.FindAllStringSubmatch(content, -1)
	for _, match := range importMatches {
		if len(match) > 1 {
			dep := match[1]
			if !depMap[dep] {
				depMap[dep] = true
				dependencies = append(dependencies, dep)
			}
		}
	}

	var symbols []Symbol
	classMatches := classRegex.FindAllStringSubmatch(content, -1)
	for _, match := range classMatches {
		if len(match) > 1 {
			symbols = append(symbols, Symbol{
				Type: "Class/Struct",
				Name: match[1],
			})
		}
	}

	funcMatches := funcRegex.FindAllStringSubmatch(content, -1)
	for _, match := range funcMatches {
		if len(match) > 1 {
			// Basic filtering to avoid matching standard control structures like if(), while()
			name := match[1]
			if name != "if" && name != "while" && name != "for" && name != "switch" && name != "catch" {
				symbols = append(symbols, Symbol{
					Type: "Function",
					Name: name,
				})
			}
		}
	}

	return &ASTGraphResult{
		File:            targetFile,
		Dependencies:    dependencies,
		ExportedSymbols: symbols,
		AstNodesParsed:  len(symbols),
	}, nil
}
