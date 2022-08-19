package template_funcs

import "html/template"

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"escape":           Escape,
		"nl2br":            Nl2br,
		"formatUrls":       FormatUrls,
		"formatHashtags":   FormatHashtags,
		"formatTimestamps": FormatTimestamps,
		"html":             Html,
	}
}
