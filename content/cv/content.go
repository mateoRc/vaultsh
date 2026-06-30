package content

import "embed"

//go:embed .motd *.txt experience/*.txt projects/*.txt
var Files embed.FS
