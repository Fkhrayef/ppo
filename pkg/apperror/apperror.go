package apperror

import "fmt"

type Kind int

const (
	KindNotFound Kind = iota
	KindValidation
	KindConflict
	KindUpstream
	KindInternal
)

type Error struct {
	Kind    Kind
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.Err }

func NewNotFound(msg string) *Error {
	return &Error{Kind: KindNotFound, Message: msg}
}

func NewValidation(msg string) *Error {
	return &Error{Kind: KindValidation, Message: msg}
}

func NewConflict(msg string) *Error {
	return &Error{Kind: KindConflict, Message: msg}
}

func NewUpstream(msg string, err error) *Error {
	return &Error{Kind: KindUpstream, Message: msg, Err: err}
}

func NewInternal(msg string, err error) *Error {
	return &Error{Kind: KindInternal, Message: msg, Err: err}
}
