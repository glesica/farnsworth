package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/jhoonb/archivex"

	"github.com/glesica/farnsworth/proxy"
)

// A Project is a Farnsworth project.
type Project struct {
	proxy.Proxy

	baseName string
	path     string
}

// Load creates a new project from a path.
func Load(projectPath string) (*Project, error) {
	absProjectPath, absProjectPathErr := filepath.Abs(projectPath)
	if absProjectPathErr != nil {
		return nil, fmt.Errorf("failed to convert '%s' to absolute path", projectPath)
	}

	projectBaseName := filepath.Base(absProjectPath)

	projectProxy, projectProxyErr := proxy.GetProxy(projectPath)
	if projectProxyErr != nil {
		return nil, projectProxyErr
	}

	project := Project{
		Proxy:    projectProxy,
		baseName: projectBaseName,
		path:     absProjectPath,
	}

	return &project, nil
}

// BaseName returns the name of the directory that contains the project root.
func (proj *Project) BaseName() string {
	return proj.baseName
}

// Path returns the absolute filesystem path to the project root.
func (proj *Project) Path() string {
	return proj.path
}

// Zip compresses a project into a Zip archive.
func (proj *Project) Zip(zipPath string) error {
	zipFile := new(archivex.ZipFile)
	zipFile.Create(zipPath)

	zipInfo, _ := os.Stat(zipPath)

	filepath.Walk(proj.path, func(filePath string, fileInfo os.FileInfo, walkErr error) error {
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

		relFilePath, relFilePathErr := filepath.Rel(proj.path, filePath)
		if relFilePathErr != nil {
			return fmt.Errorf("Failed to find relative path of '%s'", filePath)
		}

		zipFile.Add(path.Join(proj.baseName, relFilePath), fileContent)
		return nil
	})

	zipFile.Close()

	return nil
}
