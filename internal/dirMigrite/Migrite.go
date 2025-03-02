package dirMigrite

import (
	"os"
	"path/filepath"
)

func PathMigrite() (string, error) {

	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)

	migritePath := filepath.Join(exeDir, "../internal/dirMigrite")

	absMigritePath, err := filepath.Abs(migritePath)
	if err != nil {
		return "", err
	}

	return absMigritePath, nil
}
