# graphsh <img src="logo.png" width=50>

Graphsh (pronounced "graphsh") is an interactive shell for exploring GraphQL APIs.

## Installation

```console
$ go get github.com/jclem/graphsh
```

## Usage

To start a `graphsh` session, supply a GraphQL endpoint and optional headers.

```console
$ graphsh https://api.github.com/graphql -H "Authorization: Bearer $token"
› ▋
```

### Commands

#### Path traversal

Since Graphsh represents nested fields as something akin to directories, use the path traversal syntax to set your current path, relative to the current path.

```
› .repository(owner: "jclem", name: "graphsh")
› pp
.repository(owner: "jclem", name: "graphsh")
› .owner
› pp
.repository(owner: "jclem", name: "graphsh").owner
```

You can use `..` to traverse upwards:

```
› .repository(owner: "jclem", name: "graphsh").owner
› pp
.repository(owner: "jclem", name: "graphsh").owner
› ..
› pp
.repository(owner: "jclem", name: "graphsh")
```

#### `on`

The `on {ConcreteType}` command applies a concrete type to your current query. Passing no concrete type removes one.

```
› .repository(owner: "jclem", name: "graphsh").object(expression: "HEAD")
› pq
query {
  repository( name: "graphsh", owner: "jclem") {
    object(expression: "HEAD") {

    }
  }
}
› on Commit
› pq
query {
  repository( name: "graphsh", owner: "jclem") {
    object(expression: "HEAD") {
      ... on Commit {

      }
    }
  }
}
› on
› pq
query {
  repository( name: "graphsh", owner: "jclem") {
    object(expression: "HEAD") {

    }
  }
}
```

#### `pp`

The `pp` command shows your present path.

```
› .repository(owner: "jclem", name: "graphsh")
› pp
.repository(owner: "jclem", name: "graphsh")
› .owner
› pp
.repository(owner: "jclem", name: "graphsh").owner
```

#### `pq`

The `pq` command shows the current state of your query, which is a reflection of your present path.

```
› .repository(owner: "jclem", name: "graphsh").owner
› pq
query {
  repository(owner: "jclem", name: "graphsh") {
    owner {

    }
  }
}
```

#### `ls`

The `ls` command shows information about each field on the current node.

```
› .meta
› ls
gitHubServicesSha                  Returns a String that's a SHA of `github-services`
gitIpAddresses                     IP addresses that users connect to for git operations
hookIpAddresses                    IP addresses that service hooks are sent from
importerIpAddresses                IP addresses that the importer connects from
isPasswordAuthenticationVerifiable Whether or not users are verified
pagesIpAddresses                   IP addresses for GitHub Pages' A records
```

#### Querying

In order to query, use an expression surrounded by brackets.

```
› .repository(owner: "jclem", name: "graphsh")
› {name, owner {login}}
{
  "data": {
    "repository": {
      "name": "graphsh",
      "owner": {
        "login": "jclem"
      }
    }
  }
}
```
