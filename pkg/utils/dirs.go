package utils

import (
	"os"
	"path/filepath"

	"our_records/internal/config"
)

// EnsureUploadDirs 确保上传目录存在
func EnsureUploadDirs() error {
	dirs := []string{
		filepath.Join(config.AppConfig.Upload.Path, "images"),
		filepath.Join(config.AppConfig.Upload.Path, "audios"),
		filepath.Join(config.AppConfig.Upload.Path, "videos"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
