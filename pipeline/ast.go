package pipeline

import (
	"fmt"
	"go/ast"
	"go/token"
)

type AstGenerator struct {
}

type AstFile struct {
	file *ast.File
}

func CreateImportDecl(imports []string) ast.Decl {
	specs := make([]ast.Spec, len(imports))
	for i, imp := range imports {
		specs[i] = &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s"`, imp),
			},
		}
	}

	return &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: specs,
	}
}
