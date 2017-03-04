package project

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/jhoonb/archivex"

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

func (proj *Project) filterHidden(filePath string) ([]byte, error) {
	file, fileErr := os.Open(filePath)
	if fileErr != nil {
		return nil, fmt.Errorf("failed to open file '%s'", filePath)
	}
	defer file.Close()

	lineNumber := 0
	hiding := false
	fileContentBuffer := bytes.Buffer{}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		lineNumber++
		line := fileScanner.Text()
		if proj.IsHideLine(line) {
			if hiding {
				return nil, fmt.Errorf("syntax error, line %d, nested hidden blocks", lineNumber)
			}
			hiding = true
		}
		if !hiding {
			fileContentBuffer.WriteString(line)
			fileContentBuffer.WriteString("\n")
		}
		if proj.IsStopLine(line) {
			if !hiding {
				return nil, fmt.Errorf("syntax error, line %d, dangling stop", lineNumber)
			}
			hiding = false
		}
	}

	return fileContentBuffer.Bytes(), nil
}

// Merge copies parts of one project into another.
// TODO: Make sure the merge path project is validated?
func (proj *Project) Merge(mergeProj Project) error {
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
	zipFile := new(archivex.ZipFile)
	zipErr := zipFile.Create(zipPath)
	if zipErr != nil {
		return fmt.Errorf("failed to create archive '%s'", zipPath)
	}
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
			content, contentErr := proj.filterHidden(filePath)
			if contentErr != nil {
				return contentErr
			}
			fileContent = content
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

		zipFile.Add(path.Join(proj.baseName, relFilePath), fileContent)
		return nil
	})

	return nil
}
