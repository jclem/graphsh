package command

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/jclem/graphsh/querybuilder"
	"github.com/jclem/graphsh/types"
)

// Traverse traverses nodes
type Traverse struct {
	head *querybuilder.Query
	tail *querybuilder.Query
}

func testTraverse(input string) (Command, error) {
	if strings.HasPrefix(input, ".") {
		var cmd Traverse

		pathSegments := strings.Split(input, ".")[1:]

		for _, pathSegment := range pathSegments {
			if pathSegment == "" {
				return nil, errors.New("Path segment must not be empty")
			}

			query, err := parseTraversalSegment(pathSegment)

			if err != nil {
				return nil, err
			}

			if cmd.head == nil {
				cmd.head = query
			} else {
				cmd.tail.AddChild(query)
			}

			cmd.tail = query
		}

		return &cmd, nil
	}

	return nil, nil
}

// Execute implements the Command interface
func (c Traverse) Execute(s types.Session) error {
	s.CurrentQuery().AddChild(c.head)
	s.SetCurrentQuery(c.tail)
	return nil
}

var noArgPattern = regexp.MustCompile("^([_A-Za-z][_0-9A-Za-z]*)$")
var argPattern = regexp.MustCompile("^([_A-Za-z][_0-9A-Za-z]*)\\((.+)\\)$")

// TODO: This is really rough—currently parsing a float breaks because of an earlier split on "."
func parseTraversalSegment(segment string) (*querybuilder.Query, error) {
	noArgMatch := noArgPattern.FindStringSubmatch(segment)

	if noArgMatch != nil {
		return querybuilder.NewQuery(noArgMatch[0], map[string]interface{}{}), nil
	}

	argMatch := argPattern.FindStringSubmatch(segment)

	if argMatch != nil {
		argMap := map[string]interface{}{}

		pairs := strings.Split(argMatch[2], ",")

		for _, pair := range pairs {
			parts := strings.SplitN(pair, ":", 2)
			k := parts[0]
			v := parts[1]

			// Try int, since float64 is default
			var i int
			if err := json.Unmarshal([]byte(v), &i); err == nil {
				argMap[k] = i
			} else {
				var i interface{}
				if err := json.Unmarshal([]byte(v), &i); err != nil {
					return nil, err
				}
				argMap[k] = i
			}
		}

		return querybuilder.NewQuery(argMatch[1], argMap), nil
	}

	return nil, errors.New("Invalid traversal segment")
}
