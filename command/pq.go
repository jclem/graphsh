package command

import (
	"fmt"

	"github.com/jclem/graphsh/types"
)

// Pq shows present node path's query
type Pq struct{}

func testPq(input string) (Command, error) {
	if input == "pq" {
		return &Pq{}, nil
	}

	return nil, nil
}

// Execute implements the Command interface
func (c Pq) Execute(s types.Session) error {
	fmt.Println(s.RootQuery())
	return nil
}
