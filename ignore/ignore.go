package ignore

import "github.com/gobwas/glob"

// A Predicate returns `true` for paths that should be excluded from
// an archive or merge operation, and `false` otherwise.
//
// For example, users might want to keep version control metadata out of
// archives, in which case a path that ended in `.gitignore` might cause
// the predicate to return `true`.
type Predicate func(path string) bool

// TODO: Make Filter into an interface.

// A Filter is a collection of Predicates.
type Filter struct {
	predicates []Predicate
}

// AddPredicate adds a Predicate to the Filter.
func (filter *Filter) AddPredicate(predicate Predicate) {
	filter.predicates = append(filter.predicates, predicate)
}

// ShouldIgnore returns the logical disjunction of the predicates
// included in the Filter.
func (filter *Filter) ShouldIgnore(path string) bool {
	for _, predicate := range filter.predicates {
		if predicate(path) {
			return true
		}
	}
	return false
}

func GlobPredicate(globString string) Predicate {
	g := glob.MustCompile(globString)
	return func(path string) bool {
		return g.Match(path)
	}
}

func Get(rootPath string) (Filter, error) {
	filter := Filter{}
	// Something like this...
	// There's a Go package that will read a .gitignore file, which
	// would be a handy feature, but I'd rather not couple to Git.
	// Maybe just use regular expressions? But how weird are Go
	// regular expressions?
	filter.AddPredicate(GlobPredicate("*.git"))
	return filter, nil
}
