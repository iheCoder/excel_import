package pipeline

import (
	"go/ast"
)

var (
	importsStmt = []string{
		"errors",
		"gorm.io/gorm",
		"",
		"excel_import/general_framework",
	}
	structName = "SectionImporter"
)

type ItemResourceAstGenerator struct {
	f *ast.File
	// struct name and its fields relation
	relations map[string]StructFieldsRelation
}

func (i *ItemResourceAstGenerator) AddImportDecl() {
	i.f.Decls = append(i.f.Decls, CreateImportDecl(importsStmt))
}

func (i *ItemResourceAstGenerator) AddStructDecl() {
	i.f.Decls = append(i.f.Decls, CreateStructDecl(&StructInfo{
		Name: structName,
	}))
}

func (i *ItemResourceAstGenerator) AddNewFuncStructDecl() {
	// create return statement
	ret := CreateNewStructReturnStmt(&StructInfo{
		Name: structName,
	})

	// create new func declaration
	newFuncDelc := CreateFuncDecl(&FuncDef{
		FuncName: "New" + structName,
		Results:  []Field{{Type: "*" + structName}},
	})

	// add return statement to the new func declaration
	newFuncDelc.Body.List = append(newFuncDelc.Body.List, ret)

	i.f.Decls = append(i.f.Decls, newFuncDelc)
}

func (i *ItemResourceAstGenerator) createTypeAssertStmt(source, target Var) *ast.IfStmt {
	// create return error statement
	retErr := CreateReturnErrStmt("errors.New(\"type assertion failed\")")

	// create type assertion statement
	typeAssertStmt := CreateTypeAssertStmt(source, target, []ast.Stmt{retErr})
	return typeAssertStmt
}

func (i *ItemResourceAstGenerator) createStructAssignStmt(receptor StructInfo) *ast.AssignStmt {
	// get struct fields relation
	relation, ok := i.relations[receptor.Name]
	if !ok {
		return nil
	}

	// create struct assign statement
	assignStmt := CreateStructAssignStmt(relation)
	return assignStmt
}

func TransferStructFieldsRelation(info *StructInfo, graph *ModelGraph) StructFieldsRelation {
	// transfer struct field relation
	fieldsRelation := make([]FieldRelation, 0, len(info.Fields))
	for _, field := range info.Fields {
		// get edge from the field node
		to, _, ok := graph.GetOneEdge(FieldNode{
			StructName: info.Name,
			FieldName:  field.Name,
		})
		if !ok {
			continue
		}

		// create field relation
		fieldsRelation = append(fieldsRelation, FieldRelation{
			ReceptorFieldName: field.Name,
			// TODO: replace into var name
			ProviderVarName:   to.StructName,
			ProviderFieldName: to.FieldName,
		})
	}

	return StructFieldsRelation{
		Info:   *info,
		Fields: fieldsRelation,
	}
}
