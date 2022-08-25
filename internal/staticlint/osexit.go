// Package staticlint for project specific static linters
// check usage os.Exit() in main() function
package staticlint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check os.Exit in main functions",
	Run:  run,
}

// run linter check os.Exit() in main() function
func run(pass *analysis.Pass) (interface{}, error) {

	// обойти дерево AST
	// найти функцию main()
	// в функции main() найти os.Exit(int i)

	for _, file := range pass.Files {

		ast.Inspect(file, func(node ast.Node) bool {

			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.String() == "main" {

					ast.Inspect(x.Body, func(nodeMain ast.Node) bool {

						switch nx := nodeMain.(type) {
						case *ast.SelectorExpr:
							if nx.Sel.Name == "Exit" {

								ast.Inspect(nx, func(nodeExit ast.Node) bool {
									switch nxx := nodeExit.(type) {
									case *ast.Ident:
										if nxx.Name == "os" {
											pass.Reportf(nxx.Pos(), "remove call os.Exit() from main()")
										}
									}
									return true
								})
							}
						}
						return true
					})
				}
			}
			return true
		})
	}
	return nil, nil
}
