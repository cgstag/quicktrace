package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	// Stdout is a color friendly pipe.
	Stdout = colorable.NewColorableStdout()

	// Stderr is a color friendly pipe.
	Stderr = colorable.NewColorableStderr()
)

const (
	defaultInput  = "samples/micro-log.txt"
	defaultOutput = "stdout"
	defaultStats  = true
)

// Options contain the command line options passed to the program.
type Options struct {
	Input  string
	Output string
	Help   bool
	Stats  bool
}

// ParseOptions parses the command line options.
func ParseOptions() *Options {
	var opt Options
	flag.StringVar(&opt.Input, "input", defaultInput, "Stdin or filename and format of the input file.")
	flag.StringVar(&opt.Output, "output", defaultOutput, "Output is file or stdout (default stdout)")
	flag.BoolVar(&opt.Stats, "stats", defaultStats, "Toggle statistics")
	flag.Parse()

	return &opt
}

// Valid checks command line options are valid.
func (opt *Options) Valid() bool {

	// File exists
	if _, err := os.Stat(opt.Input); os.IsNotExist(err) && opt.Input != "stdin" {
		fmt.Fprintln(Stderr, color.RedString("File does not exists: %s", opt.Input))
		return false
	}
	// Correct file format for input
	if opt.Input != "stdin" && !strings.Contains(opt.Input, ".") {
		fmt.Fprintln(Stderr, color.RedString("Missing format for input file: %s", opt.Input))
		return false
	}
	// Correct file format for output
	if opt.Output != "stdout" && !strings.Contains(opt.Output, ".") {
		fmt.Fprintln(Stderr, color.RedString("Missing format for output file: %s", opt.Output))
		return false
	}

	return true
}

// PrintUsage prints the usage of the program.
func (opt *Options) PrintUsage() {
	var banner = `QuickTrace`
	color.Cyan(banner)
	flag.Usage()
}
