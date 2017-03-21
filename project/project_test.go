package project

import (
	"testing"

	"strings"

	"github.com/glesica/farnsworth/proxy"
	"github.com/glesica/farnsworth/ignore"
	"github.com/stretchr/testify/assert"
)

func filterFactory(path string) (ignore.Filter, error) {
	return nil, nil
}

func proxyFactory(path string) (proxy.Proxy, error) {
	return nil, nil
}

func TestLoad(t *testing.T) {
	path := "path/to/project"

	proj, err := Load(path, filterFactory, proxyFactory)
	if err != nil {
		t.Error("expected load to return without error")
	}

	assert.Equal(t, "project", proj.BaseName())

	if !strings.HasSuffix(proj.Path(), path) {
		t.Errorf("expected path to be '%s', found '%s'", path, proj.Path())
	}
}
