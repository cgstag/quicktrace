package quicktrace

import "errors"

func NewNoRootSpanError() error {
	return errors.New("NoRootSpan")
}
