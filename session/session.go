package session

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/textproto"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jclem/graphsh/command"
	"github.com/jclem/graphsh/graphql"
	"github.com/jclem/graphsh/introspection"
	"github.com/jclem/graphsh/querybuilder"
)

type (
	// Options represents a session's initializer options
	Options struct {
		Endpoint string
		Headers  []string
	}

	// Session represents a shell session
	Session struct {
		client       graphql.Querier
		endpoint     string
		headers      []string
		rootQuery    *querybuilder.Query
		currentQuery *querybuilder.Query
	}
)

// Client implements types.Session
func (s Session) Client() graphql.Querier {
	return s.client
}

// Endpoint implements types.Session
func (s Session) Endpoint() string {
	return s.endpoint
}

// Headers implements types.Session
func (s Session) Headers() []string {
	return s.headers
}

// RootQuery implements types.Session
func (s Session) RootQuery() *querybuilder.Query {
	return s.rootQuery
}

// CurrentQuery implements types.Session
func (s Session) CurrentQuery() *querybuilder.Query {
	return s.currentQuery
}

// SetCurrentQuery implements types.Session
func (s *Session) SetCurrentQuery(query *querybuilder.Query) {
	s.currentQuery = query
}

// NewSession creates a new session
func NewSession(options Options) (*Session, error) {
	headers, err := parseHeaders(options.Headers)
	if err != nil {
		return nil, err
	}

	client := graphql.New(options.Endpoint, headers)
	query := querybuilder.NewRootQuery()

	// Load the schema for this session
	if err := introspection.LoadSchema(client); err != nil {
		return nil, err
	}

	return &Session{
		client:       client,
		endpoint:     options.Endpoint,
		headers:      options.Headers,
		rootQuery:    query,
		currentQuery: query,
	}, nil
}

func parseHeaders(headers []string) (map[string][]string, error) {
	headerString := fmt.Sprintf("%s\r\n\r\n", strings.Join(headers, "\r\n"))
	reader := bufio.NewReader(strings.NewReader((headerString)))
	tp := textproto.NewReader(reader)
	return tp.ReadMIMEHeader()
}

// Loop starts a session loop to react to user input, using a default prompt
func Loop(options Options) {
	s, err := NewSession(options)

	isInterrupting := false

	reader, err := readline.New("› ")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	for {
		wasInterrupting := isInterrupting
		isInterrupting = false

		line, err := reader.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				if wasInterrupting {
					os.Exit(0)
				}

				isInterrupting = true
				fmt.Fprintln(os.Stderr, "To exit, press ^C again, or press ^D, or use the `exit` command.")
				continue
			}

			if err == io.EOF {
				os.Exit(0)
			}

			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if err := s.execInput(line); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (s *Session) execInput(line string) error {
	// Remove trailing input newline
	input := strings.TrimSuffix(line, "\n")

	cmd, err := command.FindCommand(input)

	if err != nil {
		return err
	}

	return cmd.Execute(s)
}
