package static

import (
	"embed"
	"io/fs"
)

//go:embed dist
var Fs embed.FS

var RootFs fs.FS

func GetFs() fs.FS {
	if RootFs == nil {
		RootFs, _ = fs.Sub(Fs, "dist")
	}
	return RootFs
}
