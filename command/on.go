package command

import (
	"regexp"

	"github.com/jclem/graphsh/types"
)

// On scopes the query with an inline fragment for a concrete type
type On struct {
	concreteType string
}

var onTest = regexp.MustCompile("^on(?: ([a-zA-Z0-9_-]+))?$")

func testOn(input string) (Command, error) {
	match := onTest.FindStringSubmatch(input)

	if len(match) == 0 {
		return nil, nil
	}

	return &On{match[1]}, nil
}

// Execute implements the Command interface
func (o On) Execute(s types.Session) error {
	s.CurrentQuery().ConcreteType = o.concreteType
	return nil
}
