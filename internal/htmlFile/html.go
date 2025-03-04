package htmlFile

import (
	"fmt"
	"os"
	"path/filepath"
)

func PathHtml(s string) (string, error) {

	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	htmlPath := filepath.Join(exeDir, "../internal/htmlFile/", s)

	absHTMLPath, err := filepath.Abs(htmlPath)
	if err != nil {
		return "", fmt.Errorf("cannot resolve absolute path for '%s': %w", htmlPath, err)
	}

	return absHTMLPath, nil
}
