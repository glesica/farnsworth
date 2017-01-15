package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/jhoonb/archivex"
)

// A Project is a Farnsworth project.
type Project struct {
	BaseName string
	Path     string
}

// Load creates a new project from a path.
func Load(projectPath string) (*Project, error) {
	absProjectPath, absProjectPathErr := filepath.Abs(projectPath)
	if absProjectPathErr != nil {
		return nil, fmt.Errorf("failed to convert '%s' to absolute path", projectPath)
	}

	projectBaseName := filepath.Base(absProjectPath)

	project := Project{
		BaseName: projectBaseName,
		Path:     absProjectPath,
	}

	return &project, nil
}

// Zip compresses a project into a Zip archive.
func (proj Project) Zip(zipPath string) error {
	zipFile := new(archivex.ZipFile)
	zipFile.Create(zipPath)

	zipInfo, _ := os.Stat(zipPath)

	filepath.Walk(proj.Path, func(filePath string, fileInfo os.FileInfo, walkErr error) error {
		if walkErr != nil {
			fmt.Fprintf(os.Stderr, "Error reading, skipping '%s'", filePath)
			return nil
		}

		// Don't accidentally try to zip the zip file.
		if os.SameFile(zipInfo, fileInfo) {
			return nil
		}

		if fileInfo.IsDir() {
			return nil
		}

		fileContent, fileContentErr := ioutil.ReadFile(filePath)
		if fileContentErr != nil {
			return fmt.Errorf("failed to read file '%s'", filePath)
		}

		newFilePath, newFilePathErr := filepath.Rel(filepath.Dir(proj.Path), filePath)
		if newFilePathErr != nil {
			return fmt.Errorf("Failed to find relative path of '%s'", filePath)
		}

		zipFile.Add(path.Join(proj.BaseName, newFilePath), fileContent)
		return nil
	})

	zipFile.Close()

	return nil
}
