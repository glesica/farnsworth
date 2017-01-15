package project

// Type is a project type interface. For instance, a Java project.
type Type interface {
	Path() string
	IsProject(proj Project) bool
	ShouldMerge(path string) bool
	IsHideLine(line string) bool
	IsStopLine(line string) bool
}
