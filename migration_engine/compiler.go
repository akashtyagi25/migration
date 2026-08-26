package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// CompileCode runs the compiler for the target language.
// Returns an error containing the stderr output if compilation fails.
func CompileCode(targetDir string, langConfig LangConfig) error {
	fmt.Printf("[Compiler Engine]: Running `%s` to verify syntax...\n", strings.Join(langConfig.CompileCmd, " "))
	
	cmd := exec.Command(langConfig.CompileCmd[0], langConfig.CompileCmd[1:]...)
	cmd.Dir = targetDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		errorMessage := strings.TrimSpace(string(output))
		if errorMessage == "" {
			errorMessage = err.Error()
		}
		fmt.Printf("[Compiler Engine]: SYNTAX ERROR DETECTED:\n%s\n", errorMessage)
		return fmt.Errorf("%s", errorMessage)
	}

	fmt.Println("[Compiler Engine]: Code Compiled Successfully! No Syntax Errors.")
	return nil
}

// TestCode runs tests for the target language and returns any logic errors
func TestCode(dir string, langConfig LangConfig) error {
	fmt.Printf("[Compiler Engine]: Running `%s` to verify business logic...\n", strings.Join(langConfig.TestCmd, " "))

	cmd := exec.Command(langConfig.TestCmd[0], langConfig.TestCmd[1:]...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("[Compiler Engine]: LOGIC ERROR DETECTED:")
		fmt.Println(string(output))
		return fmt.Errorf("%s", string(output))
	}

	fmt.Println("[Compiler Engine]: Tests Passed Successfully! Business Logic is 100% correct.")
	return nil
}

