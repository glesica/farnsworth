package ignore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"bytes"
	"strings"
)

func TestShouldIgnore(t *testing.T) {
	f := filter{}

	f.addPredicate(func(path string) bool {
		return strings.HasPrefix(path, "ignore")
	})
	f.addPredicate(func(path string) bool {
		return strings.HasSuffix(path, "ignore")
	})

	assert.True(t, f.ShouldIgnore("ignore something"))
	assert.True(t, f.ShouldIgnore("something ignore"))
	assert.True(t, f.ShouldIgnore("ignore something ignore"))
	assert.False(t, f.ShouldIgnore("something"))
}

func getNewRegexPredicate(t *testing.T, pattern string) predicate {
	predicate, err := newRegexPredicate(pattern)
	if err != nil {
		t.Fatalf("Failed to create predicate for '%s'", pattern)
	}
	return predicate
}

func TestNewRegexPredicate(t *testing.T) {
	p0 := getNewRegexPredicate(t,"simple")
	assert.True(t, p0("simple is good"))
	assert.False(t, p0("complex is bad"))

	p1 := getNewRegexPredicate(t,"^simple")
	assert.True(t, p1("simple is good"))
	assert.False(t, p1("good is simple"))

	p2 := getNewRegexPredicate(t, "[sd]imple")
	assert.True(t, p2("simple"))
	assert.True(t, p2("dimple"))
	assert.False(t, p2("complex"))
}

func TestLoad(t *testing.T) {
	buffer := bytes.NewBufferString("first\nsecond")

	f, err := load(buffer)
	if err != nil {
		t.Fatal("Expected load to complete successfully")
	}

	assert.True(t, f.ShouldIgnore("first"))
	assert.True(t, f.ShouldIgnore("second"))
	assert.True(t, f.ShouldIgnore("foo/first/bar"))
	assert.True(t, f.ShouldIgnore("foo/second/bar"))
	assert.False(t, f.ShouldIgnore("fourth"))
}
