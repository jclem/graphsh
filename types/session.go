package types

import (
	"github.com/jclem/graphsh/graphql"
	"github.com/jclem/graphsh/querybuilder"
)

// Session gets session information
type Session interface {
	Client() graphql.Querier
	Endpoint() string
	Headers() []string
	RootQuery() *querybuilder.Query
	CurrentQuery() *querybuilder.Query
	SetCurrentQuery(q *querybuilder.Query)
}
