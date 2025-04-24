package public

import (
	"embed"
)

//go:embed * */*
var fs embed.FS

func EFS() *embed.FS {
	return &fs
}
