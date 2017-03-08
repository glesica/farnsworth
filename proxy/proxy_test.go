package proxy

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

type testProxy struct{}

func (proxy *testProxy) IsHideLine(line string) bool {
	return strings.HasPrefix(line, "HIDE")
}

func (proxy *testProxy) IsStopLine(line string) bool {
	return strings.HasPrefix(line, "STOP")
}

func (proxy *testProxy) Name() string {
	return "test proxy"
}

func (proxy *testProxy) ShouldMerge(path string, content io.Reader) bool {
	return false
}

var validInputs = []string{
	`HIDE
hidden
STOP`,
	`content
HIDE
hidden
STOP
content`,
	`content
HIDE
hidden
STOP`,
	`HIDE
hidden
STOP
content`,
	`HIDE
hidden 0
hidden 1
STOP`,
}

var validOutputs = []string{
	``,
	`content
content`,
	`content`,
	`content`,
	``,
}

func TestValidRemoveHiddenLines(t *testing.T) {
	proxy := testProxy{}
	for i, input := range validInputs {
		inputBuffer := bytes.NewBufferString(input)
		output, err := RemoveHiddenLines(inputBuffer, &proxy)
		if err != nil {
			t.Errorf("expected '%s' to be valid", input)
		}
		validOutput := validOutputs[i]
		if output != validOutput {
			t.Errorf("expected '%s' -> '%s' but got '%s' instead", input, validOutput, output)
		}
	}
}

var invalidInputs = []string{
	`HIDE
HIDE
hidden
STOP`,
	`content
HIDE
HIDE
hidden
STOP
content`,
	`HIDE
hidden
STOP
STOP`,
	`content
HIDE
hidden
STOP
STOP
content`,
}

func TestInvalidRemoveHiddenLines(t *testing.T) {
	proxy := testProxy{}
	for _, input := range invalidInputs {
		inputBuffer := bytes.NewBufferString(input)
		output, err := RemoveHiddenLines(inputBuffer, &proxy)
		if err == nil {
			t.Errorf("expected '%s' -> `error` but got '%s' instead", input, output)
		}
	}
}
