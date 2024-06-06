package docker

import (
	_ "embed"
	"html/template"

	"github.com/gabe565/transsmute/internal/templatefuncs"
)

//go:embed description.html.gotmpl
var descriptionTmplStr string

//nolint:gochecknoglobals
var descriptionTmpl = template.Must(
	template.New("").Funcs(templatefuncs.FuncMap()).Parse(descriptionTmplStr),
)

type DescriptionValues struct {
	Repo string
	Tag  string
}
