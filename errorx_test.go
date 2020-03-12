package errorx

import (
	stderrors "errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapNil(t *testing.T) {
	got := Wrap(nil, "no error")
	require.Nilf(t, got, "expected nil: got %#v", got)
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrap(io.EOF, "read error"), "client error", "client error: read error: EOF"},
		{Wrap(*Wrap(io.EOF, "inner"), "outer"), "outermost", "outermost: outer: inner: EOF"},
	}

	for _, tt := range tests {
		got := Wrap(tt.err, tt.message).Error()
		assert.Equalf(t, tt.want, got, "got: %v, want %v", got, tt.want)
	}
}

func TestUnwrap(t *testing.T) {
	err := New("test")
	// wErr := Wrap(err, "invalid")
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "wrapped",
			args: args{Wrap(err, "wrapping message")},
			want: err,
		},
		{
			name: "std errors compatibility",
			args: args{err: fmt.Errorf("wrap: %w", err)},
			want: err,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := stderrors.Unwrap(tt.args.err); !reflect.DeepEqual(err, tt.want) {
				assert.Failf(t, "failed", "want %v got %v", tt.want, err)
			}

			unwrapped := stderrors.Unwrap(tt.args.err)
			assert.True(t, stderrors.Is(unwrapped, tt.want))
		})
	}
}

func TestPrint(t *testing.T) {
	err := Wrap(Wrap(New("inner"), "outer"), "outermost")

	assert.Equal(t, "outermost: outer: inner", err.Error())
	t.Logf("%+v", err)
}
