// Package migrations embeds all SQL migration files so they can be
// compiled into the binary and applied automatically on startup.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
