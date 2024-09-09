package pipeline

import (
	"fmt"
	"go/ast"
	"go/token"
)

const (
	DefaultOKVarName = "ok"
)

type AstGenerator struct {
}

type AstFile struct {
	file *ast.File
}

type StructInfo struct {
	Name    string
	Fields  []Field
	VarName string
}

type Field struct {
	Name    string
	Type    string
	Comment string
	VarName string
}

type FuncDef struct {
	SI       *StructInfo
	FuncName string
	Params   []Field
	Results  []Field
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
	fields := CreateFields(info.Fields)

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

// CreateFuncDecl creates a function declaration with the given function definition.
func CreateFuncDecl(def *FuncDef) *ast.FuncDecl {
	// If the struct info is not nil, create a receiver with the struct name.
	var recv *ast.FieldList
	if def.SI != nil {
		recv = &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(def.SI.VarName)},
					Type:  ast.NewIdent(def.SI.Name),
				},
			},
		}
	}

	// Create a field list with the given parameters.
	params := make([]*ast.Field, len(def.Params))
	for i, param := range def.Params {
		params[i] = &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(param.VarName)},
			Type:  ast.NewIdent(param.Type),
		}
	}

	// Create a field list with the given results.
	results := make([]*ast.Field, len(def.Results))
	for i, result := range def.Results {
		results[i] = &ast.Field{
			Type: ast.NewIdent(result.Type),
		}
	}

	// Create a function type with the given parameters and results.
	ftype := &ast.FuncType{
		Params:  &ast.FieldList{List: params},
		Results: &ast.FieldList{List: results},
	}

	// Create a function declaration with the given function name, type, and body.
	return &ast.FuncDecl{
		Recv: recv,
		Name: ast.NewIdent(def.FuncName),
		Type: ftype,
	}
}

func CreateFields(fields []Field) []*ast.Field {
	astFields := make([]*ast.Field, len(fields))
	for i, field := range fields {
		astFields[i] = &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(field.Name)},
			Type:  ast.NewIdent(field.Type),
			Doc:   &ast.CommentGroup{List: []*ast.Comment{{Text: fmt.Sprintf("// %s", field.Comment)}}},
		}
	}
	return astFields
}

// CreateTypeAssertStmt creates a type assertion statement with the given sourceName, targetName, and typeName.
func CreateTypeAssertStmt(sourceName, targetName, typeName string, stmt []ast.Stmt) *ast.IfStmt {
	// Create a type assertion statement with the given source name, target name, and type name.
	return &ast.IfStmt{
		Init: &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(targetName), ast.NewIdent(DefaultOKVarName)},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.TypeAssertExpr{
					X:    ast.NewIdent(sourceName),
					Type: ast.NewIdent(typeName),
				},
			},
		},
		Cond: &ast.UnaryExpr{
			Op: token.NOT,
			X:  ast.NewIdent(DefaultOKVarName),
		},
		Body: &ast.BlockStmt{
			List: stmt,
		},
	}
}

// CreateReturnErrStmt creates a return statement with the given error name.
func CreateReturnErrStmt(errName string) *ast.ReturnStmt {
	// Create a return statement with the given error name.
	return &ast.ReturnStmt{
		Results: []ast.Expr{ast.NewIdent(errName)},
	}
}
