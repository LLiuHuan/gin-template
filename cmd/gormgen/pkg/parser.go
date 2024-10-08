// Package pkg
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:46
package pkg

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/LLiuHuan/gin-template/pkg/gormDB"
)

// Parser 解析器用于解析目录并公开有关该目录文件中定义的结构的信息。
type Parser struct {
	dir         string
	pkg         *build.Package
	parsedFiles []*ast.File
}

// NewParser 创建一个新的解析器实例。
func NewParser(dir string) *Parser {
	return &Parser{
		dir: dir,
	}
}

// getPackage 解析目录获取go文件和包
func (p *Parser) getPackage() {
	pkg, err := build.Default.ImportDir(p.dir, build.ImportComment)
	if err != nil {
		log.Fatalf("cannot process directory %s: %s", p.dir, err)
	}
	p.pkg = pkg

}

// parseGoFiles 解析go文件
func (p *Parser) parseGoFiles() {
	var parsedFiles []*ast.File
	fs := token.NewFileSet()
	for _, file := range p.pkg.GoFiles {
		file = p.dir + "/" + file
		parsedFile, err := parser.ParseFile(fs, file, nil, 0)
		if err != nil {
			log.Fatalf("parsing package: %s: %s\n", file, err)
		}
		parsedFiles = append(parsedFiles, parsedFile)
	}
	p.parsedFiles = parsedFiles
}

// parseTypes 解析结构体类型
func (p *Parser) parseTypes(file *ast.File) (ret []structConfig) {
	ast.Inspect(file, func(n ast.Node) bool {
		decl, ok := n.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			return true
		}

		for _, spec := range decl.Specs {
			var (
				data structConfig
			)
			typeSpec, _ok := spec.(*ast.TypeSpec)
			if !_ok {
				continue
			}
			// 我们只关心结构声明（目前）
			var structType *ast.StructType
			if structType, ok = typeSpec.Type.(*ast.StructType); !ok {
				continue
			}

			data.StructName = typeSpec.Name.Name
			for _, v := range structType.Fields.List {
				var (
					optionField fieldConfig
				)

				if t, _ok := v.Type.(*ast.Ident); _ok {
					optionField.FieldType = t.String()
				} else {
					if v.Tag != nil {
						if strings.Contains(v.Tag.Value, "gormDB") && strings.Contains(v.Tag.Value, "time") {
							optionField.FieldType = "time.Time"
						}
					}
				}

				if len(v.Names) > 0 {
					optionField.FieldName = v.Names[0].String()
					optionField.ColumnName = gormDB.ToDBName(optionField.FieldName)
					optionField.HumpName = SQLColumnToHumpStyle(optionField.ColumnName)
				}

				data.OptionFields = append(data.OptionFields, optionField)
			}

			ret = append(ret, data)
		}
		return true
	})
	return
}

// Parse 应该在解析器的任何类型查询之前调用。它需要解析目录并提取该目录中定义的所有结构。
func (p *Parser) Parse() (ret []structConfig) {
	var (
		data []structConfig
	)
	p.getPackage()
	p.parseGoFiles()
	for _, f := range p.parsedFiles {
		data = append(data, p.parseTypes(f)...)
	}
	return data
}
