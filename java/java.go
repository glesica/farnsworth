package java

import (
	"io/ioutil"
)

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
func Factory() *Proxy {
	return &Proxy{}
}

// Proxy is a project proxy for a Gradle-based Java project.
type Proxy struct{}

// Name is the unique name of the project proxy.
func (proxy *Proxy) Name() string {
	return "java"
}

// IsHideLine indicates whether the given line begins a hidden block.
func (proxy *Proxy) IsHideLine(line string) bool {
	return false
}

// IsStopLine indicates whether the given line ends a block.
func (proxy *Proxy) IsStopLine(line string) bool {
	return false
}

// ShouldMerge indicates whether the given path should be merged.
func (proxy *Proxy) ShouldMerge(path string, content []byte) bool {
	return true
}
