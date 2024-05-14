package jsonpointer

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

// Pointer is a JSON Pointer as defined in RFC 6901.
type Pointer []string

// Parse returns a new pointer from the given string representation.
func Parse(s string) (Pointer, error) {
	if s == "" {
		return nil, ErrInvalidPointer
	}
	if s[0] != '/' {
		return nil, ErrInvalidPointer
	}
	s = s[1:]
	if s == "" {
		return Pointer{}, nil
	}
	p := Pointer{}
	for _, token := range strings.Split(s, "/") {
		if token == "" {
			return nil, ErrInvalidPointer
		}
		var err error
		token, err = unescape(token)
		if err != nil {
			return nil, err
		}
		p = append(p, token)
	}
	return p, nil
}

func (p Pointer) String() string {
	if len(p) == 0 {
		return "/"
	}
	var sb strings.Builder
	for _, token := range p {
		sb.WriteByte('/')
		sb.WriteString(escape(token))
	}
	return sb.String()
}

// Append appends the given tokens to the pointer.
func (p Pointer) Append(tokens ...string) Pointer {
	return append(p, tokens...)
}

var ErrInvalidPointer = errors.New("invalid JSON pointer")

var unescapeReplacer = strings.NewReplacer("~1", "/", "~0", "~")
var unescapeInvalid = regexp.MustCompile(`~([^01]|$)`)

func unescape(token string) (string, error) {
	// Check for invalid escape sequences
	if unescapeInvalid.MatchString(token) {
		return "", ErrInvalidPointer
	}
	return unescapeReplacer.Replace(token), nil
}

var escapeReplacer = strings.NewReplacer("~", "~0", "/", "~1")

func escape(token string) string {
	return escapeReplacer.Replace(token)
}

func (p Pointer) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func New(elements ...string) Pointer {
	return elements
}
