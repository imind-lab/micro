package orm

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// @summary table 生成 struct
// @param   db string "conf.yaml文件中db组件中对应数据库链接的key值"
// @param   args []  "需要生成表"
// @param   modelPath string "生成目录"
// @example
// table := []string{"xpt_user","xpt_order"}
// gen.Run("xpt", table, "../models/xpt/")
func Run(db string, args []string, modelPath string) (err error) {
	tables, err := GetTablesInfo(db, args)
	if err != nil {
		return err
	}
	for _, table := range tables {
		modelStr, err := genGo(table)
		if err != nil {
			return err
		}
		_, err = os.Open(modelPath)
		if err != nil {
			err = os.MkdirAll(modelPath, 0755)
			if err != nil {
				return err
			}
		}
		err = ioutil.WriteFile(modelPath+table.Name+".go", []byte(modelStr), 0644)
		if err != nil {
			return err
		}
	}
	return
}

//模板解析
func genGo(table Table) (string, error) {
	// 解析 model
	str := strings.Replace(templateContent, "${backquote}", "`", -1)
	modelFiles, err := template.New("tp").Parse(str)
	if err != nil {
		return "", err
	}
	var modelsBuf bytes.Buffer
	err = modelFiles.Execute(&modelsBuf, table)
	if err != nil {
		return "", err
	}
	return modelsBuf.String(), nil
}
