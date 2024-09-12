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
