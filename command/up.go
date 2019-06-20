package command

import (
	"regexp"
	"strings"

	"github.com/jclem/graphsh/types"
)

// Up traverses up one node
type Up struct {
	levels int
}

var upPattern = regexp.MustCompile(`^\.\.(?:\/\.\.)*$`)

func testUp(input string) (Command, error) {
	if upPattern.Match([]byte(input)) {
		length := len(strings.Split(input, "/")) - 1
		return &Up{length}, nil
	}

	return nil, nil
}

// Execute implements the Command interface
func (c Up) Execute(s types.Session) error {
	for i := 0; i <= c.levels; i++ {
		if s.CurrentQuery() == s.RootQuery() {
			return nil
		}

		s.SetCurrentQuery(s.CurrentQuery().Drop())
	}

	return nil
}
