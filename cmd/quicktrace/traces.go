package quicktrace

import (
	"fmt"
	"os"
	"strings"
)

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

func (tr *Trace) QuickTrace(ctx *Context) (*Trace, error) {
	var err error
	tr.Root, tr.Unsorted, err = matchSpan(tr.Root, tr.Unsorted)
	if err != nil && err.Error() == "NoRootSpan" {
		fmt.Fprintf(os.Stderr, "Orphan line detected for Trace: %s \n", tr.Id)
		ctx.Orphans = append(ctx.Orphans, tr)
		return nil, err
	}
	return tr, nil
}

/**
 * The challenge here is parse the log entry and assign it to the Hashmap
 */
func (tr TraceMap) PushSpan(strSpan string) {
	span := toSpan(strSpan)
	if span.Span == "" {
		return
	}
	if _, ok := tr[span.Trace]; !ok {
		// This trace is new
		tr[span.Trace] = Trace{
			Id: span.Trace,
		}
		var current = tr[span.Trace]
		if span.Caller == "null" {
			// The span is root
			current.Root = span
			tr[span.Trace] = current
		} else {
			// Add to Unsorted list of that trace
			current.Unsorted = append(tr[span.Trace].Unsorted, span)
			tr[span.Trace] = current
		}
	} else {
		// This trace is ongoing
		var onGoingTrace = tr[span.Trace]
		if span.Caller == "null" {
			// The span is root
			onGoingTrace.Root = span
			tr[span.Trace] = onGoingTrace
		} else {
			// Add to Unsorted list of that trace
			onGoingTrace.Unsorted = append(tr[span.Trace].Unsorted, span)
			tr[span.Trace] = onGoingTrace
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
