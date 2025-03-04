package dirMigrite

import (
	"fmt"
	"os"
	"path/filepath"
)

func PathMigrite() (string, error) {

	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	migritePath := filepath.Join(exeDir, "../internal/dirMigrite")

	absMigritePath, err := filepath.Abs(migritePath)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path for migrite: %v", err)
	}

	return absMigritePath, nil
}
