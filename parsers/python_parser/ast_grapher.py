import ast
import sys
import json
import os

if len(sys.argv) < 2:
    print(json.dumps({"error": "Please provide a file path."}))
    sys.exit(1)

file_path = sys.argv[1]

try:
    with open(file_path, "r", encoding="utf-8") as f:
        code = f.read()
except Exception as e:
    print(json.dumps({"error": f"Failed to read file: {e}"}))
    sys.exit(1)

try:
    tree = ast.parse(code)
except Exception as e:
    print(json.dumps({"error": f"Parse error: {e}"}))
    sys.exit(1)

dependencies = []
exported_symbols = []

class DependencyVisitor(ast.NodeVisitor):
    def visit_Import(self, node):
        for alias in node.names:
            dependencies.append(alias.name + ".py")
        self.generic_visit(node)

    def visit_ImportFrom(self, node):
        if node.module:
            dependencies.append(node.module + ".py")
        self.generic_visit(node)

    def visit_ClassDef(self, node):
        methods = [n.name for n in node.body if isinstance(n, ast.FunctionDef)]
        exported_symbols.append({
            "type": "class",
            "name": node.name,
            "methods": methods
        })
        self.generic_visit(node)

    def visit_FunctionDef(self, node):
        # Only add standalone functions (not methods inside classes)
        # We can check parent, but simple visit is fine for MVP
        if not hasattr(self, 'current_class'):
            exported_symbols.append({
                "type": "function",
                "name": node.name,
                "methods": []
            })
        self.generic_visit(node)

visitor = DependencyVisitor()
visitor.visit(tree)

print(json.dumps({
    "file": os.path.basename(file_path),
    "dependencies": dependencies,
    "exported_symbols": exported_symbols,
    "ast_nodes_parsed": len(list(ast.walk(tree)))
}, indent=4))
