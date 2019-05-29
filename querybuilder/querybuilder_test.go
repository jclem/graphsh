package querybuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQuery(t *testing.T) {
	query := NewQuery("repository", map[string]interface{}{"foo": "bar"})
	assert.Equal(t, "repository", query.Name, "name should be set")
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, query.Args, "args should be set")
}

func TestNewRootQuery(t *testing.T) {
	query := NewRootQuery()
	assert.Equal(t, "query", query.Name, "root query name should be set")
}

func TestToPath(t *testing.T) {
	var query *Query

	query = NewQuery("repository", map[string]interface{}{})
	assert.Equal(t, ".repository", query.Path())

	query = NewQuery("repository", map[string]interface{}{"owner": "jclem", "name": "graphsh"})
	assert.Equal(t, `.repository(name: "graphsh", owner: "jclem")`, query.Path())
}

func TestString(t *testing.T) {
	var query *Query

	query = NewQuery("repository", map[string]interface{}{})
	assert.Equal(t, `repository {

}`, query.String())

	query = NewQuery("repository", map[string]interface{}{"owner": "jclem", "name": "graphsh"})
	assert.Equal(t, `repository(name: "graphsh", owner: "jclem") {

}`, query.String())

	query = NewQuery("repository", map[string]interface{}{"owner": "jclem", "name": "graphsh"})
	child := NewQuery("owner", map[string]interface{}{})
	query.AddChild(child)
	assert.Equal(t, `repository(name: "graphsh", owner: "jclem") {
  owner {

  }
}`, query.String())

	query = NewQuery("repository", map[string]interface{}{})
	query.ConcreteType = "Foo"
	assert.Equal(t, `repository {
  ... on Foo {

  }
}`, query.String())

	query = NewQuery("repository", map[string]interface{}{})
	query.ConcreteType = "Foo"
	child = NewQuery("owner", map[string]interface{}{})
	query.AddChild(child)
	assert.Equal(t, `repository {
  ... on Foo {
    owner {

    }
  }
}`, query.String())
}
