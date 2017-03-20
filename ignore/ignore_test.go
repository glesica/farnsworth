package ignore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"strings"
)

func addPredicate(t *testing.T, filter *Filter, predicate Predicate) {
	startLength := len(filter.predicates)
	filter.AddPredicate(predicate)
	assert.Len(t, filter.predicates, startLength + 1)
}

func TestAddPredicate(t *testing.T) {
	filter := Filter{}
	addPredicate(t, &filter, func(path string) bool {
		return true
	})
	addPredicate(t, &filter, func(path string) bool {
		return false
	})
	assert.True(t, filter.predicates[0](""))
	assert.False(t, filter.predicates[1](""))
}

func TestShouldIgnore(t *testing.T) {
	filter := Filter{}
	filter.AddPredicate(func(path string) bool {
		return strings.HasPrefix(path, "ignore")
	})
	filter.AddPredicate(func(path string) bool {
		return strings.HasSuffix(path, "ignore")
	})
	assert.True(t, filter.ShouldIgnore("ignore something"))
	assert.True(t, filter.ShouldIgnore("something ignore"))
	assert.True(t, filter.ShouldIgnore("ignore something ignore"))
	assert.False(t, filter.ShouldIgnore("something"))
}
