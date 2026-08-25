package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DatasetEntry represents a single JSONL row for OpenAI / HuggingFace fine-tuning
type DatasetEntry struct {
	Instruction string `json:"instruction"`
	Input       string `json:"input"`
	Output      string `json:"output"`
}

// AppendToDataset saves the successful migration pair to dataset.jsonl
func AppendToDataset(bundledContext string, finalMigratedCode string) error {
	datasetDir := "../training_data"
	os.MkdirAll(datasetDir, 0755)
	
	datasetFile := filepath.Join(datasetDir, "dataset.jsonl")
	
	entry := DatasetEntry{
		Instruction: "Translate this legacy PHP file into modern, dependency-injected Go code.",
		Input:       bundledContext,
		Output:      finalMigratedCode,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal dataset entry: %v", err)
	}

	// Append to file
	f, err := os.OpenFile(datasetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open dataset file: %v", err)
	}
	defer f.Close()

	if _, err := f.Write(append(jsonData, '\n')); err != nil {
		return fmt.Errorf("failed to write to dataset file: %v", err)
	}

	fmt.Printf("[Data Harvester]: Added new perfect training pair to %s\n", datasetFile)
	return nil
}
