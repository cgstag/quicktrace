package quicktrace

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_ToTrace(t *testing.T) {
	var strSpan = "2013-10-23T10:12:35.019Z 2013-10-23T10:12:35.019Z ngr7jl6y service2 5dm5aee3->z35dizqs"
	// Test if struct correctly initialized
	if reflect.TypeOf(toSpan(strSpan).Calls).String() != "[]quicktrace.Span" {
		t.Errorf("toSpan method incorrect, got: %s, want: %s.", reflect.TypeOf(toSpan(strSpan).Calls).String(), "*[]quicktrace.Span")
	}
	// Test if values correctly assigned
	if toSpan(strSpan).Trace != "ngr7jl6y" {
		t.Errorf("toSpan method incorrect, got: %s, want: %s.", toSpan(strSpan).Trace, "ngr7jl6y")
	}
	// Test fringe values but explosion respected
	strSpan = "201ß10-23T10:12:35.019Z 2013-10-23T1µ:12:35.019Z Š service2 space†->linebreak†"
	if toSpan(strSpan).Service != "service2" {
		t.Errorf("toSpan method incorrect, special character not handled got: %s, want: %s.", toSpan(strSpan).Service, "service2")
	}
	// Test empty line
	strSpan = "                    "
	if toSpan(strSpan).Service != "" {
		t.Errorf("toSpan method incorrect, empty line not handled got: %s, want: %s.", toSpan(strSpan), "&{      []}")
	}
}

func TestTraceMap_PushSpan(t *testing.T) {
	var strSpan = "2013-10-23T10:12:35.019Z 2013-10-23T10:12:35.019Z ngr7jl6y service2 null->z35dizqs"
	tr := TraceMap{}
	tr.PushSpan(strSpan)
	// Test if Trace correctly assigned to map
	if _, ok := tr[toSpan(strSpan).Trace]; !ok {
		t.Errorf("PushSpan method incorrect, TraceKey not set - got %s, expected: %s.", toSpan(strSpan).Trace, "ngr7jl6y")
	}
	// Test if Root is correctly set during pushmap
	if tr["ngr7jl6y"].Root.Span != "z35dizqs" {
		t.Errorf("PushSpan method incorrect, Root value not set - got: %s, want: %s.", tr["ngr7jl6y"].Root.Span, "z35dizqs")
	}

	// Test if subsequent log in same trace are append
	strSpan = "2013-10-23T10:12:35.019Z 2013-10-23T10:12:35.019Z ngr7jl6y service2 z35dizqs->6yal5to5"
	tr.PushSpan(strSpan)
	if len(tr["ngr7jl6y"].Unsorted) != 1 && tr["ngr7jl6y"].Root.Span != "z35dizqs" {
		t.Errorf("PushSpan method incorrect, Unsorted length is : %d, want: %s.", len(tr["ngr7jl6y"].Unsorted), "> 1")
	}

	// Test new non-Root trace
	strSpan = "2013-10-23T10:12:35.019Z 2013-10-23T10:12:35.019Z ks2vhl7f service12 6lgnmppe->v66pw26l"
	tr.PushSpan(strSpan)
	if len(tr) != 2 {
		t.Errorf("PushSpan method incorrect, map length is : %d, want: %d.", len(tr), 2)
	} else if tr["ks2vhl7f"].Root != nil {
		t.Errorf("PushSpan method incorrect, Root invalid is : %s, expected: %s.", tr["ks2vhl7f"].Root, "%!s(*quicktrace.Span=<nil>)")
	}

	// Test ongoing root trace
	strSpan = "2013-10-23T10:12:35.019Z 2013-10-23T10:12:35.019Z ks2vhl7f service14 null->6lgnmppe"
	tr.PushSpan(strSpan)
	if len(tr) != 2 {
		t.Errorf("PushSpan method incorrect, map length is : %d, want: %d.", len(tr), 2)
	} else if tr["ks2vhl7f"].Root.Span != "6lgnmppe" {
		t.Errorf("PushSpan method incorrect, Root invalid is : %s, expected: %s.", tr["ks2vhl7f"].Root.Span, "6lgnmppe")
	}

	// Test Push empty Span
	strSpan = ""
	tr.PushSpan(strSpan)
	if len(tr) > 2 {
		t.Errorf("PushSpan method incorrect, Hashmap length is : %d, want: %d", len(tr), 2)
	}
}

func TestTrace_QuickTrace(t *testing.T) {
	ctx := new(Context)
	// Test Standard QuickTrace
	tr := createTrace()
	_, err := tr.QuickTrace(ctx)
	if err != nil || tr == nil {
		t.Errorf("QuickTrace method error")
	}
	// Root should have two calls
	if len(tr.Root.Calls) != 2 {
		t.Errorf("QuickTrace method incorrect, Root calls : %d, expected: %d.", len(tr.Root.Calls), 2)
		strTr, _ := json.Marshal(tr)
		t.Errorf("Trace: %s", string(strTr))
	}
	// Test Orphan
	tr = createOrphanTrace()
	tr, err = tr.QuickTrace(ctx)
	if len(ctx.Orphans) == 0 {
		t.Errorf("QuickTrace method not handling orphans properly")
	}

}

func createOrphanTrace() *Trace {
	tr := new(Trace)
	call1 := &Span{
		Service: "service6",
		Trace:   "6mptyd3j",
		Span:    "toa2oj25",
		Start:   "2013-10-23T10:12:37.965Z",
		End:     "2013-10-23T10:12:38.066Z",
		Caller:  "null",
	}
	tr.Id = "6mptyd3j"
	tr.Unsorted = append(tr.Unsorted, call1)
	return tr
}

func createTrace() *Trace {
	tr := new(Trace)
	root := &Span{
		Service: "service7",
		Trace:   "6mptyd3j",
		Span:    "toa2oj25",
		Start:   "2013-10-23T10:12:37.965Z",
		End:     "2013-10-23T10:12:38.066Z",
		Caller:  "cxjahaud",
	}
	call1 := &Span{
		Service: "service6",
		Trace:   "6mptyd3j",
		Span:    "toa2oj25",
		Start:   "2013-10-23T10:12:37.965Z",
		End:     "2013-10-23T10:12:38.066Z",
		Caller:  "null",
	}
	call2 := &Span{
		Service: "service8",
		Trace:   "6mptyd3j",
		Span:    "6qpd42l4",
		Start:   "2013-10-23T10:12:37.965Z",
		End:     "2013-10-23T10:12:38.066Z",
		Caller:  "ndra5nyu",
	}
	call3 := &Span{
		Service: "service5",
		Trace:   "6mptyd3j",
		Span:    "ndra5nyu",
		Start:   "2013-10-23T10:12:37.965Z",
		End:     "2013-10-23T10:12:38.066Z",
		Caller:  "cxjahaud",
	}
	call4 := &Span{
		Service: "service4",
		Trace:   "6mptyd3j",
		Span:    "6lgnmppe",
		Start:   "2013-10-23T10:12:37.965Z",
		End:     "2013-10-23T10:12:38.066Z",
		Caller:  "feypmh4c",
	}
	call5 := &Span{
		Service: "service2",
		Trace:   "6mptyd3j",
		Span:    "feypmh4c",
		Start:   "2013-10-23T10:12:37.965Z",
		End:     "2013-10-23T10:12:38.066Z",
		Caller:  "toa2oj25",
	}
	call6 := &Span{
		Service: "service1",
		Trace:   "6mptyd3j",
		Span:    "cxjahaud",
		Start:   "2013-10-23T10:12:37.965Z",
		End:     "2013-10-23T10:12:38.066Z",
		Caller:  "toa2oj25",
	}
	call7 := &Span{
		Service: "service6",
		Trace:   "6mptyd3j",
		Span:    "jus7kcva",
		Start:   "2013-10-23T10:12:37.965Z",
		End:     "2013-10-23T10:12:38.066Z",
		Caller:  "v66pw26l",
	}
	tr.Id = "6mptyd3j"
	tr.Root = root
	tr.Unsorted = append(tr.Unsorted, call1)
	tr.Unsorted = append(tr.Unsorted, call2)
	tr.Unsorted = append(tr.Unsorted, call3)
	tr.Unsorted = append(tr.Unsorted, call4)
	tr.Unsorted = append(tr.Unsorted, call5)
	tr.Unsorted = append(tr.Unsorted, call6)
	tr.Unsorted = append(tr.Unsorted, call7)
	return tr
}
