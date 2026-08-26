package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// UIEvent is a generic message struct
type UIEvent struct {
	Type    string   `json:"type"`
	Message string   `json:"message,omitempty"`
	Files   []string `json:"files,omitempty"`
	File    string   `json:"file,omitempty"`
}

func emitLog(msg string) {
	evt, _ := json.Marshal(UIEvent{Type: "log", Message: msg})
	BroadcastEvent(string(evt))
}

func RunBatchMigrationUI(targetLang string, instructions string) {
	defer func() {
		if r := recover(); r != nil {
			emitLog(fmt.Sprintf("[FATAL]: Engine crashed: %v", r))
			BroadcastEvent("DONE")
		}
	}()

	inputDir := "../legacy_app"
	outputDir := "../modern_app"
	modelName := "phi3" // Could be dynamic from UI

	// Clean output directory before batch run
	os.RemoveAll(outputDir)
	os.MkdirAll(outputDir, 0755)
	if targetLang == "go" {
		cmd := exec.Command("go", "mod", "init", "modern_app")
		cmd.Dir = outputDir
		cmd.Run()
	}

	emitLog("==================================================")
	emitLog(" DEEP TECH MIGRATION BATCH PROCESSOR (AST + LLM)  ")
	emitLog("==================================================")

	langConfig, exists := SupportedLanguages[strings.ToLower(targetLang)]
	if !exists {
		emitLog(fmt.Sprintf("[FATAL]: Unsupported target language: %s", targetLang))
		BroadcastEvent("DONE")
		return
	}

	// 1. Crawl Directory
	emitLog("[Batch Processor]: Scanning legacy codebase...")
	files, err := FindAllLegacyFiles(inputDir)
	if err != nil {
		emitLog(fmt.Sprintf("[FATAL]: Failed to scan directory: %v", err))
		BroadcastEvent("DONE")
		return
	}

	// 2. Build Global Graph
	emitLog(fmt.Sprintf("[Batch Processor]: Found %d Legacy files. Building Global AST Graph...", len(files)))
	globalGraph := make(map[string][]string)
	astDetails := make(map[string]*ASTGraphResult)

	for _, file := range files {
		absPath := filepath.Join(inputDir, file)
		graph, err := BuildDependencyGraph(absPath, inputDir)
		if err != nil {
			emitLog(fmt.Sprintf("[FATAL]: AST parsing failed for %s: %v", file, err))
			BroadcastEvent("DONE")
			return
		}
		globalGraph[file] = graph.Dependencies
		astDetails[file] = graph
	}

	// 3. Topological Sort
	emitLog("[Batch Processor]: Performing Topological Sort...")
	sortedFiles, err := TopologicalSort(globalGraph)
	if err != nil {
		emitLog(fmt.Sprintf("[FATAL]: Sorting failed: %v", err))
		BroadcastEvent("DONE")
		return
	}

	// Emit graph event
	graphEvt, _ := json.Marshal(UIEvent{Type: "graph", Files: sortedFiles})
	BroadcastEvent(string(graphEvt))

	// 4. Batch Process
	for _, relPath := range sortedFiles {
		emitLog(fmt.Sprintf("\n>>> MIGRATING: %s <<<", relPath))
		
		evt, _ := json.Marshal(UIEvent{Type: "file_start", File: relPath})
		BroadcastEvent(string(evt))

		absTargetFile := filepath.Join(inputDir, relPath)
		graph := astDetails[relPath]
		if graph == nil {
			emitLog(fmt.Sprintf("[Warning]: Skipping %s (Invalid file or external dependency)", relPath))
			continue
		}

		bundledContext, err := BundleContext(absTargetFile, graph, astDetails, inputDir)
		if err != nil {
			emitLog(fmt.Sprintf("[FATAL]: Failed to bundle context: %v", err))
			evt, _ := json.Marshal(UIEvent{Type: "file_error", File: relPath})
			BroadcastEvent(string(evt))
			BroadcastEvent("DONE")
			return
		}

		emitLog("=== STARTING AUTONOMOUS AGENTIC LOOP ===")
		maxRetries := 3
		feedbackError := ""
		var finalCode string

		for attempt := 1; attempt <= maxRetries; attempt++ {
			emitLog(fmt.Sprintf("--- Agent Attempt %d ---", attempt))
			
			evtAtt, _ := json.Marshal(UIEvent{Type: "attempt"})
			BroadcastEvent(string(evtAtt))

			migratedCode, testCode, err := CallLLM(bundledContext, feedbackError, modelName, langConfig, instructions)
			if err != nil {
				emitLog(fmt.Sprintf("[FATAL]: LLM Translation failed: %v", err))
				BroadcastEvent("DONE")
				return
			}
			if strings.TrimSpace(migratedCode) == "" {
				emitLog("[Compile Failed]: LLM returned completely empty code. The model might be failing or context is too large.")
				
				if attempt == maxRetries {
					emitLog(fmt.Sprintf("[FATAL]: LLM failed to fix logic in %s after %d attempts.", relPath, maxRetries))
					evtErr, _ := json.Marshal(UIEvent{Type: "file_error", File: relPath})
					BroadcastEvent(string(evtErr))
				}
				continue
			}

			os.MkdirAll(outputDir, 0755)
			
			outputFile := filepath.Join(outputDir, "migrated_"+filepath.Base(absTargetFile))
			outputFile = strings.Replace(outputFile, ".php", langConfig.Extension, 1)
			outputFile = strings.Replace(outputFile, ".py", langConfig.Extension, 1)
			outputFile = strings.Replace(outputFile, ".dart", langConfig.Extension, 1)

			testFile := strings.Replace(outputFile, langConfig.Extension, langConfig.TestExtension, 1)

			os.WriteFile(outputFile, []byte(migratedCode), 0644)
			os.WriteFile(testFile, []byte(testCode), 0644)

			compileErr := CompileCode(outputDir, filepath.Base(outputFile), langConfig)
			if compileErr == nil {
				testErr := TestCode(outputDir, filepath.Base(testFile), langConfig)
				if testErr == nil {
					finalCode = migratedCode
					emitLog(fmt.Sprintf("[Orchestrator]: SUCCESS! Code and Logic verified. Saved to: %s", outputFile))
					
					evtSucc, _ := json.Marshal(UIEvent{Type: "file_success", File: relPath})
					BroadcastEvent(string(evtSucc))

					AppendToDataset(bundledContext, finalCode)
					break
				} else {
					feedbackError = testErr.Error()
					emitLog(fmt.Sprintf("[Test Failed]: %s", feedbackError))
					
					os.Remove(outputFile)
					os.Remove(testFile)

					if attempt == maxRetries {
						emitLog(fmt.Sprintf("[FATAL]: LLM failed to fix logic in %s after %d attempts.", relPath, maxRetries))
						evtErr, _ := json.Marshal(UIEvent{Type: "file_error", File: relPath})
						BroadcastEvent(string(evtErr))
					}
				}
			} else {
				feedbackError = compileErr.Error()
				emitLog(fmt.Sprintf("[Compile Failed]: %s", feedbackError))
				
				os.Remove(outputFile)
				os.Remove(testFile)

				if attempt == maxRetries {
					emitLog(fmt.Sprintf("[FATAL]: LLM failed to fix syntax in %s after %d attempts.", relPath, maxRetries))
					evtErr, _ := json.Marshal(UIEvent{Type: "file_error", File: relPath})
					BroadcastEvent(string(evtErr))
				}
			}
		}
	}

	emitLog("\n==================================================")
	emitLog(" BATCH MIGRATION COMPLETE! ALL FILES PROCESSED!   ")
	emitLog("==================================================")
	BroadcastEvent("DONE")
}
