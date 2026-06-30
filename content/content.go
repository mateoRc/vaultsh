package content

import "embed"

//go:embed .motd cv/*.txt cv/experience/*.txt docs/*.md projects/*.txt
var Files embed.FS
