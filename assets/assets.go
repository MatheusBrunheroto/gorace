package assets

import (
	_ "embed"
)

//go:embed help.txt
var Help string

//go:embed logo.txt
var Logo string

//go:embed themes.txt
var Themes string
