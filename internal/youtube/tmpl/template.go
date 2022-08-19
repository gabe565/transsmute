package tmpl

import (
	_ "embed"
	"github.com/gabe565/transsmute/internal/template_funcs"
	"html/template"
)

//go:embed description.html.gotmpl
var descriptionTmplStr string

var DescriptionTmpl = template.Must(
	template.New("").Funcs(template_funcs.FuncMap()).Parse(descriptionTmplStr),
)
