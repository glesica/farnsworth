package proxy

import (
	"fmt"

	"github.com/glesica/farnsworth/java"
)

// Proxy is a project type interface. For instance, a Java project.
type Proxy interface {
	ShouldMerge(path string) bool
	IsHideLine(line string) bool
	IsStopLine(line string) bool
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
	return nil, fmt.Errorf("path '%s' is not a valid project root", path)
}
