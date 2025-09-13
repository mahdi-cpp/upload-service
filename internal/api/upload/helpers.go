package upload

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// Helper functions
func isJPEG(file *multipart.FileHeader) bool {
	// Check content type
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" {
		return false
	}

	// Also check file extension for extra safety
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return ext == ".jpg" || ext == ".jpeg"
}

func getJPEGFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".jpg" || ext == ".jpeg" {
				// Return relative path
				rel, err := filepath.Rel(dir, path)
				if err == nil {
					files = append(files, rel)
				}
			}
		}

		return nil
	})

	return files, err
}
