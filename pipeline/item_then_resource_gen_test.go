package pipeline

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"testing"
)

func TestCreateGormDBCreateIfStmt(t *testing.T) {
	db := Var{
		Name: "tx",
	}
	model := Var{
		Name: "&model",
	}
	stmt := CreateGormDBCreateBlockStmt(db, model)
	want := "err := tx.Create(&model)\nif err != nil {\n\treturn err\n}"
	got := stmtsToString(stmt)
	if got != want {
		t.Errorf("CreateGormDBCreateIfStmt(%v, %v) =\n %v\n, want\n %v", db, model, got, want)
	}
}

func ifStmtToString(stmt *ast.IfStmt) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), stmt)
	if err != nil {
		log.Fatalf("Error printing AST: %v", err)
	}

	return buf.String()
}

func TestCreateCreateModelCaseClause(t *testing.T) {
	dbVar := Var{
		Name: "tx",
	}
	condVars := []Var{
		{
			Name: "1",
		},
	}
	relation := &StructFieldsRelation{
		Info: StructInfo{
			Name: "Genshin",
		},
		Fields: []FieldRelation{
			{
				ReceptorFieldName: "Hero",
				ProviderVarName:   "excelModel",
				ProviderFieldName: "Hero",
			},
		},
	}
	modelVar := Var{
		Name: "model",
	}
	stmt := CreateCreateModelCaseClause(dbVar, modelVar, condVars, relation)
	want := "case 1:\n\tmodel := Genshin{Hero: excelModel.Hero}\n\terr := tx.Create(model)\n\tif err != nil {\n\t\treturn err\n\t}"
	got := stmtToString(stmt)
	if got != want {
		t.Errorf("CreateCreateModelCaseClause() =\n %v\n, want\n %v", got, want)
	}
}

// assume that the relation is already created
// case 1: Var{type: Resource1, name: resource1}
// default: Var{type: Resource2, name: resource2}
func TestItemResourceAstGenerator_AddSwitchCreateResourceItem(t *testing.T) {
	// create relations
	relations := make(map[string]*StructFieldsRelation)
	relations["Resource1"] = &StructFieldsRelation{
		Info: StructInfo{
			Name: "Resource1",
		},
		Fields: []FieldRelation{
			{
				ReceptorFieldName: "Name",
				ProviderVarName:   "excelModel",
				ProviderFieldName: "Name",
			},
		},
	}
	relations["Resource2"] = &StructFieldsRelation{
		Info: StructInfo{
			Name: "Resource2",
		},
		Fields: []FieldRelation{
			{
				ReceptorFieldName: "Name",
				ProviderVarName:   "excelModel",
				ProviderFieldName: "Name",
			},
		},
	}

	// create case resource items
	items := []*CaseResourceItem{
		{
			CondVars: []Var{
				{
					Name: "1",
				},
			},
			Info: &StructInfo{
				Name: "Resource1",
			},
		},
		{
			Info: &StructInfo{
				Name: "Resource2",
			},
		},
	}

	funcName := "CreateResource"
	// create var mgr
	mgr := NewVarMgr()
	mgr.AddScopeAtRoot(funcName)

	// create item resource ast generator
	generator := &ItemResourceAstGenerator{
		relations:         relations,
		caseResourceItems: items,
		mgr:               mgr,
	}

	// create ast func decl
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent(funcName),
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{},
	}

	// add switch stmt
	dbVar := Var{
		Name: "tx",
	}
	resVar := Var{
		Name: "resource",
	}
	field := &Field{
		Name: "Type",
	}
	fd := &FuncDef{
		FuncName: funcName,
	}
	generator.AddSwitchCreateResourceItem(dbVar, resVar, field, fd, funcDecl)

	// check the switch stmt
	want := "func CreateResource() {\n\tswitch resource.Type {\n\tcase 1:\n\t\tr := Resource1{Name: excelModel.Name}\n\t\terr := tx.Create(r)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\tdefault:\n\t\tresource2 := Resource2{Name: excelModel.Name}\n\t\terr := tx.Create(resource2)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n}"
	got := funcDeclToString(funcDecl)
	if got != want {
		t.Errorf("AddSwitchCreateResourceItem() =\n %v\n, want\n %v", got, want)
	}

	t.Log("PASS")
}

func funcDeclToString(fd *ast.FuncDecl) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), fd)
	if err != nil {
		log.Fatalf("Error printing AST: %v", err)
	}

	return buf.String()
}
