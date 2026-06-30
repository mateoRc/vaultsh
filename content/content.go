package content

import _ "embed"

// Database is the read-only content database shipped with the binary.
//
//go:embed vaultsh.db
var Database []byte
