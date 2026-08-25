package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// CompileCode runs `go build` in the specified directory.
// Returns an error containing the stderr output if compilation fails.
func CompileCode(targetDir string) error {
	fmt.Println("[Compiler Engine]: Running `go build` to verify syntax...")
	
	cmd := exec.Command("go", "build", ".")
	cmd.Dir = targetDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		errorMessage := strings.TrimSpace(string(output))
		fmt.Printf("[Compiler Engine]: SYNTAX ERROR DETECTED:\n%s\n", errorMessage)
		return fmt.Errorf("%s", errorMessage)
	}

	fmt.Println("[Compiler Engine]: Code Compiled Successfully! No Syntax Errors.")
	return nil
}
