package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// CompileCode runs the compiler for the target language.
// Returns an error containing the stderr output if compilation fails.
func CompileCode(targetDir string, fileName string, langConfig LangConfig) error {
	var args []string
	for _, arg := range langConfig.CompileCmd[1:] {
		args = append(args, strings.ReplaceAll(arg, "{file}", fileName))
	}

	fmt.Printf("[Compiler Engine]: Running `%s %s` to verify syntax...\n", langConfig.CompileCmd[0], strings.Join(args, " "))
	
	cmd := exec.Command(langConfig.CompileCmd[0], args...)
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
func TestCode(dir string, testFileName string, langConfig LangConfig) error {
	var args []string
	for _, arg := range langConfig.TestCmd[1:] {
		args = append(args, strings.ReplaceAll(arg, "{file}", testFileName))
	}

	fmt.Printf("[Compiler Engine]: Running `%s %s` to verify business logic...\n", langConfig.TestCmd[0], strings.Join(args, " "))

	cmd := exec.Command(langConfig.TestCmd[0], args...)
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

