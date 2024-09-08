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

type StructInfo struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name    string
	Type    string
	Comment string
}

// CreateImportDecl creates an import declaration with the given import paths.
func CreateImportDecl(imports []string) ast.Decl {
	specs := make([]ast.Spec, len(imports))
	for i, imp := range imports {
		// If the import path is empty, create an empty import spec.
		if len(imp) == 0 {
			specs[i] = &ast.ImportSpec{}
			continue
		}

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

// CreateStructDecl creates a struct declaration with the given struct info.
func CreateStructDecl(info *StructInfo) ast.Decl {
	// Create a field list with the given fields.
	fields := make([]*ast.Field, len(info.Fields))
	for i, field := range info.Fields {
		fields[i] = &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(field.Name)},
			Type:  ast.NewIdent(field.Type),
			Doc:   &ast.CommentGroup{List: []*ast.Comment{{Text: fmt.Sprintf("// %s", field.Comment)}}},
		}
	}

	// Create a type spec with the given struct name and fields.
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(info.Name),
				Type: &ast.StructType{
					Fields: &ast.FieldList{List: fields},
				},
			},
		},
	}
}
