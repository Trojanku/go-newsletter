package templates

import _ "embed"

//go:embed index.html
var Index string

//go:embed thanks.html
var Thanks string
