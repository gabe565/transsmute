package assets

import (
	"embed"
	_ "embed"
)

//go:embed favicon.ico
var Assets embed.FS
