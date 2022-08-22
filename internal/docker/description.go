package docker

import (
	_ "embed"
	"github.com/gabe565/transsmute/internal/template_funcs"
	"html/template"
)

//go:embed description.html.gotmpl
var descriptionTmplStr string

var descriptionTmpl = template.Must(
	template.New("").Funcs(template_funcs.FuncMap()).Parse(descriptionTmplStr),
)

type DescriptionValues struct {
	Repo string
	Tag  string
}
