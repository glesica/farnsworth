package java

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/glesica/farnsworth/proxy"
)

func init() {
	proxy.Register(Factory, IsValid)
}

// IsValid indicates whether the given project root is a valid project of this type.
func IsValid(path string) bool {
	dirEntries, dirEntriesErr := ioutil.ReadDir(path)
	if dirEntriesErr != nil {
		return false
	}

	for _, entry := range dirEntries {
		if entry.Name() == "build.gradle" {
			return true
		}
	}

	return false
}

// Factory returns an instance of the Proxy.
func Factory() proxy.Proxy {
	return &java{}
}

// java is a project proxy for a Gradle-based Java project.
type java struct{}

// Name is the unique name of the project proxy.
func (proxy *java) Name() string {
	return "java"
}

// IsHideLine indicates whether the given line begins a hidden block.
func (proxy *java) IsHideLine(line string) bool {
	matched, matchedErr := regexp.MatchString(`^\s*//\+\+\s*hide\s*$`, line)
	if matchedErr != nil {
		// Dangerous, but meh.
		return false
	}
	return matched
}

// IsStopLine indicates whether the given line ends a block.
func (proxy *java) IsStopLine(line string) bool {
	matched, matchedErr := regexp.MatchString(`^\s*//\+\+\s*stop\s*$`, line)
	if matchedErr != nil {
		// Dangerous, but meh.
		return false
	}
	return matched
}

// ShouldMerge indicates whether the given path should be merged.
func (proxy *java) ShouldMerge(path string, content []byte) bool {
	return strings.Contains(path, "src/test")
}
