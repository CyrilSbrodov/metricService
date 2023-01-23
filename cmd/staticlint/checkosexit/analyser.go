/*
	Package checkosexit анализатор проверки того, что в пакете main и в функции
	main нет вызова функции os.Exit
*/
package checkosexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name:     "osexitcheck",
	Doc:      "check for os.Exit func",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	oscheck := func(list []ast.Stmt) {
		for _, stmt := range list {
			if e, ok := stmt.(*ast.ExprStmt); ok {
				if c, ok := e.X.(*ast.CallExpr); ok {
					if s, ok := c.Fun.(*ast.SelectorExpr); ok {
						if i, ok := s.X.(*ast.Ident); ok {
							if i.Name == "os" {
								if s.Sel.Name == "Exit" {
									pass.Reportf(s.Pos(), "func os.Exit in main")
								}
							}
						}
					}
				}
			}
		}
	}

	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					if fd.Name.Name == "main" {
						oscheck(fd.Body.List)
						return false
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
