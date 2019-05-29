package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jclem/graphsh/introspection"
	"github.com/jclem/graphsh/types"
)

// Ls lists fields for the current node
type Ls struct{}

func testLs(input string) (Command, error) {
	if input == "ls" {
		return &Ls{}, nil
	}

	return nil, nil
}

// Execute implements the Command interface
func (c Ls) Execute(s types.Session) error {
	fields, err := introspection.GetFields(s.Client(), s.RootQuery())
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s", "NAME", "TYPE", "DESCRIPTION"))

	for _, field := range fields {
		fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s", field.Name, field.GetHumanTypeName(), field.Description))
	}

	tw.Flush()

	return nil
}
