package templates

import _ "embed"

//go:embed index.html
var Index string

//go:embed thanks.html
var Thanks string

// Confirm template parameters:
//
//	token
//
//go:embed confirm.html
var Confirm string

//go:embed confirmed.html
var Confirmed string
