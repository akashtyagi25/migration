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
}
