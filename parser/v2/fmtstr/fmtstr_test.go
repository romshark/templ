package fmtstr_test

import (
	"testing"

	"github.com/a-h/parse"
	"github.com/a-h/templ/parser/v2/fmtstr"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in           string
		want         fmtstr.FormatString
		expectErr    error
		expectSuffix string
	}{
		{
			in:   ``,
			want: fmtstr.FormatString{},
		},
		{
			in:           `not a format string`,
			want:         fmtstr.FormatString{},
			expectSuffix: `not a format string`,
		},
		{
			in: `"hello"`,
			want: fmtstr.FormatString{
				Parts: []fmtstr.Part{
					{Type: fmtstr.PartTypeLiteral, Value: "hello"},
				},
			},
		},
		{
			in: `"hello", suffix`,
			want: fmtstr.FormatString{
				Parts: []fmtstr.Part{
					{Type: fmtstr.PartTypeLiteral, Value: "hello"},
				},
			},
			expectSuffix: `, suffix`,
		},
		{
			in: `"foo %d bar %s"`,
			want: fmtstr.FormatString{
				Parts: []fmtstr.Part{
					{Type: fmtstr.PartTypeLiteral, Value: "foo "},
					{Type: fmtstr.PartTypeBase10, Value: "%d"},
					{Type: fmtstr.PartTypeLiteral, Value: " bar "},
					{Type: fmtstr.PartTypeString, Value: "%s"},
				},
			},
		},
		{
			in: `"%% escaped"`,
			want: fmtstr.FormatString{
				Parts: []fmtstr.Part{
					{Type: fmtstr.PartTypeLiteral, Value: "% escaped"},
				},
			},
		},
		{
			in: `"pi = %.2f"`,
			want: fmtstr.FormatString{
				Parts: []fmtstr.Part{
					{Type: fmtstr.PartTypeLiteral, Value: "pi = "},
					{Type: fmtstr.PartTypeFloat, Value: "%.2f"},
				},
			},
		},
		{
			in: `"bool=%t hex=%x oct=%o char=%c ptr=%p gen=%v"`,
			want: fmtstr.FormatString{
				Parts: []fmtstr.Part{
					{Type: fmtstr.PartTypeLiteral, Value: "bool="},
					{Type: fmtstr.PartTypeBool, Value: "%t"},
					{Type: fmtstr.PartTypeLiteral, Value: " hex="},
					{Type: fmtstr.PartTypeHex, Value: "%x"},
					{Type: fmtstr.PartTypeLiteral, Value: " oct="},
					{Type: fmtstr.PartTypeOctal, Value: "%o"},
					{Type: fmtstr.PartTypeLiteral, Value: " char="},
					{Type: fmtstr.PartTypeChar, Value: "%c"},
					{Type: fmtstr.PartTypeLiteral, Value: " ptr="},
					{Type: fmtstr.PartTypePointer, Value: "%p"},
					{Type: fmtstr.PartTypeLiteral, Value: " gen="},
					{Type: fmtstr.PartTypeGeneric, Value: "%v"},
				},
			},
		},
		{
			in:        `"trailing %`, // no verb after '%'
			expectErr: fmtstr.ErrTrailingPercent,
		},
		{
			in:        `"trailing %"`, // no verb after '%'
			expectErr: fmtstr.ErrUntermVerb,
		},
		{
			in:        `"%"`, // only percent
			expectErr: fmtstr.ErrUntermVerb,
		},
		{
			in:        `"unterminated`, // missing closing quote
			expectErr: fmtstr.ErrUnterm,
		},
	}

	for _, tt := range tests {
		pi := parse.NewInput(tt.in)
		fstr, err := fmtstr.Parse(pi)
		if tt.expectErr != nil {
			require.ErrorIs(t, err, tt.expectErr, "input=%q", tt.in)
			require.Zero(t, fstr)
			continue
		}
		require.NoError(t, err, "input=%q", tt.in)
		require.Equal(t, tt.want, fstr, "input=%q", tt.in)
		rest, _ := pi.Peek(-1)
		require.Equal(t, tt.expectSuffix, rest, "input=%q", tt.in)
	}
}
