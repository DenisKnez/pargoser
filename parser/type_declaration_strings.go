package parser

const StructTemplateString = "type {{ .Name }} struct { \n" +
	"{{ range .Fields }}  {{ if .Name }}{{ .Name }} {{ .Type }} {{ else }}{{ .Type }}{{ end }}  " +
	"{{ if .Tag }}`{{ .Tag.Value.Type }}:\"{{ .Tag.Value.Value }}\"`{{ end }}\n" +
	"{{ end }}\n" +
	"}\n"

const FunctionTemplateString = "r"

const InterfaceTemplateString = "type {{ .Name }} interface { \n" +
	"{{ range .Methods }}  {{ .Name }}({{ range .Params }}{{ .Name }} {{ .Type }}{{ end }})" +
	"{{ range .Results }}{{ .Name }} {{ .Type }} {{ end }}\n{{ end }}" +
	"}\n"

const ImportTemplateString = "{{ .Name }} {{ .Path }}"
