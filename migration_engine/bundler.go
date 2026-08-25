package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// BundleContext reads the target file and all dependencies returned by the AST parser,
// and concatenates them into a single Mega-Prompt string.
func BundleContext(targetFile string, graph *ASTGraphResult, baseDir string) (string, error) {
	fmt.Printf("[Bundler]: Assembling mega-context for %s...\n", targetFile)

	contextString := "=== LEGACY SYSTEM CONTEXT (AST EXTRACTED) ===\n\n"

	// Read dependencies first
	for _, dep := range graph.Dependencies {
		depPath := filepath.Join(baseDir, dep)
		content, err := os.ReadFile(depPath)
		if err != nil {
			fmt.Printf("[Bundler Warning]: Could not read dependency %s: %v\n", dep, err)
			contextString += fmt.Sprintf("--- Dependency: %s (FILE NOT FOUND) ---\n\n", dep)
			continue
		}
		contextString += fmt.Sprintf("--- Dependency: %s ---\n%s\n\n", dep, string(content))
	}

	// Read the main target file
	mainContent, err := os.ReadFile(targetFile)
	if err != nil {
		return "", fmt.Errorf("failed to read main target file: %v", err)
	}

	contextString += "=== MAIN FILE TO MIGRATE ===\n"
	contextString += fmt.Sprintf("--- File: %s ---\n%s\n", filepath.Base(targetFile), string(mainContent))

	fmt.Println("[Bundler]: Context bundling complete.")
	return contextString, nil
}
