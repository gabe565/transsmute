package tmpl

import (
	_ "embed"
	"html/template"

	"gabe565.com/transsmute/internal/templatefuncs"
)

//go:embed description.html.gotmpl
var descriptionTmplStr string

//nolint:gochecknoglobals
var DescriptionTmpl = template.Must(
	template.New("").Funcs(templatefuncs.FuncMap()).Parse(descriptionTmplStr),
)
