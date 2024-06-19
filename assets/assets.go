package assets

import "embed"

//go:embed favicon.ico robots.txt
var Assets embed.FS
