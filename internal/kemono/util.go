package kemono

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//nolint:gochecknoglobals
var serviceNameReplacer = strings.NewReplacer(
	"fans", "Fans",
	"star", "Star",
)

func formatServiceName(name string) string {
	caser := cases.Title(language.English)
	name = caser.String(name)
	name = serviceNameReplacer.Replace(name)
	return name
}
