package command

import (
	"fmt"

	"github.com/jclem/graphsh/types"
)

// Pp shows present node path
type Pp struct{}

func testPp(input string) (Command, error) {
	if input == "pp" {
		return &Pp{}, nil
	}

	return nil, nil
}

// Execute implements the Command interface
func (c Pp) Execute(s types.Session) error {
	fmt.Println(s.RootQuery().Path())
	return nil
}
