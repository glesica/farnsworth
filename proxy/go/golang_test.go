package golang

import "testing"
import "fmt"

func TestName(t *testing.T) {
	proxy := golang{}
	name := proxy.Name()
	if name != "golang" {
		t.Errorf("Expected Name() to return 'golang', found '%s'", name)
	}
}

func getPassingTags(label string) []string {
	return []string{
		fmt.Sprintf("//++%s", label),
		fmt.Sprintf("//++ %s", label),
		fmt.Sprintf(" //++ %s", label),
		fmt.Sprintf("  //++ %s", label),
		fmt.Sprintf("\t//++ %s", label),
		fmt.Sprintf("\t\t//++ %s", label),
		fmt.Sprintf("//++ %s", label),
		fmt.Sprintf("//++ %s ", label),
		fmt.Sprintf("//++ %s  ", label),
		fmt.Sprintf("//++ %s\t", label),
		fmt.Sprintf("//++ %s\t\t", label),
		fmt.Sprintf("//++\t%s", label),
		fmt.Sprintf("//++\t%s ", label),
		fmt.Sprintf("//++\t%s  ", label),
		fmt.Sprintf("//++\t%s\t", label),
		fmt.Sprintf("//++\t%s\t\t", label),
		fmt.Sprintf("//++ %s\n", label),
	}
}

func getFailingTags(label string) []string {
	return []string{
		fmt.Sprintf("// +%s", label),
		fmt.Sprintf("//+ %s", label),
		fmt.Sprintf("// + %s", label),
		fmt.Sprintf("// ++%s", label),
		fmt.Sprintf("///++ %s", label),
		fmt.Sprintf("/// ++ %s", label),
		fmt.Sprintf("// +++%s", label),
		fmt.Sprintf("//+++ %s", label),
		fmt.Sprintf("// +++ %s", label),
		fmt.Sprintf("// %s", label),
	}
}

func TestIsHideLine(t *testing.T) {
	proxy := golang{}
	// True
	for _, line := range getPassingTags("hide") {
		t.Run(fmt.Sprintf("line='%s'", line), func(t *testing.T) {
			if !proxy.IsHideLine(line) {
				t.Errorf("Expected IsHideLine() to return `true` for '%s'", line)
			}
		})
	}
	// False
	for _, line := range getFailingTags("hide") {
		t.Run(fmt.Sprintf("line='%s'", line), func(t *testing.T) {
			if proxy.IsHideLine(line) {
				t.Errorf("Expected IsHideLine() to return `false` for '%s'", line)
			}
		})
	}
}

func TestIsStopLine(t *testing.T) {
	proxy := golang{}
	// True
	for _, line := range getPassingTags("stop") {
		t.Run(fmt.Sprintf("line='%s'", line), func(t *testing.T) {
			if !proxy.IsStopLine(line) {
				t.Errorf("Expected IsStopLine() to return `true` for '%s'", line)
			}
		})
	}
	// False
	for _, line := range getFailingTags("stop") {
		t.Run(fmt.Sprintf("line='%s'", line), func(t *testing.T) {
			if proxy.IsStopLine(line) {
				t.Errorf("Expected IsStopLine() to return `false` for '%s'", line)
			}
		})
	}
}

func TestShouldMerge(t *testing.T) {
	proxy := golang{}
	// True
	for _, path := range []string{
		"module_test.go",
		"module/submodule_test.go",
	} {
		t.Run(fmt.Sprintf("path='%s'", path), func(t *testing.T) {
			if !proxy.ShouldMerge(path, []byte{}) {
				t.Errorf("Expected ShouldMerge to return `true` for '%s'", path)
			}
		})
	}
	// False
	for _, path := range []string{
		"module.go",
		"module/submodule.go",
		"test.go",
		"readme.md",
	} {
		t.Run(fmt.Sprintf("path='%s'", path), func(t *testing.T) {
			if proxy.ShouldMerge(path, []byte{}) {
				t.Errorf("Expected ShouldMerge to return `false` for '%s'", path)
			}
		})
	}
}
