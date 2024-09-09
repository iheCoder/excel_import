package pipeline

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"testing"
)

func TestCreateImportDecl(t *testing.T) {
	type testData struct {
		imports []string
		want    string
	}

	tests := []testData{
		{
			imports: []string{"fmt"},
			want:    `import "fmt"`,
		},
		{
			imports: []string{"fmt", "strings"},
			want: `import (
	"fmt"
	"strings"
)`,
		},
	}

	for _, test := range tests {
		decl := CreateImportDecl(test.imports)
		got := declToString(decl)
		if got != test.want {
			t.Errorf("CreateImportDecl(%v) =\n %v\n, want\n %v", test.imports, got, test.want)
		}
	}
}

func declToString(decl ast.Decl) string {
	// 创建一个 bytes.Buffer 来存储生成的代码
	var buf bytes.Buffer

	// 使用 go/printer 包将 ast.Decl 写入 buffer
	err := printer.Fprint(&buf, token.NewFileSet(), decl)
	if err != nil {
		log.Fatalf("Error printing AST: %v", err)
	}

	// 返回 buffer 中的内容作为字符串
	return buf.String()
}

func TestCreateStructDecl(t *testing.T) {
	type testData struct {
		name   string
		fields []Field
		want   string
	}

	tests := []testData{
		{
			name: "test",
			fields: []Field{
				{
					Name:    "ID",
					Type:    "int",
					Comment: "id",
				},
				{
					Name:    "Name",
					Type:    "string",
					Comment: "name",
				},
			},
			want: "type test struct {\n\tID\tint\t//\tid\n\tName\tstring\t//\tname\n}",
		},
		{
			name: "test",
			fields: []Field{
				{
					Name:    "A",
					Type:    "int",
					Comment: "你好",
				},
			},
			want: "type test struct {\n\tA int\t//\t你好\n}",
		},
	}

	for _, test := range tests {
		decl := CreateStructDecl(&StructInfo{
			Name:   test.name,
			Fields: test.fields,
		})
		got := declToString(decl)
		if got != test.want {
			t.Errorf("CreateStructDecl(%s, %v) =\n %v\n, want\n %v", test.name, test.fields, got, test.want)
		}
	}
}
