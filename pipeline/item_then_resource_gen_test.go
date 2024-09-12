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
