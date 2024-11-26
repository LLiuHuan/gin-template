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
	Name    string         `gorm:"column:TABLE_NAME"`    // 表名
	Comment sql.NullString `gorm:"column:TABLE_COMMENT"` // 表描述
}

// tableColumn 表字段信息
type tableColumn struct {
	OrdinalPosition uint16         `gorm:"column:ORDINAL_POSITION"` // 列在表中的位置（从1开始）, 表示该列是表中第几列
	ColumnName      string         `gorm:"column:COLUMN_NAME"`      // 列的名称
	ColumnType      string         `gorm:"column:COLUMN_TYPE"`      // 列的完整数据类型定义，包括长度、是否为 UNSIGNED、是否为 ZEROFILL 等信息, 例如 varchar(255) int(11) unsigned
	DataType        string         `gorm:"column:DATA_TYPE"`        // 列的数据类型（不带长度）。常见的数据类型包括 int、varchar、date、text 等
	IsNullable      string         `gorm:"column:IS_NULLABLE"`      // 该列是否允许存储 NULL 值。值为 YES 或 NO
	ColumnKey       sql.NullString `gorm:"column:COLUMN_KEY"`       // 列在表中的索引类型。可能的值包括： PRI: 主键。 UNI: 唯一索引。 MUL: 可以重复的非唯一索引（普通索引）。
	ColumnDefault   sql.NullString `gorm:"column:COLUMN_DEFAULT"`   // 列的默认值。可以是 NULL 或其他常量。如果列没有设置默认值，这里将显示 NULL。
	Extra           sql.NullString `gorm:"column:EXTRA"`            // 列的额外信息，如自动递增（auto_increment）或其他特殊属性。
	Privileges      string         `gorm:"column:PRIVILEGES"`       // 列的权限信息，如 SELECT、INSERT、UPDATE、REFERENCES 等
	ColumnComment   sql.NullString `gorm:"column:COLUMN_COMMENT"`   // 列的注释，可以由用户通过 COMMENT 子句定义。
}

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
		log.Fatal("new gorm DB err", err)
	}

	defer func() {
		if err := db.DBClose(); err != nil {
			log.Println("gorm DB close err", err)
		}
	}()

	tables, err := queryTables(db.GetDB(), dbName, genTables)
	if err != nil {
		log.Println("query tables of gorm DB err", err)
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

			columnJson := "`gorm:\""
			if info.ColumnKey.String == "PRI" {
				columnJson += "primaryKey;"
			}
			if info.Extra.String == "auto_increment" {
				columnJson += "autoIncrement;"
			}
			columnJson += "column:" + info.ColumnName + ";type:" + info.ColumnType + ";"
			if info.IsNullable == "NO" && info.ColumnDefault.Valid {
				columnJson += "default:" + info.ColumnDefault.String + ";"
			}
			if info.IsNullable == "NO" {
				columnJson += "NOT NULL;"
			} else {
				columnJson += "NULL;"
			}
			if info.ColumnComment.Valid {
				columnJson += "comment:" + info.ColumnComment.String
			}
			columnJson += "\" json:\"" + info.ColumnName + "\"`"

			// 代码注释
			columnComment := ""
			if info.ColumnComment.Valid {
				columnComment = "// " + info.ColumnComment.String
			}

			// TODO:
			modelContent += fmt.Sprintf(
				"%s %s %s %s\n",
				capitalize(info.ColumnName),
				textType(info.DataType, info.ColumnType, info.IsNullable == "YES"),
				columnJson,
				columnComment,
			)
		}

		mdFile.WriteString(tableContent)
		mdFile.Close()

		modelContent += "}\n"
		modelContent += fmt.Sprintf("\nfunc (*%s) TableName() string {\n\treturn \"%s\"\n}\n", capitalize(table.Name), table.Name)
		modelFile.WriteString(modelContent)
		modelFile.Close()

	}

}

// queryTables 查询表信息
func queryTables(db *gorm.DB, dbName string, tableName string) ([]tableInfo, error) {
	var tables []tableInfo
	tx := db.Raw("SELECT TABLE_NAME,TABLE_COMMENT FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ?", dbName).Scan(&tables)
	if tx.Error != nil {
		return tables, tx.Error
	}

	// 当指定表参数时过滤表
	if tableName != "*" {
		tableCollect := make([]tableInfo, 0)
		chooseTables := strings.Split(tableName, ",")

		for _, info := range tables {
			for _, item := range chooseTables {
				if info.Name == item {
					tableCollect = append(tableCollect, info)
				}
			}
		}

		return tableCollect, nil
	}

	return tables, nil
}

// queryTableColumn 查询表字段信息
func queryTableColumn(db *gorm.DB, dbName string, tableName string) ([]tableColumn, error) {
	// 定义承载列信息的切片
	columns := make([]tableColumn, 0)

	tx := db.Raw(`
		SELECT 
		ORDINAL_POSITION,COLUMN_NAME,COLUMN_TYPE,DATA_TYPE,IS_NULLABLE,COLUMN_KEY,COLUMN_DEFAULT,EXTRA,PRIVILEGES,COLUMN_COMMENT 
		FROM information_schema.columns 
		WHERE table_schema= ? AND table_name= ? ORDER BY ORDINAL_POSITION
		`, dbName, tableName).Scan(&columns)
	if tx.Error != nil {
		fmt.Printf("execute query table column action error, detail is [%v]\n", tx.Error.Error())
		return columns, tx.Error
	}
	return columns, nil
}

// capitalize 首字母大写
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

// parseAny2Ptr 解析数据类型,如果是可为空,则返回指针类型
func parseAny2Ptr(isNull bool, dataType string) string {
	if isNull {
		return "*" + dataType
	}
	return dataType
}

// textType 获取数据库类型
func textType(dataType, columnType string, isNull bool) string {
	var unsigned string
	if strings.Contains(columnType, "unsigned") {
		unsigned = "u"
	}

	switch dataType {
	case "int", "integer", "mediumint", "year":
		return parseAny2Ptr(isNull, unsigned+"int")
	case "tinyint":
		if strings.HasPrefix(strings.TrimSpace(columnType), "tinyint(1)") {
			return "bool"
		}
		return parseAny2Ptr(isNull, unsigned+"int8")
	case "smallint":
		return parseAny2Ptr(isNull, unsigned+"int16")
	case "bigint":
		return parseAny2Ptr(isNull, unsigned+"int64")
	case "double", "float", "real", "numeric":
		if unsigned == "u" {
			return "float64"
		}
		return "float32"
	case "decimal":
		//return "float64"
		return "decimal.Decimal"
	case "timestamp", "datetime", "date", "time":
		return parseAny2Ptr(isNull, unsigned+"time.Time")
	case "bool", "boolean":
		return "bool"
	default:
		return parseAny2Ptr(isNull, "string")
	}
}
