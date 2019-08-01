// @author Flavio Deroo https://github.com/cgstag 08/2019
// This program has been developed for SimScale GmbH as an job interview technical assessment

// Program provides functionality for parsing and formatting log files into JSON trees
// It is meant to be used by a CLI with flag parameters

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"quicktrace/cmd/quicktrace"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var input = flag.String("input", "file", "Input is file or stdin ? (default stdin)")
var filenameInput = flag.String("filenameInput", "samples/micro-log.txt", "filename and format")
var output = flag.String("output", "stdout", "Output is file or stdout (default stdout)")
var filenameOutput = flag.String("filenameOutput", "quick-traces.txt", "filename and format")

func main() {

	ctx := &quicktrace.Context{
		Start: time.Now(),
	}

	// Parse Flags
	flag.Parse()
	inputFlag := string(*input)
	filenameInputFlag := string(*filenameInput)
	outputFlag := string(*output)
	filenameOutputFlag := string(*filenameOutput)

	// Read Input
	chrono("Reading Input", ctx.Start)
	scanner := new(bufio.Scanner)
	if inputFlag == "file" {
		if filenameInputFlag == "" {
			os.Stderr.WriteString("You have chosen file as an input but no filename was provided. --filenameInput missing. go mod help for help")
		} else {
			file, err := os.Open(filenameInputFlag)
			check(err)
			defer file.Close()
			scanner = bufio.NewScanner(file)
		}
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	// Bufferize Scanner
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	// Filter all entries into an Hashmap
	traces := quicktrace.TraceMap{}
	chrono("Parsing Logs Entries... This can take a while", ctx.Start)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			break
		} else {
			traces.PushSpan(scanner.Text())
		}

	}
	// Progress Bar
	chrono("Building JSON Traces", ctx.Start)
	bar := pb.StartNew(len(traces))

	// Create Output
	w := new(bufio.Writer)
	f := new(os.File)
	if outputFlag == "file" {
		if filenameOutputFlag == "" {
			os.Stderr.WriteString("You have chosen file as an output but no filename was provided. --filenameOutput missing. go mod help for help")
		} else {
			output, err := os.Create(filenameOutputFlag)
			check(err)
			defer output.Close()
			f = output
		}
	} else {
		f = os.Stdout
		defer w.Flush()
	}
	// Prepare writer
	w = bufio.NewWriter(f)
	written := int64(0)
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}
	wg.Add(len(traces))
	// Organize each trace in tree
	for _, trace := range traces {
		bar.Increment()
		go func(trace quicktrace.Trace) {
			_, err := trace.QuickTrace(ctx)
			if err != nil {
				// If any error occur, we delete the trace from final result - Example Orphans
				delete(traces, trace.Id)
			} else {
				strTraces, err := json.Marshal(trace)
				check(err)
				mutex.Lock()
				nn, err := w.WriteString(string(strTraces) + "\n")
				mutex.Unlock()
				written += int64(nn)
				check(err)
			}
			wg.Done()
		}(trace)
	}
	wg.Wait()
	err := w.Flush()
	check(err)
	bar.Finish()
	chrono("Program terminated", ctx.Start)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func chrono(text string, t time.Time) {
	log.Printf("[%s] %s ", time.Since(t), text)
}
