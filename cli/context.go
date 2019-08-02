package cli

import "time"

type Context struct {
	Start   time.Time
	IO      int64
	Elapsed time.Duration
}

func NewContext() *Context {
	return &Context{
		Start: time.Now(),
		IO:    0,
	}
}
