package main

type LangConfig struct {
	LanguageName  string
	Extension     string
	TestExtension string
	CompileCmd    []string
	TestCmd       []string
}

var SupportedLanguages = map[string]LangConfig{
	"go": {
		LanguageName:  "Go",
		Extension:     ".go",
		TestExtension: "_test.go",
		CompileCmd:    []string{"go", "build", "./..."},
		TestCmd:       []string{"go", "test", "./..."},
	},
	"rust": {
		LanguageName:  "Rust",
		Extension:     ".rs",
		TestExtension: "_test.rs",
		CompileCmd:    []string{"cargo", "check"},
		TestCmd:       []string{"cargo", "test"},
	},
	"node": {
		LanguageName:  "TypeScript/Node.js",
		Extension:     ".ts",
		TestExtension: ".test.ts",
		CompileCmd:    []string{"tsc", "--noEmit"},
		TestCmd:       []string{"npm", "test"},
	},
	"python": {
		LanguageName:  "Python",
		Extension:     ".py",
		TestExtension: "_test.py",
		CompileCmd:    []string{"python", "-m", "py_compile"},
		TestCmd:       []string{"pytest"},
	},
	"java": {
		LanguageName:  "Java",
		Extension:     ".java",
		TestExtension: "Test.java",
		CompileCmd:    []string{"mvn", "compile"},
		TestCmd:       []string{"mvn", "test"},
	},
	"csharp": {
		LanguageName:  "C# (.NET)",
		Extension:     ".cs",
		TestExtension: "Tests.cs",
		CompileCmd:    []string{"dotnet", "build"},
		TestCmd:       []string{"dotnet", "test"},
	},
	"ruby": {
		LanguageName:  "Ruby",
		Extension:     ".rb",
		TestExtension: "_spec.rb",
		CompileCmd:    []string{"ruby", "-c"},
		TestCmd:       []string{"rspec"},
	},
	"cpp": {
		LanguageName:  "C++",
		Extension:     ".cpp",
		TestExtension: "_test.cpp",
		CompileCmd:    []string{"g++", "-fsyntax-only"},
		TestCmd:       []string{"make", "test"},
	},
}
