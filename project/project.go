package project

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"strings"

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

// MergeFrom copies parts of another project into the receiver project.
func (proj *Project) MergeFrom(mergeProj Project) error {
	if proj.Name() != mergeProj.Name() {
		return fmt.Errorf(
			"cannot merge project of type '%s' into project of type '%s'",
			mergeProj.Name(),
			proj.Name())
	}

	filepath.Walk(mergeProj.path, func(
		filePath string,
		fileInfo os.FileInfo,
		walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("failed to stat '%s', skipping", filePath)
		}

		if fileInfo.IsDir() {
			return nil
		}

		fileContent, fileContentErr := ioutil.ReadFile(filePath)
		if fileContentErr != nil {
			return fmt.Errorf("failed to read '%s'", filePath)
		}

		if !proj.ShouldMerge(filePath, fileContent) {
			return nil
		}

		// Merge the file

		destFilePath := strings.Replace(filePath, mergeProj.path, proj.path, 1)

		destFile, destFileErr := os.Create(destFilePath)
		if destFileErr != nil {
			return fmt.Errorf("failed to open '%s'", destFilePath)
		}

		_, destWriteErr := destFile.Write(fileContent)
		if destWriteErr != nil {
			return fmt.Errorf("failed to write to '%s'", destFilePath)
		}

		return nil
	})

	return nil
}

// Path returns the absolute filesystem path to the project root.
func (proj *Project) Path() string {
	return proj.path
}

// Zip compresses a project into a Zip archive.
func (proj *Project) Zip(zipPath string, private bool) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create archive '%s'", zipPath)
	}

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	defer zipFile.Close()

	zipInfo, _ := os.Stat(zipPath)

	filepath.Walk(proj.path, func(filePath string, fileInfo os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("failed to stat '%s', skipping", filePath)
		}

		// Don't accidentally try to zip the zip file.
		if os.SameFile(zipInfo, fileInfo) {
			return nil
		}

		if fileInfo.IsDir() {
			return nil
		}

		var fileContent []byte

		if private {
			content, contentErr := proxy.RemoveHiddenLinesFromFile(filePath, proj)
			if contentErr != nil {
				return contentErr
			}
			fileContent = []byte(content)
		} else {
			content, contentErr := ioutil.ReadFile(filePath)
			if contentErr != nil {
				return fmt.Errorf("failed to read file '%s'", filePath)
			}
			fileContent = content
		}

		relFilePath, relFilePathErr := filepath.Rel(proj.path, filePath)
		if relFilePathErr != nil {
			return fmt.Errorf("failed to find relative path to file '%s'", filePath)
		}

		f, err := zipWriter.Create(path.Join(proj.baseName, relFilePath))
		if err != nil {
			return fmt.Errorf("failed add file to zip %s", filePath)
		}

		_, err = f.Write(fileContent)
		if err != nil {
			return fmt.Errorf("failed to write '%s' contents to zip", filePath)
		}

		return nil
	})

	return nil
}
