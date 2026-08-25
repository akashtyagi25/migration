package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

// BuildDependencyGraph runs the PHP AST parser and extracts dependencies deterministically.
func BuildDependencyGraph(targetFile string, parserScriptPath string) (*ASTGraphResult, error) {
	fmt.Printf("[Go Engine]: Invoking PHP AST parser on %s...\n", targetFile)

	cmd := exec.Command("php", parserScriptPath, targetFile)
	cmd.Dir = filepath.Dir(parserScriptPath)

	absTarget, _ := filepath.Abs(targetFile)
	cmd.Args = []string{"php", filepath.Base(parserScriptPath), absTarget}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("PHP AST execution failed: %v, Output: %s", err, string(output))
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
	fmt.Println("==================================================")
	fmt.Println(" DEEP TECH MIGRATION BATCH PROCESSOR (AST + LLM)  ")
	fmt.Println("==================================================")

	baseLegacyDir := "../legacy_app"
	parserScript := "../parsers/php_parser/ast_grapher.php"
	outputDir := "../modern_app"

	// 1. Crawl Directory
	fmt.Println("[Batch Processor]: Scanning legacy codebase...")
	files, err := FindAllPHPFiles(baseLegacyDir)
	if err != nil {
		log.Fatalf("[FATAL]: Failed to scan directory: %v", err)
	}

	// 2. Build Global Graph
	fmt.Printf("[Batch Processor]: Found %d PHP files. Building Global AST Graph...\n", len(files))
	globalGraph := make(map[string][]string)
	astDetails := make(map[string]*ASTGraphResult)

	for _, file := range files {
		absPath := filepath.Join(baseLegacyDir, file)
		graph, err := BuildDependencyGraph(absPath, parserScript)
		if err != nil {
			log.Fatalf("[FATAL]: AST parsing failed for %s: %v", file, err)
		}
		globalGraph[file] = graph.Dependencies
		astDetails[file] = graph
	}

	// 3. Topological Sort
	fmt.Println("[Batch Processor]: Performing Topological Sort...")
	sortedFiles, err := TopologicalSort(globalGraph)
	if err != nil {
		log.Fatalf("[FATAL]: Sorting failed: %v", err)
	}

	fmt.Println("\n=== MIGRATION ORDER (Dependencies First) ===")
	for i, file := range sortedFiles {
		fmt.Printf("%d. %s\n", i+1, file)
	}
	fmt.Println("============================================\n")

	// 4. Batch Process Each File
	for _, relPath := range sortedFiles {
		fmt.Printf("\n>>> MIGRATING: %s <<<\n", relPath)
		absTargetFile := filepath.Join(baseLegacyDir, relPath)
		graph := astDetails[relPath]

		bundledContext, err := BundleContext(absTargetFile, graph, astDetails, baseLegacyDir)
		if err != nil {
			log.Fatalf("[FATAL]: Failed to bundle context: %v", err)
		}

		fmt.Println("=== STARTING AUTONOMOUS AGENTIC LOOP ===")
		maxRetries := 3
		feedbackError := ""
		var finalCode string

		for attempt := 1; attempt <= maxRetries; attempt++ {
			fmt.Printf("--- Agent Attempt %d ---\n", attempt)
			
			migratedCode, err := CallLLM(bundledContext, feedbackError)
			if err != nil {
				log.Fatalf("[FATAL]: LLM Translation failed: %v", err)
			}

			// Clean old files before compiling the new one to prevent conflicts
			os.MkdirAll(outputDir, 0755)
			
			outputFile := filepath.Join(outputDir, "migrated_"+filepath.Base(absTargetFile))
			outputFile = strings.Replace(outputFile, ".php", ".go", 1)

			err = os.WriteFile(outputFile, []byte(migratedCode), 0644)
			if err != nil {
				log.Fatalf("[FATAL]: Failed to save migrated code: %v", err)
			}

			// Compile the generated code
			compileErr := CompileCode(outputDir)
			if compileErr == nil {
				finalCode = migratedCode
				fmt.Printf("[Orchestrator]: SUCCESS! Code is verified and saved to: %s\n", outputFile)
				
				fmt.Println("=== PHASE 4: DATA HARVESTING ===")
				if err := AppendToDataset(bundledContext, finalCode); err != nil {
					fmt.Printf("[Warning]: Failed to save to dataset: %v\n", err)
				}
				break
			} else {
				feedbackError = compileErr.Error()
				// Remove the broken file so it doesn't break subsequent compiles
				os.Remove(outputFile)
				if attempt == maxRetries {
					log.Fatalf("[FATAL]: LLM failed to fix %s after %d attempts. Last Error: %s", relPath, maxRetries, feedbackError)
				}
			}
		}
	}

	fmt.Println("\n==================================================")
	fmt.Println(" BATCH MIGRATION COMPLETE! ALL FILES PROCESSED!   ")
	fmt.Println("==================================================")
}
