package content

import "embed"

//go:embed *.txt experience/*.txt projects/*.txt
var Files embed.FS
