package command

import (
	"os"

	"github.com/jclem/graphsh/types"
)

// Exit exits the shell
type Exit struct{}

func testExit(input string) (Command, error) {
	if input == "exit" {
		return &Exit{}, nil
	}

	return nil, nil
}

// Execute implements the Command interface
func (c Exit) Execute(s types.Session) error {
	os.Exit(0)
	return nil
}
