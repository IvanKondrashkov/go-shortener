package exitcheck

import (
	"strings"

	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// ExitCheckAnalyzer определяет анализатор для проверки прямого вызова os.Exit
var ExitCheckAnalyzer = &analysis.Analyzer{
	Name:     "exitcheck",
	Doc:      "check call os.Exit in main package",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Пропускаем тестовые файлы
	for _, filename := range pass.Files {
		if strings.HasSuffix(pass.Fset.File(filename.Pos()).Name(), "_test.go") {
			return nil, nil
		}
	}

	// Пропускаем если это не пакет main
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	ins.Preorder(nodeFilter, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)

		pos := pass.Fset.Position(callExpr.Pos())
		if isTempGoBuildFile(pos.Filename) {
			return
		}

		// Проверяем, является ли вызов функцией Exit из пакета os
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selExpr.X.(*ast.Ident); ok {
				if ident.Name == "os" && selExpr.Sel.Name == "Exit" {
					// Проверяем, находимся ли мы в функции main
					if isExitMain(pass, callExpr.Pos()) {
						pass.Reportf(callExpr.Pos(), "check call os.Exit in main package")
					}
				}
			}
		}
	})

	return nil, nil
}

func isTempGoBuildFile(filename string) bool {
	return strings.Contains(filename, "\\go-build\\") ||
		strings.Contains(filename, "/go-build/") ||
		strings.Contains(filename, "\\Temp\\go-build") ||
		strings.Contains(filename, "/Temp/go-build")
}

func isExitMain(pass *analysis.Pass, pos token.Pos) bool {
	for _, file := range pass.Files {
		if pass.Fset.Position(pos).Filename != pass.Fset.Position(file.Pos()).Filename {
			continue
		}

		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == "main" {
					// Проверяем, находится ли вызов внутри этой функции main
					funcStart := funcDecl.Pos()
					funcEnd := funcDecl.End()
					if pos >= funcStart && pos <= funcEnd {
						return true
					}
				}
			}
		}
	}
	return false
}
