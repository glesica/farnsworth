package project

import (
	"archive/zip"
	"fmt"
	"io"
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

// A ProxyFactory creates and returns the appropriate kind of
// Proxy for a given path.
type ProxyFactory func(projectPath string) (proxy.Proxy, error)

func loadWithFactory(projectPath string, proxyFactory ProxyFactory) (*Project, error) {
	absProjectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to convert '%s' to absolute path", projectPath)
	}

	projectBaseName := filepath.Base(absProjectPath)

	proxy, err := proxyFactory(projectPath)
	if err != nil {
		return nil, err
	}

	project := Project{
		Proxy:    proxy,
		baseName: projectBaseName,
		path:     absProjectPath,
	}

	return &project, nil
}

// Load creates a new project from a path.
func Load(projectPath string) (*Project, error) {
	return loadWithFactory(projectPath, proxy.Get)
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
		err error) error {
		if err != nil {
			return fmt.Errorf("failed to stat '%s', skipping", filePath)
		}

		if fileInfo.IsDir() {
			return nil
		}

		fileReader, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to read '%s'", filePath)
		}

		if !proj.ShouldMerge(filePath, fileReader) {
			return nil
		}

		fileReader.Close()

		// Merge the file

		srcFile, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to read '%s'", filePath)
		}
		defer srcFile.Close()

		destFilePath := strings.Replace(filePath, mergeProj.path, proj.path, 1)

		destFile, err := os.Create(destFilePath)
		if err != nil {
			return fmt.Errorf("failed to open '%s'", destFilePath)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
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
	// TODO: This stuff can be abstracted by Archiver.
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create archive '%s'", zipPath)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	zipInfo, _ := os.Stat(zipPath)

	filepath.Walk(proj.path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to stat '%s', skipping", filePath)
		}

		// Don't accidentally try to zip the zip file.
		// TODO: Figure out how to do this check with the Archiver.
		//       Maybe it isn't needed any more if all file operations are
		//       held until after the tree has been walked?
		if os.SameFile(zipInfo, fileInfo) {
			return nil
		}

		if fileInfo.IsDir() {
			return nil
		}

		var fileContent []byte

		// TODO: We can extract this block into a method that operates on a Reader.
		if private {
			content, err := proxy.RemoveHiddenLinesFromFile(filePath, proj)
			if err != nil {
				return err
			}
			fileContent = []byte(content)
		} else {
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file '%s'", filePath)
			}
			fileContent = content
		}

		relFilePath, err := filepath.Rel(proj.path, filePath)
		if err != nil {
			return fmt.Errorf("failed to find relative path to file '%s'", filePath)
		}

		// TODO: We can abstract the zip file for testing.
		writer, err := zipWriter.Create(path.Join(proj.baseName, relFilePath))
		if err != nil {
			return fmt.Errorf("failed add file to zip %s", filePath)
		}

		_, err = writer.Write(fileContent)
		if err != nil {
			return fmt.Errorf("failed to write '%s' contents to zip", filePath)
		}

		return nil
	})

	return nil
}
