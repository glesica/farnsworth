package proxy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"bufio"
)

var (
	proxyRegistryLock sync.Mutex
	proxyRegistry     = []*proxyRegistryEntry{}
)

type proxyRegistryEntry struct {
	factory   Factory
	validator Validator
}

// Proxy is a project type interface. For instance, a Java project.
type Proxy interface {
	IsHideLine(line string) bool
	IsStopLine(line string) bool
	Name() string
	ShouldMerge(path string, content io.Reader) bool
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

// Get returns an instance of the correct proxy for the project
//rooted at the given path.
func Get(rootPath string) (Proxy, error) {
	for _, entry := range proxyRegistry {
		if entry.validator(rootPath) {
			return entry.factory(), nil
		}
	}

	return nil, fmt.Errorf("path '%s' is not a valid project root", rootPath)
}

// Register adds a given factory to be considered when making a call
// to GetProxy. The factory will be invoked if the associated validator returns
// true.
func Register(factory Factory, validator Validator) {
	proxyRegistryLock.Lock()
	defer proxyRegistryLock.Unlock()

	proxyRegistry = append(proxyRegistry, &proxyRegistryEntry{factory, validator})
}
