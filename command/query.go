package command

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/jclem/graphsh/types"
)

var headerPattern = regexp.MustCompile("^(.+): (.+)$")
var queryPattern = regexp.MustCompile("^{.+}$")

// Query executes a GraphQL query
type Query struct {
	query string
}

func testQuery(input string) (Command, error) {
	if queryPattern.Match([]byte(input)) {
		cmd := &Query{query: input[1 : len(input)-1]}
		return cmd, nil
	}

	return nil, nil
}

// Execute implements the Command interface
func (c Query) Execute(s types.Session) error {
	body, err := executeQuery(s, c.query)
	if err != nil {
		return err
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return err
	}

	json, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(json))

	return nil
}

func executeQuery(s types.Session, query string) ([]byte, error) {
	fullQuery := s.RootQuery().WithQuery(query)
	return s.Client().Query(fullQuery)
}
