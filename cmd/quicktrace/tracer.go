package quicktrace

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"quicktrace/cli"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type Tracer struct {
	TraceMap
	Options          *cli.Options
	Context          *cli.Context
	Orphans          []*Trace
	Malformed        int
	AverageTraceSize int
	AverageDepth     int
}

type TraceMap map[string]Trace

type Trace struct {
	Id       string  `json:"id"`
	Root     *Span   `json:"root"`
	Unsorted []*Span `json:"-"`
}

// ParseAndTrace parse an input log and build a trace
func (tr Tracer) ParseAndTrace() error {
	tr.PrintProgress("Setting up Scanner")
	scanner := new(bufio.Scanner)
	if tr.Options.Input == "stdin" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(tr.Options.Input)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)
	tr.TraceMap = TraceMap{}
	tr.Malformed = 0
	tr.AverageTraceSize = 0
	tr.PrintProgress("Parsing input... this could take a while")
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			break
		} else {
			tr.pushSpan(scanner.Text())
		}
	}
	if err := tr.BuildTrace(); err != nil {
		return err
	}
	return nil
}

// NewTracer creates a new Tracer
func NewTracer(opt *cli.Options, ctx *cli.Context) *Tracer {
	var t = &Tracer{
		TraceMap: TraceMap{},
		Options:  opt,
		Context:  ctx,
	}

	return t
}

func (tr Tracer) BuildTrace() (err error) {
	bar := pb.StartNew(len(tr.TraceMap))
	tr.PrintProgress("Parsing input... this could take a while")
	f := new(os.File)
	if tr.Options.Output != "stdout" {
		output, err := os.Create(tr.Options.Output)
		if err != nil {
			return err
		}
		defer output.Close()
		f = output
	} else {
		f = os.Stdout
	}

	w := bufio.NewWriter(f)
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	wg.Add(len(tr.TraceMap))

	for _, trace := range tr.TraceMap {
		bar.Increment()
		go func(trace Trace) {
			if tr.Options.Stats {
				tr.AverageTraceSize += len(trace.Unsorted) + 1
			}
			trace.Root, trace.Unsorted, err = matchSpan(trace.Root, trace.Unsorted)
			if err != nil && err.Error() == "NoRootSpan" {
				fmt.Fprintf(os.Stderr, "Orphan line detected for Trace: %s \n", trace.Id)
				tr.Orphans = append(tr.Orphans, &trace)
			}
			if err != nil {
				// If any error occur, we delete the trace from final result - Example Orphans
				delete(tr.TraceMap, trace.Id)
			} else {
				strTraces, _ := json.Marshal(trace)
				mutex.Lock()
				nn, _ := w.WriteString(string(strTraces) + "\n")
				mutex.Unlock()
				tr.Context.IO += int64(nn)
			}
			wg.Done()
		}(trace)
	}
	wg.Wait()
	_ = w.Flush()
	bar.Finish()
	tr.PrintProgress("Program terminated")
	if tr.Options.Stats {
		if tr.AverageTraceSize > 0 {
			tr.AverageTraceSize /= len(tr.TraceMap)
		}

		PrintStats(tr.Context, tr.Orphans, tr.Malformed, tr.AverageTraceSize)
	}
	return nil
}

func PrintStats(ctx *cli.Context, orphans []*Trace, empty int, averageSize int) {
	ctx.Elapsed = time.Since(ctx.Start)
	elapsedStats := fmt.Sprintf("Sucessfully executed QuickTrace in %s\n", ctx.Elapsed)
	IOStats := fmt.Sprintf("I/O %.2fMB/s\n", (float64(ctx.IO)/1000000)/(float64(ctx.Elapsed)/float64(time.Second)))
	orphanStats := fmt.Sprintf("%d Orphan traces\n", len(orphans))
	emptyStats := fmt.Sprintf("%d Empty Span\n", empty)
	averageSizeStats := fmt.Sprintf("Average trace size %d spans\n", averageSize)
	os.Stderr.WriteString(elapsedStats)
	os.Stderr.WriteString(IOStats)
	os.Stderr.WriteString(orphanStats)
	os.Stderr.WriteString(emptyStats)
	os.Stderr.WriteString(averageSizeStats)
}

func (tr *Tracer) PrintProgress(msg string) {
	if tr.Options.Stats {
		stderr := fmt.Sprintf("[%s] %s \n", time.Since(tr.Context.Start), msg)
		os.Stderr.WriteString(stderr)
	}
}
