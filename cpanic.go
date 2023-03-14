// cpanic is a package for gracefully handling panics in Go.
//
// There are two main functions in this package: `Recover` and `Forward`. `Recover` can
// be provided a callback that executes when a panic is recovered. `Forward` can be
// provided an error pointer that will be set to a `*Panic` type when a panic is
// recovered. `Forward` is useful when you want to return an error from a function
// that may panic. `Go` is an application of `Forward` that accepts a function that may
// panic and returns an error instead.
package cpanic

import (
	"fmt"
	"runtime"
	"time"
)

// Handler is a function that handles a panic.
type Handler func(p *Panic)

// Recover is a defer function that recovers from a panic and calls the handler. If no
// handler is provided, `recover` is never called and the panic is allowed to continue.
func Recover(handler Handler) {
	if handler == nil {
		return
	}

	if value := recover(); value != nil {
		handler(New(value))
	}
}

// Go calls the provided function and recovers from any panics. If the function panics,
// the error returned will be a `*Panic` type otherwise the error returned, if any, will
// be from the function.
func Go(fn func() error) (err error) {
	defer Forward(&err)
	return fn()
}

// Forward is a defer function that recovers from a panic and sets the provided error
// pointer to a `*Panic` type. If the error pointer is nil, `recover` is never called
// and the panic is allowed to continue.
func Forward(errPtr *error) {
	if errPtr == nil {
		return
	}

	if value := recover(); value != nil {
		if *errPtr == nil {
			*errPtr = New(value)
		}
	}
}

// Panic is an error type that is returned when a panic is recovered.
type Panic struct {
	// Time is the time the panic occurred.
	Time time.Time `json:"time" yaml:"time"`
	// Value is the value of the panic. This is usually a `string` or an `error` but can
	// be any type.
	Value interface{} `json:"value" yaml:"value"`
	// Trace is the stack trace of all goroutines at the time of the panic.
	Trace string `json:"trace" yaml:"trace"`
}

// Error implements the `error` interface and returns a string representation of the
// panic value. This does not include the stack traces.
func (p *Panic) Error() string {
	return fmt.Sprintf("panic: %v", p.Value)
}

// String implements the `fmt.Stringer` interface and returns a string representation
// of the panic with all of the collected stack traces from when the panic occurred.
func (p *Panic) String() string {
	return fmt.Sprintf("%s\n\n%s", p.Error(), p.Trace)
}

// New creates a new `*Panic` from the provided value. Stack traces for all goroutines
// are collected during construction. This is expected to be used during panic recovery.
func New(v interface{}) *Panic {
	var trace [1 << 16]byte
	n := runtime.Stack(trace[:], true)
	p := &Panic{
		Time:  time.Now(),
		Value: v,
		Trace: string(trace[:n]),
	}
	return p
}
