package completer

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jclem/graphsh/introspection"
	"github.com/jclem/graphsh/types"
)

// NewCompleter creates a new completer for the given session
func NewCompleter(s types.Session) *Completer {
	return &Completer{s}
}

// Completer provides tab-completion for a session
type Completer struct {
	session types.Session
}

var traversePattern = regexp.MustCompile(`^\.[a-zA-Z0-9_-]*`)

// Do implements the readline AutoCompleter interface
// See: https://github.com/chzyer/readline/blob/2972be24d48e78746da79ba8e24e8b488c9880de/complete.go#L11-L18
func (c Completer) Do(input []rune, offset int) ([][]rune, int) {
	strInput := string(input)

	if traversePattern.Match([]byte(strInput)) {
		return c.completeTraversal(strInput)
	}

	return nil, 0
}

func (c Completer) completeTraversal(input string) ([][]rune, int) {
	chars := strings.TrimLeft(input, ".")

	fields, err := introspection.GetFields(c.session.Client(), c.session.RootQuery())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, 0
	}

	var opts [][]rune

	for _, field := range fields {
		if strings.HasPrefix(field.Name, chars) {
			name := field.Name[len(input)-1:]
			opts = append(opts, []rune(name))
		}
	}

	return opts, len(input)
}
