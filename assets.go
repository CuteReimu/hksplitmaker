package main

import "embed"

//go:embed translate.tsv
var transLateData []byte

//go:embed hk-split-maker/src/asset/hollowknight
var assets embed.FS
