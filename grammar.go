package pgtwixt

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Grammar struct {
	Value func(key, value string) error
}

func (g Grammar) Parse(text string) error {
	text = g.trim(text)

	var err error
	for err == nil && len(text) > 0 {
		text, err = g.parseKey(text)
	}

	return err
}

func (g Grammar) parseKey(text string) (string, error) {
	var r rune
	var w int

	if r, w = utf8.DecodeRuneInString(text); r == '\'' {
		text, key, err := g.parseQuoted(text[w:])
		text = g.trim(text)

		if err != nil {
			return text, err
		}

		if r, w = utf8.DecodeRuneInString(text); r == '=' {
			return g.parseValue(g.trim(text[w:]), key)
		}

		return text, fmt.Errorf("expected '=', got %q", r)
	}

	for i := 0; i < len(text); i += w {
		r, w = utf8.DecodeRuneInString(text[i:])

		if r == '=' {
			return g.parseValue(g.trim(text[i+w:]), text[:i])
		}

		if unicode.IsSpace(r) {
			key := text[:i]
			text = g.trim(text[i:])

			if r, w = utf8.DecodeRuneInString(text); r == '=' {
				return g.parseValue(g.trim(text[w:]), key)
			}

			return text, fmt.Errorf("expected '=', got %q", r)
		}
	}

	return "", errors.New("expected '=' before end")
}

func (g Grammar) parseQuoted(text string) (string, string, error) {
	var escape bool
	var result string
	var r rune
	var w int

	for i := 0; i < len(text); i += w {
		r, w = utf8.DecodeRuneInString(text[i:])

		if escape {
			result += string(r)
			escape = false
		} else {
			switch r {
			default:
				result += string(r)
			case '\\':
				escape = true
			case '\'':
				return g.trim(text[i+w:]), result, nil
			}
		}
	}

	return text, "", errors.New("expected matching quote before end")
}

func (g Grammar) parseValue(text string, key string) (string, error) {
	var r rune
	var w int

	if r, w = utf8.DecodeRuneInString(text); r == '\'' {
		text, value, err := g.parseQuoted(text[w:])
		text = g.trim(text)

		if err != nil {
			return text, err
		}

		return text, g.Value(key, value)
	}

	for i := 0; i < len(text); i += w {
		r, w = utf8.DecodeRuneInString(text[i:])

		if unicode.IsSpace(r) {
			return g.trim(text[i:]), g.Value(key, text[:i])
		}
	}

	if len(text) < 1 {
		return "", errors.New("expected value before end")
	}

	return "", g.Value(key, text)
}

func (Grammar) trim(text string) string {
	return strings.TrimLeftFunc(text, unicode.IsSpace)
}
