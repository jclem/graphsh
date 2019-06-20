package command

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/jclem/graphsh/types"
)

// Help prints out explanations for various commands
type Help struct {
	command string
}

type helpInfo struct {
	usage       string
	description string
}

var helpMap = map[string]helpInfo{
	"exit": {
		usage:       "exit",
		description: "Exits the graphsh shell",
	},
	"help": {
		usage:       "help | help <command>",
		description: "Displays help for a command",
	},
	"ls": {
		usage:       "ls",
		description: "Lists the fields for the current query node",
	},
	"on": {
		usage:       "on <ConcreteType>",
		description: "Applies a concrete type to the current query node",
	},
	"pp": {
		usage:       "pp",
		description: "Prints the current query path",
	},
	"pq": {
		usage:       "pq",
		description: "Prints the current query",
	},
	".": {
		usage: ".<field>[...]",
		description: `Traverses through fields of the current query

For example, ".foo.bar(first: 10).baz"`,
	},
	"..": {
		usage:       "..[/..]",
		description: "Traverses upwards one or more query nodes",
	},
	"{}": {
		usage: "{<field>}",
		description: `Execute a query in the current query node

For example, to query the current node's URL and its app name: { url, app { name } }`,
	},
}

var helpPattern = regexp.MustCompile(`^h(?:elp)?( .+)?$`)

func testHelp(input string) (Command, error) {
	match := helpPattern.FindStringSubmatch(input)

	if len(match) == 0 {
		return nil, nil
	}

	return &Help{strings.TrimPrefix(match[1], " ")}, nil
}

// Execute implements the Command interface
func (c Help) Execute(s types.Session) error {
	if c.command == "" {
		helpKeys := make([]string, 0, len(helpMap))
		for key := range helpMap {
			helpKeys = append(helpKeys, key)
		}
		sort.Strings(helpKeys)

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		fmt.Fprintln(tw, fmt.Sprintf("%s\t%s", "COMMAND", "DESCRIPTION"))

		for _, key := range helpKeys {
			helpInfo := helpMap[key]
			description := strings.Split(helpInfo.description, "\n")[0]
			fmt.Fprintln(tw, fmt.Sprintf("%s\t%s", key, description))
		}

		tw.Flush()

		return nil
	}

	helpInfo, ok := helpMap[c.command]

	if !ok {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("No such command %q exists", c.command))
		return nil
	}

	fmt.Printf("Usage: %s\n\n%s\n", helpInfo.usage, helpInfo.description)

	return nil
}
