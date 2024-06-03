package template_funcs

import (
	"html/template"

	"github.com/Masterminds/sprig/v3"
)

func FuncMap() template.FuncMap {
	funcs := sprig.FuncMap()

	funcs["escape"] = Escape
	funcs["nl2br"] = Nl2br
	funcs["formatUrls"] = FormatUrls
	funcs["formatHashtags"] = FormatHashtags
	funcs["formatTimestamps"] = FormatTimestamps
	funcs["html"] = Html

	return funcs
}
