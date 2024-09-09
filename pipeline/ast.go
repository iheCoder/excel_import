package pipeline

import (
	"fmt"
	"go/ast"
	"go/token"
)

const (
	DefaultOKVarName = "ok"
	NilToken         = "nil"
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
	Struct  *StructInfo
}

type Var struct {
	Name string
	Type string
}

type StructFieldsRelation struct {
	Info   StructInfo
	Fields []FieldRelation
}

type FieldRelation struct {
	Source, Target Field
}

type FuncDef struct {
	Receiver *StructInfo
	FuncName string
	Params   []Field
	Results  []Field
}

type FuncCall struct {
	FuncName   string
	Args       []Var
	ReturnVars []Var
	Receiver   *StructInfo
}

// CreateImportDecl creates an import declaration with the given import paths.
func CreateImportDecl(imports []string) ast.Decl {
	specs := make([]ast.Spec, len(imports))
	for i, imp := range imports {
		// If the import path is empty, create an empty import spec.
		if len(imp) == 0 {
			specs[i] = &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind: token.STRING,
				},
			}
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
	if def.Receiver != nil {
		recv = &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(def.Receiver.VarName)},
					Type:  ast.NewIdent(def.Receiver.Name),
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
			Comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{
						Text: fmt.Sprintf("\t//\t%s", field.Comment),
					},
				},
			},
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

// CreateSwitchStmt creates a switch statement with the given selector and cases.
func CreateSwitchStmt(x, sel string, cases []*ast.CaseClause) *ast.SwitchStmt {
	list := make([]ast.Stmt, len(cases))
	for i, c := range cases {
		list[i] = c
	}

	return &ast.SwitchStmt{
		Tag: &ast.SelectorExpr{
			X:   ast.NewIdent(x),
			Sel: ast.NewIdent(sel),
		},
		Body: &ast.BlockStmt{
			List: list,
		},
	}
}

// CreateStructValueSpec creates a value specification with the given struct fields relation.
func CreateStructValueSpec(relation StructFieldsRelation) *ast.ValueSpec {
	// Create a composite literal with the given struct name.
	cl := &ast.CompositeLit{
		Type: ast.NewIdent(relation.Info.Name),
	}

	// Create a key-value expression with the given source and target fields.
	elts := make([]ast.Expr, len(relation.Fields))
	for i, f := range relation.Fields {
		elts[i] = &ast.KeyValueExpr{
			Key:   ast.NewIdent(f.Target.Name),
			Value: &ast.SelectorExpr{X: ast.NewIdent(f.Source.Struct.VarName), Sel: ast.NewIdent(f.Source.Name)},
		}
	}

	// Set the elements to the composite literal.
	cl.Elts = elts

	// Create a value specification with the composite literal.
	return &ast.ValueSpec{
		Names:  []*ast.Ident{ast.NewIdent(relation.Info.VarName)},
		Type:   ast.NewIdent(relation.Info.Name),
		Values: []ast.Expr{cl},
	}
}

// CreateFuncCallStmt creates a function call statement with the given function call.
func CreateFuncCallStmt(call *FuncCall) *ast.ExprStmt {
	// Create a function call expression with the given function name and arguments.
	fc := &ast.CallExpr{
		Fun: ast.NewIdent(call.FuncName),
	}

	// If the receiver is not nil, set the receiver to the function call expression.
	if call.Receiver != nil {
		fc.Fun = &ast.SelectorExpr{
			X:   ast.NewIdent(call.Receiver.VarName),
			Sel: ast.NewIdent(call.FuncName),
		}
	}

	// Set the arguments to the function call expression.
	for _, arg := range call.Args {
		fc.Args = append(fc.Args, ast.NewIdent(arg.Name))
	}

	// Set the return variables to the function call expression.
	for _, rv := range call.ReturnVars {
		fc.Args = append(fc.Args, ast.NewIdent(rv.Name))
	}

	// Create an expression statement with the function call expression.
	return &ast.ExprStmt{
		X: fc,
	}
}

// CreateIfErrIsNotNilStmt creates an if statement with the given error name and return err statements.
func CreateIfErrIsNotNilStmt(errName string) *ast.IfStmt {
	// Create an if statement with the given error name and return err statements.
	return &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  ast.NewIdent(errName),
			Op: token.NEQ,
			Y:  ast.NewIdent(NilToken),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				CreateReturnErrStmt(errName),
			},
		},
	}
}

// CreateNewStructReturnStmt creates a return statement with the given struct info.
func CreateNewStructReturnStmt(info *StructInfo) *ast.ReturnStmt {
	// Create a composite literal with the given struct name.
	cl := &ast.CompositeLit{
		Type: ast.NewIdent(info.Name),
	}

	// Create a key-value expression with the given struct fields.
	elts := make([]ast.Expr, len(info.Fields))
	for i, f := range info.Fields {
		elts[i] = &ast.KeyValueExpr{
			Key:   ast.NewIdent(f.Name),
			Value: ast.NewIdent(f.VarName),
		}
	}

	// Set the elements to the composite literal.
	cl.Elts = elts

	// Create a return statement with the composite literal.
	return &ast.ReturnStmt{
		Results: []ast.Expr{cl},
	}
}
