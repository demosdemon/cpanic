package cpanic_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/demosdemon/cpanic"
)

func TestForward(t *testing.T) {
	tests := []struct {
		name   string
		fn     func() error
		panics bool
		errMsg string
	}{
		{
			name: "does nothing",
			fn:   func() error { return nil },
		},
		{
			name:   "returns panic",
			panics: true,
			fn:     func() error { panic("not at a disco") },
			errMsg: "panic: not at a disco",
		},
		{
			name:   "doesn't override error",
			errMsg: "test",
			fn: func() (err error) {
				defer cpanic.Forward(&err)
				err = errors.New("test")
				panic("not at a disco")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				assert.Panics(t, func() {
					defer cpanic.Forward(nil)
					_ = tt.fn()
				})
			}

			err := cpanic.Go(tt.fn)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errMsg)
			}
		})
	}
}
