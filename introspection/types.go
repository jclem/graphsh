package introspection

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

var schemaQuery = `
query IntrospectionQuery {
  __schema {
    queryType {
      name
    }
    mutationType {
      name
    }
    subscriptionType {
      name
    }
    types {
      ...FullType
    }
    directives {
      name
      description
      locations
      args {
        ...InputValue
      }
    }
  }
}

fragment FullType on __Type {
  kind
  name
  description
  fields(includeDeprecated: true) {
    name
    description
    args {
      ...InputValue
    }
    type {
      ...TypeRef
    }
    isDeprecated
    deprecationReason
  }
  inputFields {
    ...InputValue
  }
  interfaces {
    ...TypeRef
  }
  enumValues(includeDeprecated: true) {
    name
    description
    isDeprecated
    deprecationReason
  }
  possibleTypes {
    ...TypeRef
  }
}

fragment InputValue on __InputValue {
  name
  description
  type {
    ...TypeRef
  }
  defaultValue
}

fragment TypeRef on __Type {
  kind
  name
  ofType {
    kind
    name
    ofType {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
									ofType {
										kind
										name
									}
                }
              }
            }
          }
        }
      }
    }
  }
}
`

type introspection struct {
	Data struct {
		Schema Schema `json:"__schema"`
	}
}

// Schema is a GraphQL schema, the result of an introspection query
type Schema struct {
	QueryType struct {
		Name string
	}

	MutationType *struct {
		Name string
	}

	SubscriptionType *struct {
		Name string
	}

	Types []FullType

	Directives []struct {
		Name        string
		Description string
		Locations   []string
		Args        []inputValue
	}
}

// GetQueryType returns the full query type
func (s Schema) GetQueryType() (*FullType, bool) {
	return s.GetType(s.QueryType.Name)
}

// GetType returns a type with the given name if one exists
func (s Schema) GetType(name string) (*FullType, bool) {
	for _, t := range s.Types {
		if t.Name == name {
			return &t, true
		}
	}

	return nil, false
}

// FullType is a full type description
type FullType struct {
	Kind        string
	Name        string
	Description string
	Fields      []Field
	InputFields []inputValue
	Interfaces  []typeRef
	EnumValues  []struct {
		Name              string
		Description       string
		IsDeprecated      bool
		DeprecationReason string
	}
	PossibleTypes []typeRef
}

// GetField gets a field with the given name, if it exists
func (t FullType) GetField(name string) (*Field, bool) {
	for _, f := range t.Fields {
		if f.Name == name {
			return &f, true
		}
	}

	return nil, false
}

// Field represents a field of a GraphQL type
type Field struct {
	Name              string
	Description       string
	Args              []inputValue
	Type              typeRef
	IsDeprecated      bool
	DeprecationReason string
}

// GetTypeName gets the name of the field's type
func (f Field) GetTypeName() string {
	t := f.Type

	for {
		if t.OfType == nil {
			return t.Name
		}

		t = *t.OfType
	}
}

var emptyKinds = []string{"INTERFACE", "NON_NULL", "OBJECT"}

// GetHumanTypeName is a human-readable type name
func (f Field) GetHumanTypeName() string {
	t := &f.Type

	name := f.Type.Name

	wrap := func(s string) string {
		return s
	}

	for {
		wrapFn := wrap
		kind := t.Kind

		if isEmptyKind(kind) {
			wrap = func(s string) string {
				return wrapFn(s)
			}
		} else if kind == "SCALAR" {
			wrap = func(s string) string {
				return fmt.Sprintf("{%s}", wrapFn(s))
			}
		} else if kind == "LIST" {
			wrap = func(s string) string {
				return fmt.Sprintf("[]%s", wrapFn(s))
			}
		} else {
			kind = strcase.ToCamel(strings.ToLower(kind))
			wrap = func(s string) string {
				return fmt.Sprintf("%s<%s>", kind, wrapFn(s))
			}
		}

		name = t.Name

		if t.OfType == nil {
			return wrap(name)
		}

		t = t.OfType
	}
}

func isEmptyKind(kind string) bool {
	for _, emptyKind := range emptyKinds {
		if emptyKind == kind {
			return true
		}
	}

	return false
}

type inputValue struct {
	Name         string
	Description  string
	Type         typeRef
	DefaultValue string
}

type typeRef struct {
	Kind   string
	Name   string
	OfType *typeRef
}

// 	Kind   string
// 	Name   string
// 	OfType *struct {
// 		Kind   string
// 		Name   string
// 		OfType *struct {
// 			Kind   string
// 			Name   string
// 			OfType *struct {
// 				Kind   string
// 				Name   string
// 				OfType *struct {
// 					Kind   string
// 					Name   string
// 					OfType *struct {
// 						Kind   string
// 						Name   string
// 						OfType *struct {
// 							Kind   string
// 							Name   string
// 							OfType *struct {
// 								Kind   string
// 								Name   string
// 								OfType *struct {
// 									Kind string
// 									Name string
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }
// }
