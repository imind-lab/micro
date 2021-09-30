package orm

const templateContent = `{{ $TableName := .BigCamelName -}}
package model

type {{ $TableName }} struct {
{{ range $index, $field := .Fields }}{{"\t"}}{{ $field.BigCamelName }}{{ range $field.BigCamelSpaces }} {{ end }}{{ $field.DataType }}{{ range $field.TypeSpaces }} {{ end }}${backquote}gorm:"column:{{ $field.Name }}{{ if eq $field.ColKey "PRI" }};primary_key{{ end }}" redis:"{{ $field.Name }}{{ if ne $field.ColKey "PRI" }},omitempty{{ end }}"${backquote} {{ $le:= len $field.Comment }}{{ if gt $le 0 }}// {{ $field.Comment }}{{ end }}
{{ end }}}

func (t *{{ $TableName }}) TableName() string {
	return "{{ .Name }}"
}
`
