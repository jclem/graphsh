package command

import (
	"github.com/jclem/graphsh/types"
)

// Up traverses up one node
type Up struct{}

func testUp(input string) (Command, error) {
	if input == ".." {
		return &Up{}, nil
	}

	return nil, nil
}

// Execute implements the Command interface
func (c Up) Execute(s types.Session) error {
	if s.CurrentQuery() == s.RootQuery() {
		return nil
	}

	s.SetCurrentQuery(s.CurrentQuery().Drop())
	return nil
}
