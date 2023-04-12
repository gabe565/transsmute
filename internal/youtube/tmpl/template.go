package tmpl

import (
	_ "embed"
	"html/template"

	"github.com/gabe565/transsmute/internal/template_funcs"
)

//go:embed description.html.gotmpl
var descriptionTmplStr string

var DescriptionTmpl = template.Must(
	template.New("").Funcs(template_funcs.FuncMap()).Parse(descriptionTmplStr),
)
