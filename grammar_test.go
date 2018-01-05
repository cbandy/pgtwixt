package pgtwixt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrammarEmptyDoesNothing(t *testing.T) {
	t.Parallel()
	assert.NoError(t, Grammar{}.Parse(""))
}

func TestGrammarRequiresValue(t *testing.T) {
	t.Parallel()
	assert.Error(t, Grammar{}.Parse(`a`))
	assert.Error(t, Grammar{}.Parse(`a `))
	assert.Error(t, Grammar{}.Parse(`'a'`))
	assert.Error(t, Grammar{}.Parse(`a=`))
}

func TestGrammarRequiresMatchingQuotes(t *testing.T) {
	t.Parallel()
	assert.Error(t, Grammar{}.Parse("'a"))
	assert.Error(t, Grammar{}.Parse("'a'='b"))
}

func TestGrammarValues(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		input    string
		expected []string
	}{
		{`a=b`, []string{`a`, `b`}},
		{` a=b`, []string{`a`, `b`}},
		{`a=b `, []string{`a`, `b`}},
		{`a =b`, []string{`a`, `b`}},
		{`a= b`, []string{`a`, `b`}},
		{`a = b`, []string{`a`, `b`}},
		{` a = b `, []string{`a`, `b`}},

		{`ab=cd ef=gh`, []string{`ab`, `cd`, `ef`, `gh`}},
		{` ab=cd ef=gh`, []string{`ab`, `cd`, `ef`, `gh`}},
		{`ab=cd ef=gh `, []string{`ab`, `cd`, `ef`, `gh`}},
		{` ab = cd ef = gh `, []string{`ab`, `cd`, `ef`, `gh`}},

		{`a=''`, []string{`a`, ``}},
		{`a= 'b c' `, []string{`a`, `b c`}},
		{`a= 'b\'c'`, []string{`a`, `b'c`}},
		{`a= 'b\\c'`, []string{`a`, `b\c`}},
		{`a= 'b\c'`, []string{`a`, `bc`}},
		{`a= 'b\'c' d=e`, []string{`a`, `b'c`, `d`, `e`}},
		{` 'a b' = 'c d' `, []string{`a b`, `c d`}},
	} {
		t.Run(tt.input, func(t *testing.T) {
			var result []string
			g := Grammar{Value: func(key, value string) error {
				result = append(result, key, value)
				return nil
			}}

			require.NoError(t, g.Parse(tt.input))
			assert.Equal(t, tt.expected, result)
		})
	}
}
