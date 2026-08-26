package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	inputDir     string
	outputDir    string
	modelName    string
	targetLang   string
	instructions string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the migration factory on a legacy codebase",
	Long:  `Run the batch migration processor to convert legacy code to Go using AI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("==================================================")
		fmt.Println(" DEEP TECH MIGRATION BATCH PROCESSOR (AST + LLM)  ")
		fmt.Println("==================================================")
		fmt.Printf("Input Directory: %s\n", inputDir)
		fmt.Printf("Output Directory: %s\n", outputDir)
		fmt.Printf("Model: %s\n", modelName)
		fmt.Printf("Target Language: %s\n", targetLang)
		fmt.Println("==================================================")

		langConfig, exists := SupportedLanguages[strings.ToLower(targetLang)]
		if !exists {
			log.Fatalf("[FATAL]: Unsupported target language: %s", targetLang)
		}

		// 1. Crawl Directory
		fmt.Println("[Batch Processor]: Scanning legacy codebase...")
		files, err := FindAllLegacyFiles(inputDir)
		if err != nil {
			log.Fatalf("[FATAL]: Failed to scan directory: %v", err)
		}

		// 2. Build Global Graph
		fmt.Printf("[Batch Processor]: Found %d Legacy files. Building Global AST Graph...\n", len(files))
		globalGraph := make(map[string][]string)
		astDetails := make(map[string]*ASTGraphResult)

		for _, file := range files {
			absPath := filepath.Join(inputDir, file)
			graph, err := BuildDependencyGraph(absPath, inputDir)
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
		os.RemoveAll(outputDir)
		os.MkdirAll(outputDir, 0755)
		if targetLang == "go" {
			cmd := exec.Command("go", "mod", "init", "modern_app")
			cmd.Dir = outputDir
			cmd.Run()
		}

		for _, relPath := range sortedFiles {
			fmt.Printf("\n>>> MIGRATING: %s <<<\n", relPath)
			absTargetFile := filepath.Join(inputDir, relPath)
			graph := astDetails[relPath]
			if graph == nil {
				fmt.Printf("[Warning]: Skipping %s (Invalid file or external dependency)\n", relPath)
				continue
			}

			bundledContext, err := BundleContext(absTargetFile, graph, astDetails, inputDir)
			if err != nil {
				log.Fatalf("[FATAL]: Failed to bundle context: %v", err)
			}

			fmt.Println("=== STARTING AUTONOMOUS AGENTIC LOOP ===")
			maxRetries := 3
			feedbackError := ""
			var finalCode string

			for attempt := 1; attempt <= maxRetries; attempt++ {
				fmt.Printf("--- Agent Attempt %d ---\n", attempt)
				migratedCode, testCode, err := CallLLM(bundledContext, feedbackError, modelName, langConfig, instructions)
				if err != nil {
					log.Fatalf("[FATAL]: LLM Translation failed: %v", err)
				}
				if strings.TrimSpace(migratedCode) == "" {
					fmt.Println("[Compile Failed]: LLM returned completely empty code. The model might be failing or context is too large.")
					if attempt == maxRetries {
						log.Fatalf("[FATAL]: LLM failed to fix logic in %s after %d attempts. Last Error: Empty Code", relPath, maxRetries)
					}
					continue
				}

				// Clean old files before compiling the new one to prevent conflicts
				os.MkdirAll(outputDir, 0755)
				
				outputFile := filepath.Join(outputDir, "migrated_"+filepath.Base(absTargetFile))
				outputFile = strings.Replace(outputFile, ".php", langConfig.Extension, 1)
				outputFile = strings.Replace(outputFile, ".py", langConfig.Extension, 1) // For Python files

				testFile := strings.Replace(outputFile, langConfig.Extension, langConfig.TestExtension, 1)

				err = os.WriteFile(outputFile, []byte(migratedCode), 0644)
				if err != nil {
					log.Fatalf("[FATAL]: Failed to save migrated code: %v", err)
				}

				err = os.WriteFile(testFile, []byte(testCode), 0644)
				if err != nil {
					log.Fatalf("[FATAL]: Failed to save test code: %v", err)
				}

				// 1. Check Syntax (Compile)
				compileErr := CompileCode(outputDir, langConfig)
				if compileErr == nil {
					// 2. Check Logic (Test)
					testErr := TestCode(outputDir, langConfig)
					if testErr == nil {
						finalCode = migratedCode
						fmt.Printf("[Orchestrator]: SUCCESS! Code and Logic verified. Saved to: %s\n", outputFile)
						
						fmt.Println("=== PHASE 4: DATA HARVESTING ===")
						if err := AppendToDataset(bundledContext, finalCode); err != nil {
							fmt.Printf("[Warning]: Failed to save to dataset: %v\n", err)
						}
						break
					} else {
						feedbackError = testErr.Error()
						
						os.Remove(outputFile)
						os.Remove(testFile)

						if attempt == maxRetries {
							log.Fatalf("[FATAL]: LLM failed to fix logic in %s after %d attempts. Last Error: %s", relPath, maxRetries, feedbackError)
						}
					}
				} else {
					feedbackError = compileErr.Error()
					
					os.Remove(outputFile)
					os.Remove(testFile)

					if attempt == maxRetries {
						log.Fatalf("[FATAL]: LLM failed to fix syntax in %s after %d attempts. Last Error: %s", relPath, maxRetries, feedbackError)
					}
				}
			}
		}

		fmt.Println("\n==================================================")
		fmt.Println(" BATCH MIGRATION COMPLETE! ALL FILES PROCESSED!   ")
		fmt.Println("==================================================")
	},
}

func init() {
	runCmd.Flags().StringVarP(&inputDir, "input", "i", "../legacy_app", "Path to legacy codebase")
	runCmd.Flags().StringVarP(&outputDir, "output", "o", "../modern_app", "Path to save modern code")
	runCmd.Flags().StringVarP(&modelName, "model", "m", "phi3", "Ollama model name to use")
	runCmd.Flags().StringVarP(&targetLang, "target", "t", "go", "Target programming language (go, rust, node, python, java, csharp, ruby, cpp)")
	runCmd.Flags().StringVarP(&instructions, "instructions", "p", "", "Custom instructions for the AI (e.g. Version upgrades, specific frameworks)")
}
