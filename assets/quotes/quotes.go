package quotes

import (
	"embed"
)

//go:embed quotes.txt
var Quotes embed.FS
