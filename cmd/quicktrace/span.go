package quicktrace

import (
	"strings"
)

type Span struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	Trace   string `json:"-"`
	Service string `json:"service"`
	Span    string `json:"span"`
	Caller  string `json:"-"`
	Calls   []Span `json:"calls"`
}

var current string

func matchSpan(root *Span, unsorted []*Span) (*Span, []*Span, error) {
	if root == nil {
		return root, unsorted, NewNoRootSpanError()
	}
	for index, currentSpan := range unsorted {
		if currentSpan.Caller == root.Span {
			unsorted = remove(unsorted, index)
			newSpan, newUnsorted, _ := matchSpan(currentSpan, unsorted)
			unsorted = newUnsorted
			root.Calls = append(root.Calls, *newSpan)
		}
	}
	return root, unsorted, nil
}

/**
* The challenge here is parse the log entry and assign it to the TraceMap
 */
func (tr Tracer) pushSpan(strSpan string) {
	span := toSpan(strSpan)
	if span.Span == "" {
		tr.Malformed++
		return
	}
	// Set timestamp of treatment as current
	current = span.End
	if _, ok := tr.TraceMap[span.Trace]; !ok {
		// This trace is new
		tr.TraceMap[span.Trace] = Trace{
			Id: span.Trace,
		}
		var current = tr.TraceMap[span.Trace]
		if span.Caller == "null" {
			// The span is root
			current.Root = span
			tr.TraceMap[span.Trace] = current
		} else {
			// Add to Unsorted list of that trace
			current.Unsorted = append(tr.TraceMap[span.Trace].Unsorted, span)
			tr.TraceMap[span.Trace] = current
		}
	} else {
		// This trace already exists
		var onGoingTrace = tr.TraceMap[span.Trace]
		if span.Caller == "null" {
			// The span is root
			onGoingTrace.Root = span
			tr.TraceMap[span.Trace] = onGoingTrace
		} else {
			// Add to Unsorted list of that trace
			onGoingTrace.Unsorted = append(tr.TraceMap[span.Trace].Unsorted, span)
			tr.TraceMap[span.Trace] = onGoingTrace
		}
	}
}

func toSpan(strSpan string) *Span {
	span := new(Span)
	parts := strings.Split(strSpan, " ")
	if len(parts) == 5 {
		calls := strings.Split(parts[4], "->")
		span.Start = parts[0]
		span.End = parts[1]
		span.Trace = parts[2]
		span.Service = parts[3]
		span.Caller = calls[0]
		span.Span = calls[1]
		span.Calls = make([]Span, 0)
	}

	return span
}

func remove(s []*Span, i int) []*Span {
	s[i] = new(Span)
	return s
}
