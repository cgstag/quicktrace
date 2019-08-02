// @author Flavio Deroo https://github.com/cgstag 08/2019
// This program has been developed for SimScale GmbH as an job interview technical assessment

// Program provides functionality for parsing and formatting log files into JSON trees
// It is meant to be used by a CLI with flag parameters

package main

import (
	"fmt"
	"quicktrace/cli"
	"quicktrace/cmd/quicktrace"

	"github.com/fatih/color"
)

func main() {
	options := cli.ParseOptions()
	context := cli.NewContext()

	if options.Help {
		options.PrintUsage()
	} else if options.Valid() {
		err := quicktrace.NewTracer(options, context).ParseAndTrace()

		if err != nil {
			fmt.Fprintln(cli.Stderr, color.RedString(err.Error()))
		}
	}
}
