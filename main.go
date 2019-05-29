package main

import (
	"fmt"
	"os"

	"github.com/jclem/graphsh/session"
	flag "github.com/spf13/pflag"
)

var headers = flag.StringArrayP("header", "H", []string{}, "Set a custom request header")
var help = flag.BoolP("help", "h", false, "Print this help message")

func main() {
	flag.Parse()

	flag.Usage = usage

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	endpoint := flag.Arg(0)
	if endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}

	session.Loop(session.Options{
		Endpoint: endpoint,
		Headers:  *headers,
	})
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: graphsh <endpoint> [<options>]")
	flag.PrintDefaults()
}
