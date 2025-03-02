package htmlFile

import (
	"os"
	"path/filepath"
)

func PathHtml(s string) (string, error) {

	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)

	htmlPath := filepath.Join(exeDir, "../internal/htmlFile/", s)

	absHTMLPath, err := filepath.Abs(htmlPath)
	if err != nil {
		return "", err
	}

	return absHTMLPath, nil
}
