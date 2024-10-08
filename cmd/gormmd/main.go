// Package main
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:51
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/LLiuHuan/gin-template/pkg/gormDB"

	"gorm.io/gorm"
)

// TODO: 支持多个数据库

// tableInfo 表基础信息
type tableInfo struct {
	Name    string         `gorm:"table_name"`    // 表名
	Comment sql.NullString `gorm:"table_comment"` // 表描述
}

// tableColumn 表字段信息
type tableColumn struct {
	OrdinalPosition uint16         `gorm:"ORDINAL_POSITION"` // 列在表中的位置（从1开始）, 表示该列是表中第几列
	ColumnName      string         `gorm:"COLUMN_NAME"`      // 列的名称
	ColumnType      string         `gorm:"COLUMN_TYPE"`      // 列的完整数据类型定义，包括长度、是否为 UNSIGNED、是否为 ZEROFILL 等信息, 例如 varchar(255) int(11) unsigned
	DataType        string         `gorm:"DATA_TYPE"`        // 列的数据类型（不带长度）。常见的数据类型包括 int、varchar、date、text 等
	ColumnKey       sql.NullString `gorm:"COLUMN_KEY"`       // 列在表中的索引类型。可能的值包括： PRI: 主键。 UNI: 唯一索引。 MUL: 可以重复的非唯一索引（普通索引）。
	IsNullable      string         `gorm:"IS_NULLABLE"`      // 该列是否允许存储 NULL 值。值为 YES 或 NO
	Extra           sql.NullString `gorm:"EXTRA"`            // 列的额外信息，如自动递增（auto_increment）或其他特殊属性。
	ColumnComment   sql.NullString `gorm:"COLUMN_COMMENT"`   // 列的注释，可以由用户通过 COMMENT 子句定义。
	ColumnDefault   sql.NullString `gorm:"COLUMN_DEFAULT"`   // 列的默认值。可以是 NULL 或其他常量。如果列没有设置默认值，这里将显示 NULL。
}

type dataTypeMapping func(detailType string) (finalType string)

var (
	dbDriver  string
	dbHost    string
	dbPort    string
	dbUser    string
	dbPass    string
	dbName    string
	genTables string
)

func init() {
	flag.StringVar(&dbDriver, "driver", "mysql", "请输入数据库类型，例如：mysql\n")
	flag.StringVar(&dbHost, "host", "", "请输入 gorm IP/Path，例如：127.0.0.1\n")
	flag.StringVar(&dbPort, "port", "", "请输入 gorm 端口，例如：3306\n")
	flag.StringVar(&dbUser, "user", "", "请输入 gorm 数据库用户名\n")
	flag.StringVar(&dbPass, "pass", "", "请输入 gorm 数据库密码，例如：123456\n")
	flag.StringVar(&dbName, "db", "", "请输入 gorm 数据库名称\n")
	flag.StringVar(&genTables, "tables", "*", "请输入要生成的表名，默认为“*”，多个可用“,”分割\n")

	if !flag.Parsed() {
		flag.Parse()
	}
	dbName = strings.ToLower(dbName)
	genTables = strings.ToLower(genTables)
}

func main() {
	// 初始化 DB
	db, err := gormDB.NewDB(dbDriver, dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatal("new gormDB err", err)
	}

	defer func() {
		if err := db.DBClose(); err != nil {
			log.Println("gormDB close err", err)
		}
	}()

	tables, err := queryTables(db.GetDB(), dbName, genTables)
	if err != nil {
		log.Println("query tables of gormDB err", err)
		return
	}

	for _, table := range tables {

		filepath := "./internal/repository/gormDB/" + table.Name
		_ = os.Mkdir(filepath, 0766)
		fmt.Println("create dir : ", filepath)

		mdName := fmt.Sprintf("%s/gen_table.md", filepath)
		mdFile, err := os.OpenFile(mdName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
		if err != nil {
			fmt.Printf("markdown file error %v\n", err.Error())
			return
		}
		fmt.Println("  └── file : ", table.Name+"/gen_table.md")

		modelName := fmt.Sprintf("%s/gen_model.go", filepath)
		modelFile, err := os.OpenFile(modelName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
		if err != nil {
			fmt.Printf("create and open model file error %v\n", err.Error())
			return
		}
		fmt.Println("  └── file : ", table.Name+"/gen_model.go")

		modelContent := fmt.Sprintf("package %s\n", table.Name)
		modelContent += fmt.Sprintf(`import "time"`)
		modelContent += fmt.Sprintf("\n\n// %s %s \n", capitalize(table.Name), table.Comment.String)
		modelContent += fmt.Sprintf("//go:generate gormgen -structs %s -input . \n", capitalize(table.Name))
		modelContent += fmt.Sprintf("type %s struct {\n", capitalize(table.Name))

		tableContent := fmt.Sprintf("#### %s.%s \n", dbName, table.Name)
		if table.Comment.String != "" {
			tableContent += table.Comment.String + "\n"
		}
		tableContent += "\n" +
			"| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |\n" +
			"| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |\n"

		columnInfo, columnInfoErr := queryTableColumn(db.GetDB(), dbName, table.Name)
		if columnInfoErr != nil {
			continue
		}
		for _, info := range columnInfo {
			tableContent += fmt.Sprintf(
				"| %d | %s | %s | %s | %s | %s | %s | %s |\n",
				info.OrdinalPosition,
				info.ColumnName,
				strings.ReplaceAll(strings.ReplaceAll(info.ColumnComment.String, "|", "\\|"), "\n", ""),
				info.ColumnType,
				info.ColumnKey.String,
				info.IsNullable,
				info.Extra.String,
				info.ColumnDefault.String,
			)

			if textType(info.DataType) == "time.Time" {
				modelContent += fmt.Sprintf("%s %s `%s` // %s\n", capitalize(info.ColumnName), textType(info.DataType), "gormDB:\"time\"", info.ColumnComment.String)
			} else {
				modelContent += fmt.Sprintf("%s %s // %s\n", capitalize(info.ColumnName), textType(info.DataType), info.ColumnComment.String)
			}
		}

		mdFile.WriteString(tableContent)
		mdFile.Close()

		modelContent += "}\n"
		modelContent += fmt.Sprintf("\nfunc (*%s) TableName() string {\n\treturn \"%s\"\n}\n", capitalize(table.Name), table.Name)
		modelFile.WriteString(modelContent)
		modelFile.Close()

	}

}

func queryTables(db *gorm.DB, dbName string, tableName string) ([]tableInfo, error) {
	var tableCollect []tableInfo
	var tableArray []string
	var commentArray []sql.NullString
	sqlTables := fmt.Sprintf("SELECT `table_name`,`table_comment` FROM `information_schema`.`tables` WHERE `table_schema`= '%s'", dbName)
	rows, err := db.Raw(sqlTables).Rows()
	if err != nil {
		return tableCollect, err
	}
	defer rows.Close()

	for rows.Next() {
		var info tableInfo
		err = rows.Scan(&info.Name, &info.Comment)
		if err != nil {
			fmt.Printf("execute query tables action error,had ignored, detail is [%v]\n", err.Error())
			continue
		}

		tableCollect = append(tableCollect, info)
		tableArray = append(tableArray, info.Name)
		commentArray = append(commentArray, info.Comment)
	}

	// filter tables when specified tables params
	if tableName != "*" {
		tableCollect = nil
		chooseTables := strings.Split(tableName, ",")
		indexMap := make(map[int]int)
		for _, item := range chooseTables {
			subIndexMap := getTargetIndexMap(tableArray, item)
			for k, v := range subIndexMap {
				if _, ok := indexMap[k]; ok {
					continue
				}
				indexMap[k] = v
			}
		}

		if len(indexMap) != 0 {
			for _, v := range indexMap {
				var info tableInfo
				info.Name = tableArray[v]
				info.Comment = commentArray[v]
				tableCollect = append(tableCollect, info)
			}
		}
	}

	return tableCollect, err
}

func queryTableColumn(db *gorm.DB, dbName string, tableName string) ([]tableColumn, error) {
	// 定义承载列信息的切片
	var columns []tableColumn

	sqlTableColumn := fmt.Sprintf("SELECT `ORDINAL_POSITION`,`COLUMN_NAME`,`COLUMN_TYPE`,`DATA_TYPE`,`COLUMN_KEY`,`IS_NULLABLE`,`EXTRA`,`COLUMN_COMMENT`,`COLUMN_DEFAULT` FROM `information_schema`.`columns` WHERE `table_schema`= '%s' AND `table_name`= '%s' ORDER BY `ORDINAL_POSITION` ASC",
		dbName, tableName)

	rows, err := db.Raw(sqlTableColumn).Rows()
	if err != nil {
		fmt.Printf("execute query table column action error, detail is [%v]\n", err.Error())
		return columns, err
	}
	defer rows.Close()

	for rows.Next() {
		var column tableColumn
		err = rows.Scan(
			&column.OrdinalPosition,
			&column.ColumnName,
			&column.ColumnType,
			&column.DataType,
			&column.ColumnKey,
			&column.IsNullable,
			&column.Extra,
			&column.ColumnComment,
			&column.ColumnDefault)
		if err != nil {
			fmt.Printf("query table column scan error, detail is [%v]\n", err.Error())
			return columns, err
		}
		columns = append(columns, column)
	}

	return columns, err
}

func getTargetIndexMap(tableNameArr []string, item string) map[int]int {
	indexMap := make(map[int]int)
	for i := 0; i < len(tableNameArr); i++ {
		if tableNameArr[i] == item {
			if _, ok := indexMap[i]; ok {
				continue
			}
			indexMap[i] = i
		}
	}
	return indexMap
}

// capitalize
func capitalize(s string) string {
	var upperStr string
	chars := strings.Split(s, "_")
	for _, val := range chars {
		vv := []rune(val)
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				if vv[i] >= 97 && vv[i] <= 122 {
					vv[i] -= 32
				}
				upperStr += string(vv[i])
			} else {
				upperStr += string(vv[i])
			}
		}
	}
	return upperStr
}

func textType(s string) string {
	var databaseTypeToGoType = map[string]dataTypeMapping{
		"numeric":    func(string) string { return "int32" },
		"integer":    func(string) string { return "int32" },
		"int":        func(string) string { return "int32" },
		"smallint":   func(string) string { return "int32" },
		"mediumint":  func(string) string { return "int32" },
		"bigint":     func(string) string { return "int64" },
		"float":      func(string) string { return "float32" },
		"real":       func(string) string { return "float64" },
		"double":     func(string) string { return "float64" },
		"decimal":    func(string) string { return "float64" },
		"char":       func(string) string { return "string" },
		"varchar":    func(string) string { return "string" },
		"tinytext":   func(string) string { return "string" },
		"mediumtext": func(string) string { return "string" },
		"longtext":   func(string) string { return "string" },
		"binary":     func(string) string { return "[]byte" },
		"varbinary":  func(string) string { return "[]byte" },
		"tinyblob":   func(string) string { return "[]byte" },
		"blob":       func(string) string { return "[]byte" },
		"mediumblob": func(string) string { return "[]byte" },
		"longblob":   func(string) string { return "[]byte" },
		"text":       func(string) string { return "string" },
		"json":       func(string) string { return "string" },
		"enum":       func(string) string { return "string" },
		"time":       func(string) string { return "time.Time" },
		"date":       func(string) string { return "time.Time" },
		"datetime":   func(string) string { return "time.Time" },
		"timestamp":  func(string) string { return "time.Time" },
		"year":       func(string) string { return "int32" },
		"bit":        func(string) string { return "[]uint8" },
		"boolean":    func(string) string { return "bool" },
		"tinyint": func(detailType string) string {
			if strings.HasPrefix(strings.TrimSpace(detailType), "tinyint(1)") {
				return "bool"
			}
			return "int32"
		},
	}

	if convert, ok := databaseTypeToGoType[strings.ToLower(s)]; ok {
		return convert(s)
	}
	return "string"
}
