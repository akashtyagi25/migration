<?php
require 'vendor/autoload.php';

use PhpParser\Error;
use PhpParser\NodeDumper;
use PhpParser\ParserFactory;
use PhpParser\Node;
use PhpParser\NodeVisitorAbstract;
use PhpParser\NodeTraverser;

if ($argc < 2) {
    echo json_encode(["error" => "Please provide a file path."]);
    exit(1);
}

$code = file_get_contents($argv[1]);

$parser = (new ParserFactory())->createForNewestSupportedVersion();

class DependencyVisitor extends NodeVisitorAbstract {
    public $dependencies = [];
    public $exported_symbols = [];

    public function leaveNode(Node $node) {
        if ($node instanceof Node\Expr\Include_) {
            // Include types: 1 = include, 2 = include_once, 3 = require, 4 = require_once
            if ($node->expr instanceof Node\Scalar\String_) {
                $this->dependencies[] = $node->expr->value;
            } elseif ($node->expr instanceof Node\Expr\BinaryOp\Concat) {
                // E.g., include $base_path . 'config.php';
                $right = $node->expr->right;
                if ($right instanceof Node\Scalar\String_) {
                    $this->dependencies[] = "[Dynamic Concat] -> " . $right->value;
                } else {
                    $this->dependencies[] = "[Dynamic Variable Include]";
                }
            } else {
                $this->dependencies[] = "[Complex Dynamic Include]";
            }
        }
        
        // Smart Chunking: Extract Class and Method signatures
        if ($node instanceof Node\Stmt\Class_) {
            $className = $node->name ? $node->name->toString() : 'AnonymousClass';
            $methods = [];
            foreach ($node->getMethods() as $method) {
                $methods[] = $method->name->toString();
            }
            $this->exported_symbols[] = [
                "type" => "class",
                "name" => $className,
                "methods" => $methods
            ];
        }

        // Smart Chunking: Extract standalone Functions
        if ($node instanceof Node\Stmt\Function_) {
            $this->exported_symbols[] = [
                "type" => "function",
                "name" => $node->name->toString(),
                "methods" => []
            ];
        }
    }
}

try {
    $ast = $parser->parse($code);
    $traverser = new NodeTraverser();
    $visitor = new DependencyVisitor();
    $traverser->addVisitor($visitor);
    $traverser->traverse($ast);
    
    echo json_encode([
        "file" => basename($argv[1]),
        "dependencies" => $visitor->dependencies,
        "exported_symbols" => $visitor->exported_symbols,
        "ast_nodes_parsed" => count($ast)
    ], JSON_PRETTY_PRINT);
} catch (Error $error) {
    echo json_encode(["error" => "Parse error: {$error->getMessage()}"]);
}
?>
