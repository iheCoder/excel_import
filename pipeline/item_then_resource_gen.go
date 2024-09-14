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
	structName           = "SectionImporter"
	dbField              = Field{VarName: "tx", Type: "*gorm.DB"}
	rcField              = Field{VarName: "rc", Type: "*general_framework.RawContent"}
	importSectionFuncDef = FuncDef{
		FuncName: "ImportSection",
		Params:   []Field{dbField, rcField},
		Results:  []Field{{Type: "error"}},
	}

	emptyStmt = &ast.EmptyStmt{}
)

type CaseResourceItem struct {
	CondVars []Var
	Info     *StructInfo
}

type ItemResourceAstGenerator struct {
	f *ast.File
	// struct name and its fields relation
	relations map[string]*StructFieldsRelation
	// case resource items
	caseResourceItems []*CaseResourceItem
	// var manager
	mgr *VarMgr
	// resource struct info
	resourceInfo *StructInfo
	// switch field
	switchField *Field
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

func (i *ItemResourceAstGenerator) CreateStructAssignStmt(receptor *StructInfo) *ast.AssignStmt {
	// get struct fields relation
	relation, ok := i.relations[receptor.Name]
	if !ok {
		return nil
	}

	// create struct assign statement
	assignStmt := CreateStructAssignStmt(relation)
	return assignStmt
}

func CreateGormDBCreateBlockStmt(db, model Var) []ast.Stmt {
	// create err := tx.Create(&model) statement
	errName := DefaultErrToken
	fc := &FuncCall{
		FuncName: "Create",
		Args:     []Var{model},
		ReturnVars: []Var{
			{
				Name: errName,
			},
		},
		Receiver: &StructInfo{
			VarName: db.Name,
		},
	}
	fcs := CreateFuncCallStmt(fc)

	// create if statement
	ifStmt := CreateIfErrIsNotNilStmt(errName)

	// Combine the assignment and if statement
	return []ast.Stmt{fcs, ifStmt}
}

func (i *ItemResourceAstGenerator) AddImportSectionFunc(receiver *StructInfo) {
	funcDef := importSectionFuncDef
	funcDef.Receiver = receiver

	// register var
	i.mgr.AddScopeAtRoot(funcDef.FuncName)
	for _, param := range funcDef.Params {
		i.mgr.AddVarInScope(param.VarName, funcDef.FuncName)
	}

	// create func declaration
	funcDecl := CreateFuncDecl(&funcDef)

	// create type assertion statement
	sourceVar := Var{Name: rcField.VarName + ".GetModel()"}
	mvName, _ := i.mgr.GenerateVarNameInScope(i.resourceInfo.Name, funcDef.FuncName)
	modelVar := Var{Name: mvName, Type: i.resourceInfo.Name}
	typeAssertStmt := i.createTypeAssertStmt(sourceVar, modelVar)

	// create switch statement
	i.setProviderVarName(modelVar)
	switchStmt := i.CreateSwitchCreateResourceItem(Var{Name: dbField.VarName}, modelVar, i.switchField, &funcDef)

	// add statements to the function body
	funcDecl.Body.List = append(funcDecl.Body.List, typeAssertStmt, emptyStmt, switchStmt)

	// add function declaration to the file
	i.f.Decls = append(i.f.Decls, funcDecl)
}

func (i *ItemResourceAstGenerator) setProviderVarName(modelVar Var) {
	for _, relation := range i.relations {
		for _, field := range relation.Fields {
			field.ProviderVarName = modelVar.Name
		}
	}
}

func (i *ItemResourceAstGenerator) CreateSwitchCreateResourceItem(dbVar, resVar Var, field *Field, fd *FuncDef) *ast.SwitchStmt {
	// create case clauses
	caseCluases := make([]*ast.CaseClause, 0, len(i.caseResourceItems))
	for _, item := range i.caseResourceItems {
		// get struct fields relation
		relation, ok := i.relations[item.Info.Name]
		if !ok {
			continue
		}

		// generate model var
		modelVar, _ := i.mgr.GenerateVarNameInScope(item.Info.Name, fd.FuncName)

		// create case clause
		caseClause := CreateCreateModelCaseClause(dbVar, Var{Name: modelVar}, item.CondVars, relation)
		caseCluases = append(caseCluases, caseClause)
	}

	// create switch statement
	switchStmt := CreateSwitchStmt(resVar.Name, field.Name, caseCluases)

	return switchStmt
}

func CreateCreateModelCaseClause(dbVar, modelVar Var, condVars []Var, relation *StructFieldsRelation) *ast.CaseClause {
	// create struct assign statement
	relation.Info.VarName = modelVar.Name
	assignStmt := CreateStructAssignStmt(relation)

	// create gorm db create block statement
	createBlockStmt := CreateGormDBCreateBlockStmt(dbVar, modelVar)

	// create case clause
	stmts := []ast.Stmt{assignStmt}
	stmts = append(stmts, createBlockStmt...)
	return CreateCaseClause(condVars, stmts)
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
