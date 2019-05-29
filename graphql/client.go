package graphql

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Querier is an interface that makes GraphQL requests
type Querier interface {
	Query(query string) ([]byte, error)
}

// New creates a new Querier client
func New(endpoint string, header http.Header) Querier {
	return client{
		client:   http.Client{},
		endpoint: endpoint,
		header:   header,
	}
}

type client struct {
	client   http.Client
	endpoint string
	header   http.Header
}

type newOptions struct {
	Endpoint string
	Header   http.Header
}

func (c client) Query(query string) ([]byte, error) {
	reqBody := []byte(fmt.Sprintf(`{"query":%q}`, query))
	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header = c.header

	if req.Header.Get("content-type") == "" {
		req.Header.Set("content-type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
