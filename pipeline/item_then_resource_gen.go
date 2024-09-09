package pipeline

import "go/ast"

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
