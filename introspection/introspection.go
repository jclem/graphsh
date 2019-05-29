package introspection

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jclem/graphsh/graphql"
	"github.com/jclem/graphsh/querybuilder"
)

var schema *Schema

// GetFields gets the fields for a given query
func GetFields(q graphql.Querier, query *querybuilder.Query) ([]Field, error) {
	typ, ok := schema.GetQueryType()
	if !ok {
		return nil, errors.New("No QueryType present in schema")
	}

	for _, node := range query.List() {
		if node.ConcreteType == "" {
			field, ok := typ.GetField(node.Name)
			if !ok {
				return nil, fmt.Errorf("Missing field %q from type %q", node.Name, typ.Name)
			}

			typ, ok = schema.GetType(field.GetTypeName())
			if !ok {
				return nil, fmt.Errorf("Missing type %q", field.GetTypeName())
			}
		} else {
			typ, ok = schema.GetType(node.ConcreteType)
			if !ok {
				return nil, fmt.Errorf("Missing type %q", node.ConcreteType)
			}
		}
	}

	return typ.Fields, nil
}

// LoadSchema pre-loads the schema struct
func LoadSchema(q graphql.Querier) error {
	if schema != nil {
		return nil
	}

	respBody, err := q.Query(schemaQuery)
	if err != nil {
		return err
	}

	var introspection introspection
	if err := json.Unmarshal(respBody, &introspection); err != nil {
		return err
	}

	schema = &introspection.Data.Schema

	return nil
}
