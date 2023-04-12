package docker

import (
	_ "embed"
	"html/template"

	"github.com/gabe565/transsmute/internal/template_funcs"
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
