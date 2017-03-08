package project

import (
	"testing"

	"strings"

	"github.com/glesica/farnsworth/proxy"
)

func TestLoadWithInvalidPath(t *testing.T) {
	_, err := Load("/invalid/path/to/project")
	if err == nil {
		t.Error("expected Load to fail on invalid path")
	}
}

func TestLoadWithFactory(t *testing.T) {
	path := "path/to/project"
	proj, err := loadWithFactory(path, func(path string) (proxy.Proxy, error) {
		return nil, nil
	})

	if err != nil {
		t.Error("expected loadWithFactory to return without error")
	}

	if proj.BaseName() != "project" {
		t.Errorf("expected basename to be 'path', found '%s'", proj.BaseName())
	}

	if !strings.HasSuffix(proj.Path(), path) {
		t.Errorf("expected path to be '%s', found '%s'", path, proj.Path())
	}
}
