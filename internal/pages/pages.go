package pages

import (
	"path/filepath"
	"runtime"
)

func GetHtmlTemplate() string {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		panic("não foi possível obter o caminho do arquivo")
	}

	dir := filepath.Dir(filename)
	path := filepath.Join(dir, "./index.html")

	return path
}
