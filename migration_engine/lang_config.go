package main

type LangConfig struct {
	LanguageName  string
	Extension     string
	TestExtension string
	CompileCmd    []string
	TestCmd       []string
	DummyTestCode string
}

var SupportedLanguages = map[string]LangConfig{
	"go": {
		LanguageName:  "Go",
		Extension:     ".go",
		TestExtension: "_test.go",
		CompileCmd:    []string{"go", "build", "./..."},
		TestCmd:       []string{"go", "test", "./..."},
		DummyTestCode: "package main\nimport \"testing\"\nfunc TestDummy(t *testing.T) {}\n",
	},
	"rust": {
		LanguageName:  "Rust",
		Extension:     ".rs",
		TestExtension: "_test.rs",
		CompileCmd:    []string{"rustc", "--emit=mir", "{file}"},
		TestCmd:       []string{"rustc", "--test", "{file}"},
		DummyTestCode: "#[test]\nfn test_dummy() {}\n",
	},
	"node": {
		LanguageName:  "TypeScript/Node.js",
		Extension:     ".ts",
		TestExtension: ".test.ts",
		CompileCmd:    []string{"npx", "tsc", "--noEmit", "{file}"},
		TestCmd:       []string{"npx", "ts-node", "{file}"},
		DummyTestCode: "console.log('Dummy test passed');\n",
	},
	"python": {
		LanguageName:  "Python",
		Extension:     ".py",
		TestExtension: "_test.py",
		CompileCmd:    []string{"python", "-m", "py_compile", "{file}"},
		TestCmd:       []string{"pytest", "{file}"},
		DummyTestCode: "def test_dummy():\n    pass\n",
	},
	"java": {
		LanguageName:  "Java",
		Extension:     ".java",
		TestExtension: "Test.java",
		CompileCmd:    []string{"javac", "{file}"},
		TestCmd:       []string{"java", "{file}"},
		DummyTestCode: "public class DummyTest { public static void main(String[] args) {} }\n",
	},
	"csharp": {
		LanguageName:  "C# (.NET)",
		Extension:     ".cs",
		TestExtension: "Tests.cs",
		CompileCmd:    []string{"csc", "/t:library", "{file}"},
		TestCmd:       []string{"dotnet", "test"},
		DummyTestCode: "using System;\nclass DummyTest { static void Main() {} }\n",
	},
	"ruby": {
		LanguageName:  "Ruby",
		Extension:     ".rb",
		TestExtension: "_spec.rb",
		CompileCmd:    []string{"ruby", "-c", "{file}"},
		TestCmd:       []string{"rspec", "{file}"},
		DummyTestCode: "require 'rspec'\nRSpec.describe 'Dummy' do\n  it 'passes' do\n    expect(true).to eq(true)\n  end\nend\n",
	},
	"cpp": {
		LanguageName:  "C++",
		Extension:     ".cpp",
		TestExtension: "_test.cpp",
		CompileCmd:    []string{"g++", "-fsyntax-only", "{file}"},
		TestCmd:       []string{"g++", "{file}"},
		DummyTestCode: "int main() { return 0; }\n",
	},
}
