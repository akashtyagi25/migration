package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindAllLegacyFiles recursively finds all supported source files in a given directory
func FindAllLegacyFiles(rootDir string) ([]string, error) {
	var files []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(info.Name()))
			
			// A list of common enterprise languages
			validExts := map[string]bool{
				".php": true, ".py": true, ".java": true, ".cs": true, 
				".ts": true, ".js": true, ".rb": true, ".dart": true, 
				".cpp": true, ".c": true, ".go": true, ".rs": true,
			}

			if validExts[ext] {
			// Convert to relative path for cleaner graph keys
			relPath, _ := filepath.Rel(rootDir, path)
			// Ensure forward slashes for consistency
			relPath = strings.ReplaceAll(relPath, "\\", "/")
				files = append(files, relPath)
			}
		}
		return nil
	})
	return files, err
}

// TopologicalSort performs DFS on the global dependency graph to determine migration order.
// Independent files (no dependencies) will be at the start of the sorted list.
func TopologicalSort(globalGraph map[string][]string) ([]string, error) {
	var sorted []string
	visited := make(map[string]bool)
	tempMark := make(map[string]bool)

	var visit func(node string) error
	visit = func(node string) error {
		if tempMark[node] {
			return fmt.Errorf("circular dependency detected involving %s", node)
		}
		if !visited[node] {
			tempMark[node] = true
			for _, dep := range globalGraph[node] {
				if err := visit(dep); err != nil {
					return err
				}
			}
			visited[node] = true
			tempMark[node] = false
			sorted = append(sorted, node) // Append after dependencies are processed
		}
		return nil
	}

	for node := range globalGraph {
		if !visited[node] {
			if err := visit(node); err != nil {
				return nil, err
			}
		}
	}

	return sorted, nil
}
