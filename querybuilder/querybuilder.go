package querybuilder

import (
	"fmt"
	"sort"
	"strings"
)

// Query represents a GraphQL query object
type Query struct {
	Name         string
	Args         map[string]interface{}
	ConcreteType string
	child        *Query
	parent       *Query
	isRoot       bool
}

type queryArgs = map[string]interface{}

// NewRootQuery creates a new root Query struct
func NewRootQuery() *Query {
	return &Query{Name: "query", isRoot: true}
}

// NewQuery creates a new Query struct
func NewQuery(name string, args map[string]interface{}) *Query {
	return &Query{Name: name, Args: args}
}

// AddChild adds a new child Query to another Query
func (q *Query) AddChild(child *Query) {
	q.child = child
	child.parent = q
}

// Drop removes this query from its parent and returns the parent
func (q *Query) Drop() *Query {
	p := q.parent
	q.parent.child = nil
	q.parent = nil
	return p
}

// Child returns the query's child
func (q Query) Child() *Query {
	return q.child
}

// Parent returns the query's parent
func (q Query) Parent() *Query {
	return q.parent
}

// Path converts a query to a human-readable path
func (q *Query) Path() string {
	var p strings.Builder

	for {
		if q == nil {
			return p.String()
		}

		p.WriteString(fmt.Sprintf(".%s", q.Name))

		var queryArgs strings.Builder
		argsToString(q.Args, &queryArgs)
		p.WriteString(queryArgs.String())

		q = q.child
	}
}

// List converts a query to a list, excluding the root node
func (q *Query) List() []*Query {
	var list []*Query

	node := q

	for {
		if node == nil {
			break
		}

		if node.isRoot {
			node = node.Child()
			continue
		}

		list = append(list, node)
		node = node.Child()
	}

	return list
}

func (q *Query) String() string {
	return q.ToString("", "")
}

// ToString converts a query to a string with the given indentation
func (q *Query) ToString(tailQuery string, indent string) string {
	if q == nil {
		return ""
	}

	var query strings.Builder
	var queryArgs strings.Builder

	query.WriteString(fmt.Sprintf("%s%s", indent, q.Name))

	argsToString(q.Args, &queryArgs)

	query.WriteString(fmt.Sprintf("%s {", queryArgs.String()))

	if q.ConcreteType != "" {
		query.WriteRune('\n')
		query.WriteString(fmt.Sprintf("%s  ... on %s {", indent, q.ConcreteType))
	}

	if q.child == nil {
		query.WriteRune('\n')
		query.WriteString(tailQuery)
	} else {
		childIndent := indent
		if q.ConcreteType != "" {
			childIndent = fmt.Sprintf("  %s", indent)
		}

		query.WriteRune('\n')
		query.WriteString(q.child.ToString(tailQuery, fmt.Sprintf("%s  ", childIndent)))
	}

	if q.ConcreteType != "" {
		query.WriteString(fmt.Sprintf("\n%s  }", indent))
	}

	query.WriteString(fmt.Sprintf("\n%s}", indent))

	return query.String()
}

// WithQuery stringifies the full query with the given query in lowest child
func (q Query) WithQuery(tailQuery string) string {
	return q.ToString(tailQuery, "")
}

func argsToString(m map[string]interface{}, b *strings.Builder) {
	eachSortedKey(m, func(k string, v interface{}) {
		if b.Len() > 0 {
			b.WriteString(", ")
		} else {
			b.WriteRune('(')
		}

		b.WriteString(fmt.Sprintf("%s: ", k))

		switch t := v.(type) {
		case bool:
			b.WriteString(fmt.Sprintf("%t", t))
		case int:
			b.WriteString(fmt.Sprintf("%d", t))
		case string:
			b.WriteString(fmt.Sprintf("%q", t))
		default:
			panic(fmt.Sprintf("Unrecognized type in query args %s", t))
		}
	})

	if b.Len() > 0 {
		b.WriteRune(')')
	}
}

func eachSortedKey(m map[string]interface{}, fn func(key string, value interface{})) {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := m[k]
		fn(k, v)
	}
}
