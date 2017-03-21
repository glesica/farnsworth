package ignore

import (
	"regexp"
	"io"
	"bufio"
	"os"
	"path"
)

const IGNORE_FILE_NAME = ".farnsworthignore"

// A Filter determines whether or not a given path should be
// ignored.
type Filter interface {
	ShouldIgnore(filePath string) bool
}

// A predicate returns `true` for paths that should be excluded from
// an archive or merge operation, and `false` otherwise.
//
// For example, users might want to keep version control metadata out of
// archives, in which case a path that ended in `.gitignore` might cause
// the predicate to return `true`.
type predicate func(path string) bool

// A Filter is a collection of Predicates.
type filter struct {
	predicates []predicate
}

// ShouldIgnore returns the logical disjunction of the predicates
// included in the Filter.
func (f *filter) ShouldIgnore(filePath string) bool {
	for _, predicate := range f.predicates {
		if predicate(filePath) {
			return true
		}
	}
	return false
}

// addPredicate adds a predicate to the Filter.
func (f *filter) addPredicate(predicate predicate) {
	f.predicates = append(f.predicates, predicate)
}

// newRegexPredicate creates a predicate that checks to see if the path
// matches the given regular expression. The predicate will return
// `true` if they match.
func newRegexPredicate(patternString string) (predicate, error) {
	pattern, err := regexp.CompilePOSIX(patternString)
	if err != nil {
		return nil, err
	}

	return func(filePath string) bool {
		return pattern.MatchString(filePath)
	}, nil
}

func Load(ignoreFile io.Reader) (Filter, error) {
	f := filter{}

	ignoreScanner := bufio.NewScanner(ignoreFile)
	for ignoreScanner.Scan() {
		nextPatternString := ignoreScanner.Text()

		nextPredicate, err := newRegexPredicate(nextPatternString)
		if err != nil {
			return nil, err
		}

		f.addPredicate(nextPredicate)
	}

	return &f, nil
}

func Get(rootPath string) (Filter, error) {
	ignoreFile, err := os.Open(path.Join(rootPath, IGNORE_FILE_NAME))
	if err != nil {
		// If the file doesn't exist, that's not really an error,
		// we just return an "empty" filter.
		if os.IsNotExist(err) {
			return &filter{}, nil
		}
		return nil, err
	}

	return Load(ignoreFile)
}
