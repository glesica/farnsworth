package proxy

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"bufio"

	"github.com/glesica/farnsworth/proxy/java"
	"github.com/glesica/farnsworth/proxy/go"
)

// Proxy is a project type interface. For instance, a Java project.
type Proxy interface {
	IsHideLine(line string) bool
	IsStopLine(line string) bool
	Name() string
	ShouldMerge(path string, content []byte) bool
}

// RemoveHiddenLines returns the contents of a Reader with all hidden
// lines removed, based on the provided Proxy instance.
func RemoveHiddenLines(fileContent io.Reader, proxy Proxy) (string, error) {
	fileScanner := bufio.NewScanner(fileContent)

	lineNumber := 0
	isHiding := false
	filteredFileContentBuffer := bytes.Buffer{}

	for fileScanner.Scan() {
		lineNumber++
		fileLine := fileScanner.Text()

		if proxy.IsHideLine(fileLine) {
			if isHiding {
				return "", fmt.Errorf("error, line %d, nested 'hide' blocks", lineNumber)
			}
			isHiding = true
		}

		if !isHiding {
			// If this is the first non-hidden line, do not
			// add a leading newline.
			if filteredFileContentBuffer.Len() > 0 {
				filteredFileContentBuffer.WriteString("\n")
			}
			filteredFileContentBuffer.WriteString(fileLine)
		}

		if proxy.IsStopLine(fileLine) {
			if !isHiding {
				return "", fmt.Errorf("error, line %d, dangling 'stop'", lineNumber)
			}
			isHiding = false
		}
	}

	return filteredFileContentBuffer.String(), nil
}

// RemoveHiddenLinesFromFile returns the contents of a file as a string,
// with all hidden lines removed, based on the provided Proxy instance.
func RemoveHiddenLinesFromFile(filePath string, proxy Proxy) (string, error) {
	file, fileErr := os.Open(filePath)
	if fileErr != nil {
		return "", fmt.Errorf("failed to open file '%s'", filePath)
	}
	defer file.Close()

	return RemoveHiddenLines(file, proxy)
}

// Validator indicates whether the project rooted at the given path
// supports a given proxy.
type Validator func(path string) bool

// Factory returns a new Proxy instance.
type Factory func() Proxy

// GetProxy returns an instance of the correct proxy for the project
//rooted at the given path.
func GetProxy(path string) (Proxy, error) {
	if java.IsValid(path) {
		return java.Factory(), nil
	}
	if golang.IsValid(path) {
		return golang.Factory(), nil
	}
	return nil, fmt.Errorf("path '%s' is not a valid project root", path)
}
