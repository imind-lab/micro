package orm

import (
	"strings"

	"github.com/imind-lab/micro/dao"
)

type Table struct {
	Dbs                   string  `gorm:"type:string"` //数据库链接的key值
	Name                  string  `gorm:"type:string"` //表名
	BigCamelName          string  `gorm:"type:string"` //驼峰后表名
	Fields                []Field `gorm:"type:string"` //字段
	LongestBigCamelColLen int     `gorm:"type:int"`    //字段空格数长度
	LongestTagGORMLen     int     `gorm:"type:int"`    //gorm空格数长度
	LongestTypeLen        int     `gorm:"type:int"`    //类型空格数长度
}

// 获取表相的相关信息
func GetTablesInfo(db string, Tables []string) ([]Table, error) {
	connection := dao.NewDatabase().DB(db)
	var tables []Table
	query := connection.Table("information_schema.tables").Select("table_name name")
	if len(Tables) > 0 {
		query = query.Where("table_schema = ? AND table_name IN (?)", db, Tables)
	} else {
		query = query.Where("table_schema = ? ", db)
	}
	query.Scan(&tables)
	// 查出每个表的字段信息
	for i, table := range tables {
		var columns []Field
		connection.Table("information_schema.columns").
			Select("COLUMN_NAME name, COLUMN_KEY col_key, COLUMN_COMMENT comment, DATA_TYPE data_type ").
			Where("table_schema = ? AND table_name = ?", db, table.Name).
			Order("ordinal_position ASC").Scan(&columns)
		tables[i].Fields = columns
	}
	return InitTables(db, tables), nil
}

// 初始化表的相关数据，使其能匹配模板进行替换
func InitTables(db string, tables []Table) []Table {
	for i := range tables {
		longestBigCamelColLen, longestTagGORMLen, longestTypeLen := 0, 0, 0
		tables[i].Dbs = db
		tables[i].BigCamelName = ToBigCamelCase(tables[i].Name)
		for j := range tables[i].Fields {
			tables[i].Fields[j].BigCamelName = ToBigCamelCase(tables[i].Fields[j].Name)
			tables[i].Fields[j].DataType = TransformType(tables[i].Fields[j].DataType)
			tagGORMLen := len(tables[i].Fields[j].Name)
			if tables[i].Fields[j].ColKey == "PRI" {
				tagGORMLen += len(";primary_key")
			}
			longestBigCamelColLen = MaxFunc(longestBigCamelColLen, len(tables[i].Fields[j].BigCamelName))
			longestTagGORMLen = MaxFunc(longestTagGORMLen, tagGORMLen)
			longestTypeLen = MaxFunc(longestTypeLen, len(tables[i].Fields[j].DataType))
		}
		tables[i].LongestBigCamelColLen = longestBigCamelColLen
		tables[i].LongestTagGORMLen = longestTagGORMLen
		tables[i].LongestTypeLen = longestTypeLen
	}
	for i := range tables {
		for j := range tables[i].Fields {
			tagGORMLen := len(tables[i].Fields[j].Name)
			if tables[i].Fields[j].ColKey == "PRI" {
				tagGORMLen += len(";primary_key")
			}
			tables[i].Fields[j].BigCamelSpaces = make([]string, tables[i].LongestBigCamelColLen-len(tables[i].Fields[j].BigCamelName)+1)
			tables[i].Fields[j].TagGormSpaces = make([]string, tables[i].LongestTagGORMLen-tagGORMLen+1)
			tables[i].Fields[j].TypeSpaces = make([]string, tables[i].LongestTypeLen-len(tables[i].Fields[j].DataType)+1)
		}
	}
	return tables
}

// 字符串转换为驼峰写法
func ToBigCamelCase(str string) string {
	result := ""
	if str == "" {
		return result
	}
	strs := strings.Split(str, "_")
	for _, s := range strs {
		r := []rune(s)
		if len(r) > 0 {
			if r[0] >= 'a' && r[0] <= 'z' {
				r[0] -= 32
			}
			result += string(r)
		}
	}
	return result
}
func MaxFunc(i, j int) int {
	if i > j {
		return i
	}
	return j
}
