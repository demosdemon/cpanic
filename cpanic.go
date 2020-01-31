package cpanic

import (
	"fmt"
	"runtime"
	"time"
)

type Handler func(p *Panic)

func Recover(handler Handler) {
	if handler == nil {
		return
	}

	if value := recover(); value != nil {
		var trace [1 << 16]byte
		n := runtime.Stack(trace[:], true)
		p := &Panic{
			Time:  time.Now(),
			Value: value,
			Trace: string(trace[:n]),
		}
		handler(p)
	}
}

type Panic struct {
	Time  time.Time   `json:"time" yaml:"time"`
	Value interface{} `json:"value" yaml:"value"`
	Trace string      `json:"trace" yaml:"trace"`
}

func (p *Panic) Error() string {
	return fmt.Sprintf("panic: %v", p.Value)
}

func (p *Panic) String() string {
	return fmt.Sprintf("%s\n\n%s", p.Error(), p.Trace)
}
