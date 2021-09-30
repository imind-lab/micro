package orm

import "strings"

type Field struct {
	Name           string   `gorm:"type:string"` //字段名
	BigCamelName   string   `gorm:"type:string"` //驼峰后字段名
	Comment        string   `gorm:"type:string"` //字段备注
	DataType       string   `gorm:"type:string"` //字段类型
	ColKey         string   `gorm:"type:string"` //主键
	BigCamelSpaces []string `gorm:"type:string"`
	TagGormSpaces  []string `gorm:"type:string"`
	TypeSpaces     []string `gorm:"type:string"`
}

//表字段类型转换成GO相关类型
func TransformType(typeStr string) string {
	switch typeStr {
	case "bit", "int", "tinyint", "small_int", "smallint", "medium_int", "mediumint":
		if strings.Contains(typeStr, "unsigned") {
			return "uint32"
		} else {
			return "int32"
		}
	case "big_int", "bigint":
		if strings.Contains(typeStr, "unsigned") {
			return "uint64"
		} else {
			return "int64"
		}
	case "varchar", "text", "char":
		return "string"
	case "float", "double", "decimal":
		return "float64"
	case "bool":
		return "bool"
	case "datetime", "timestamp", "date", "time":
		return "string"
	default:
		return "string"
	}
}
