package quicktrace

import "time"

type Context struct {
	Start   time.Time `json:"-"`
	Elapsed time.Time `json:"elapsed"`
	Errors  int       `json:"errors"`
	Orphans []*Trace  `json:"orphans"`
}

type TraceMap map[string]Trace

type Trace struct {
	Id       string  `json:"id"`
	Root     *Span   `json:"root"`
	Unsorted []*Span `json:"u"`
}

type Span struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	Trace   string `json:"-"`
	Service string `json:"service"`
	Span    string `json:"span"`
	Caller  string `json:"-"`
	Calls   []Span `json:"calls"`
}
