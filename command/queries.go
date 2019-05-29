package command

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jclem/graphsh/types"
)

type parsedSchema struct {
	Data struct {
		Type struct {
			Fields []struct {
				Name        string
				Description string
			}
		} `json:"__type"`
	}
}

func getSchema(s types.Session, typename string) (parsedSchema, error) {
	var schema parsedSchema

	schemaResult, err := s.Client().Query(fmt.Sprintf(`{
		__type(name: %q) {
			fields {
				name
				description
			}
		}
	}`, typename))

	if err != nil {
		return schema, err
	}

	if err := json.Unmarshal(schemaResult, &schema); err != nil {
		return schema, err
	}

	return schema, nil
}

func getTypename(s types.Session) (string, error) {
	// Execute the __typename query
	queryResp, err := executeQuery(s, "__typename")
	if err != nil {
		return "", err
	}

	var payload struct {
		Data map[string]interface{}
	}

	if err := json.Unmarshal(queryResp, &payload); err != nil {
		return "", err
	}

	obj := payload.Data
	node := s.RootQuery().Child()

	// Loop through the __typename response until we get to its root
	for {
		if node == nil {
			break
		}

		switch t := obj[node.Name].(type) {
		case map[string]interface{}:
			obj = t
		case []interface{}:
			elem, ok := t[0].(map[string]interface{})

			if !ok {
				return "", errors.New("Unexpected type")
			}

			obj = elem
		}

		node = node.Child()
	}

	typename, ok := obj["__typename"].(string)
	if !ok {
		return "", errors.New("No __typename key in query")
	}

	return typename, nil
}
