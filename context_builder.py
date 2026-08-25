import re
import os

def build_context(entry_file_path: str, base_dir: str = "legacy_app") -> str:
    """
    Parses a PHP file for include/require statements, reads those files,
    and bundles them into a single context string for the LLM.
    """
    print(f"[Context Builder]: Analyzing AST and dependencies for {entry_file_path}...")
    
    with open(entry_file_path, "r") as f:
        main_code = f.read()
        
    # Regex to find include 'file.php'; or require "file.php";
    pattern = r"(?:include|require|include_once|require_once)\s*['\"]([^'\"]+)['\"]\s*;"
    dependencies = re.findall(pattern, main_code)
    
    context = "=== DEPENDENCY GRAPH CONTEXT ===\n"
    
    if dependencies:
        print(f"[Context Builder]: Found dependencies: {', '.join(dependencies)}")
        for dep in dependencies:
            dep_path = os.path.join(base_dir, dep)
            if os.path.exists(dep_path):
                with open(dep_path, "r") as dep_file:
                    context += f"\n--- File: {dep} ---\n{dep_file.read()}\n"
            else:
                context += f"\n--- File: {dep} (NOT FOUND) ---\n"
    else:
        print("[Context Builder]: No external dependencies found.")
        
    context += "\n=== MAIN FILE TO MIGRATE ===\n"
    context += f"--- File: {os.path.basename(entry_file_path)} ---\n{main_code}\n"
    
    print("[Context Builder]: Context bundling complete.")
    return context

if __name__ == "__main__":
    ctx = build_context("legacy_app/ProcessOrder.php")
    print("\n[Preview of Context String sent to AI]:\n")
    print(ctx)
