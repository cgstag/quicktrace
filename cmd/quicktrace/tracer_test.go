package quicktrace

import (
	"os"
	"quicktrace/cli"
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
	opt := &cli.Options{
		Input:  "/tmp/log.txt",
		Output: "/tmp/trace.txt",
		Stats:  false,
		Help:   false,
	}
	tr := NewTracer(opt, new(cli.Context))
	tr.pushSpan(strSpan)
	// Test if Trace correctly assigned to map
	if _, ok := tr.TraceMap[toSpan(strSpan).Trace]; !ok {
		t.Errorf("PushSpan method incorrect, TraceKey not set - got %s, expected: %s.", toSpan(strSpan).Trace, "ngr7jl6y")
	}
	// Test if Root is correctly set during pushmap
	if tr.TraceMap["ngr7jl6y"].Root.Span != "z35dizqs" {
		t.Errorf("PushSpan method incorrect, Root value not set - got: %s, want: %s.", tr.TraceMap["ngr7jl6y"].Root.Span, "z35dizqs")
	}

	// Test if subsequent log in same trace are append
	strSpan = "2013-10-23T10:12:35.019Z 2013-10-23T10:12:35.019Z ngr7jl6y service2 z35dizqs->6yal5to5"
	tr.pushSpan(strSpan)
	if len(tr.TraceMap["ngr7jl6y"].Unsorted) != 1 && tr.TraceMap["ngr7jl6y"].Root.Span != "z35dizqs" {
		t.Errorf("PushSpan method incorrect, Unsorted length is : %d, want: %s.", len(tr.TraceMap["ngr7jl6y"].Unsorted), "> 1")
	}

	// Test new non-Root trace
	strSpan = "2013-10-23T10:12:35.019Z 2013-10-23T10:12:35.019Z ks2vhl7f service12 6lgnmppe->v66pw26l"
	tr.pushSpan(strSpan)
	if len(tr.TraceMap) != 2 {
		t.Errorf("PushSpan method incorrect, map length is : %d, want: %d.", len(tr.TraceMap), 2)
	} else if tr.TraceMap["ks2vhl7f"].Root != nil {
		t.Errorf("PushSpan method incorrect, Root invalid is : %s, expected: %s.", tr.TraceMap["ks2vhl7f"].Root, "%!s(*quicktrace.Span=<nil>)")
	}

	// Test ongoing root trace
	strSpan = "2013-10-23T10:12:35.019Z 2013-10-23T10:12:35.019Z ks2vhl7f service14 null->6lgnmppe"
	tr.pushSpan(strSpan)
	if len(tr.TraceMap) != 2 {
		t.Errorf("PushSpan method incorrect, map length is : %d, want: %d.", len(tr.TraceMap), 2)
	} else if tr.TraceMap["ks2vhl7f"].Root.Span != "6lgnmppe" {
		t.Errorf("PushSpan method incorrect, Root invalid is : %s, expected: %s.", tr.TraceMap["ks2vhl7f"].Root.Span, "6lgnmppe")
	}

	// Test Push empty Span
	strSpan = ""
	tr.pushSpan(strSpan)
	if len(tr.TraceMap) > 2 {
		t.Errorf("PushSpan method incorrect, Hashmap length is : %d, want: %d", len(tr.TraceMap), 2)
	}
}

func TestTracer_ParseAndTraceStdout(t *testing.T) {
	opt := &cli.Options{
		Input:  "/tmp/log.txt",
		Output: "stdout",
		Stats:  true,
		Help:   false,
	}

	input, err := os.Create(opt.Input)
	if err != nil {
		panic(err)
	}
	_, err = input.WriteString("2013-10-23T10:12:35.958Z 2013-10-23T10:12:36.012Z 7pcahxgp service7 null->dhbwylzy")
	tr := NewTracer(opt, new(cli.Context))
	if tr.ParseAndTrace() != nil {
		t.Errorf("Error running ParseAndTrace: %s", err.Error())
	}
	_ = input.Close()
}

func TestTracer_ParseAndTraceFile(t *testing.T) {
	opt := &cli.Options{
		Input:  "/tmp/log.txt",
		Output: "/tmp/trace.txt",
		Stats:  false,
		Help:   false,
	}

	input, err := os.Create(opt.Input)
	if err != nil {
		panic(err)
	}
	_, err = input.WriteString("2013-10-23T10:12:35.958Z 2013-10-23T10:12:36.012Z 7pcahxgp service7 null->dhbwylzy")
	tr := NewTracer(opt, new(cli.Context))
	if tr.ParseAndTrace() != nil {
		t.Errorf("Error running ParseAndTrace: %s", err.Error())
	}
}

func TestTracer_ParseAndTraceOrphan(t *testing.T) {
	opt := &cli.Options{
		Input:  "/tmp/log.txt",
		Output: "/tmp/trace.txt",
		Stats:  false,
		Help:   false,
	}

	input, err := os.Create(opt.Input)
	if err != nil {
		panic(err)
	}
	_, err = input.WriteString("2013-10-23T10:12:35.958Z 2013-10-23T10:12:36.012Z 7pcahxgp service7 dhbwylas->dhbwylzy")
	tr := NewTracer(opt, new(cli.Context))
	if tr.ParseAndTrace() != nil {
		t.Errorf("Error running ParseAndTrace: %s", err.Error())
	}
}

func TestTracer_ParseAndTraceMultiple(t *testing.T) {
	opt := &cli.Options{
		Input:  "/tmp/log.txt",
		Output: "/tmp/trace.txt",
		Stats:  false,
		Help:   false,
	}

	input, err := os.Create(opt.Input)
	if err != nil {
		panic(err)
	}
	_, err = input.WriteString("2013-10-23T10:12:35.271Z 2013-10-23T10:12:35.471Z eckakaau service6 null->bm6il56t\n")
	_, err = input.WriteString("2013-10-23T10:12:35.293Z 2013-10-23T10:12:35.302Z eckakaau service7 zfjlsiev->d6m3shqy\n")
	_, err = input.WriteString("2013-10-23T10:12:37.708Z 2013-10-23T10:12:37.724Z nhdyl6hs service9 34z4ib4a->ai67mto3\n")
	_, err = input.WriteString("2013-10-23T10:12:37.127Z 2013-10-23T10:12:37.891Z nhdyl6hs service7 null->34z4ib4a\n")
	_, err = input.WriteString("2013-10-23T10:12:37.723Z 2013-10-23T10:12:37.724Z rjopvy3w service6 pncdppmf->zltdmrcv\n")
	tr := NewTracer(opt, new(cli.Context))
	if tr.ParseAndTrace() != nil {
		t.Errorf("Error running ParseAndTraceMultiple: %s", err.Error())
	}
}

func TestTracer_ParseAndTraceWrongOutput(t *testing.T) {
	opt := &cli.Options{
		Input:  "/tmp/log.txt",
		Output: "/asd/trace.txt",
		Stats:  false,
		Help:   false,
	}

	input, err := os.Create(opt.Input)
	if err != nil {
		panic(err)
	}
	_, err = input.WriteString("2013-10-23T10:12:35.958Z 2013-10-23T10:12:36.012Z 7pcahxgp service7 null->dhbwylzy")
	tr := NewTracer(opt, new(cli.Context))

	if tr.ParseAndTrace() == nil {
		t.Errorf("This ParseAndTrace should return an error")
	}
}

func TestTracer_ParseAndTraceStdin(t *testing.T) {
	opt := &cli.Options{
		Input:  "stdin",
		Output: "/tmp/trace.txt",
		Stats:  false,
		Help:   false,
	}

	tr := NewTracer(opt, new(cli.Context))

	if tr.ParseAndTrace() != nil {
		t.Errorf("This ParseAndTrace should return an error")
	}
}

func TestTracer_ParseAndTraceFileDoesntExist(t *testing.T) {
	opt := &cli.Options{
		Input:  "/abc/log.txt",
		Output: "/tmp/trace.txt",
		Stats:  false,
		Help:   false,
	}

	tr := NewTracer(opt, new(cli.Context))

	if tr.ParseAndTrace() == nil {
		t.Errorf("This ParseAndTrace should return an error")
	}
}
