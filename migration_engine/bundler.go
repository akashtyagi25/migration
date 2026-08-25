package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BundleContext reads the target file and injects only the AST-extracted stubs of its dependencies
// into a single Mega-Prompt string. This is Smart Chunking (RAG).
func BundleContext(targetFile string, graph *ASTGraphResult, astDetails map[string]*ASTGraphResult, baseDir string) (string, error) {
	fmt.Printf("[Bundler]: Assembling RAG context for %s...\n", targetFile)

	contextString := "=== LEGACY SYSTEM CONTEXT (AST EXTRACTED) ===\n\n"

	// 1. Inject Dependency Stubs (Instead of raw code)
	for _, dep := range graph.Dependencies {
		depGraph, exists := astDetails[dep]
		if !exists {
			fmt.Printf("[Bundler Warning]: No AST data for dependency %s\n", dep)
			contextString += fmt.Sprintf("--- Dependency Summary: %s (FILE NOT FOUND) ---\n\n", dep)
			continue
		}

		contextString += fmt.Sprintf("--- Dependency Summary: %s ---\n", dep)
		contextString += "Exported Symbols:\n"
		
		if len(depGraph.ExportedSymbols) == 0 {
			contextString += " - No classes or functions found.\n"
		} else {
			for _, sym := range depGraph.ExportedSymbols {
				if sym.Type == "class" {
					methods := strings.Join(sym.Methods, ", ")
					contextString += fmt.Sprintf(" - Class: %s (Methods: %s)\n", sym.Name, methods)
				} else if sym.Type == "function" {
					contextString += fmt.Sprintf(" - Function: %s\n", sym.Name)
				}
			}
		}
		contextString += "\n"
	}

	// 2. Read the main target file (We need the full code for the file we are actually migrating)
	mainContent, err := os.ReadFile(targetFile)
	if err != nil {
		return "", fmt.Errorf("failed to read main target file: %v", err)
	}

	contextString += "=== MAIN FILE TO MIGRATE ===\n"
	contextString += fmt.Sprintf("--- File: %s ---\n%s\n", filepath.Base(targetFile), string(mainContent))

	fmt.Println("[Bundler]: Smart context bundling complete.")
	return contextString, nil
}
