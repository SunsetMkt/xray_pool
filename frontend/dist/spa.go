package dist

import (
	"embed"
	_ "embed"
	"io/fs"
	"net/http"
)

//go:embed index.html
var SpaIndexHtml []byte

//go:embed assets
var SpaJS embed.FS

func Assets(dirName string, emFS embed.FS) http.FileSystem {
	// even uiAssets is empty, fs.Sub won't fail
	stripped, err := fs.Sub(emFS, dirName)
	if err != nil {
		panic(err)
	}
	return http.FS(stripped)
}

const (
	SpaFolderName = "assets"
)
