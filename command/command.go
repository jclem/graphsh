package command

import (
	"fmt"

	"github.com/jclem/graphsh/types"
)

// Command is an interface for all commands to follow
type Command interface {
	Execute(s types.Session) error
}

// A list of tests per-command that determines if input matches a command
// Define one per command file with the name `test$CommandName`.
var tests = []func(input string) (Command, error){testExit, testHelp, testLs, testOn, testPp, testPq, testUp, testTraverse, testQuery}

// FindCommand finds a command for a given input
func FindCommand(input string) (Command, error) {
	for _, test := range tests {
		if cmd, err := test(input); err != nil {
			return nil, err
		} else if cmd != nil {
			return cmd, nil
		}
	}

	return nil, fmt.Errorf("No command for input %q", input)
}
