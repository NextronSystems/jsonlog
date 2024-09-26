package thorlog

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/NextronSystems/jsonlog"
)

type MatchData struct {
	Data    []byte
	FullHex bool
}

func (f MatchData) MarshalJSON() ([]byte, error) {
	matchingString := f.String()
	return InvalidUnicodeString(matchingString).MarshalJSON()
}

func (f *MatchData) UnmarshalJSON(data []byte) error {
	var matchingString string
	err := json.Unmarshal(data, &matchingString)
	if err != nil {
		return err
	}
	f.Data = []byte(matchingString)
	return nil
}

func (f MatchData) JSONSchemaAlias() any {
	return ""
}

var notOnlyASCII = regexp.MustCompile(`[^\x20-\x7E\x0d\x0a\x09]+`) // printable chars + \r,\n,\t

func (f MatchData) String() string {
	if f.FullHex {
		return hex.EncodeToString(f.Data)
	}
	data := f.Data
	matchingString := string(data) // Try to directly convert

	if !f.FullHex && notOnlyASCII.MatchString(matchingString) { // Check if any non-printable chars occur
		var utf16Data = data
		// Try UTF16 encoding
		if len(utf16Data) > 1 && utf16Data[0] == 0xFF && utf16Data[1] == 0xFE {
			// Remove byte order mark
			utf16Data = utf16Data[2:]
		}
		if len(utf16Data) > 0 && utf16Data[0] == 0 {
			// Might be UTF16 shifted by one byte
			utf16Data = utf16Data[1:]
		}
		matchingString, _ = decodeUTF16(utf16Data)
		if notOnlyASCII.MatchString(matchingString) || len(matchingString) == 0 {
			// Can't cleanly be rendered as UTF-16
			matchingString = string(data)
		}
	}
	return matchingString
}

// https://gist.github.com/bradleypeabody/185b1d7ed6c0c2ab6cec
func decodeUTF16(b []byte) (string, error) {
	if len(b)%2 != 0 {
		b = b[:len(b)-1]
	}
	u16s := make([]uint16, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)
	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}
	return ret.String(), nil
}

func (f MatchData) QuotedString() string {
	matchingString := f.String()
	matchingString = escaper.Replace(matchingString)
	var replacedString bytes.Buffer
	for _, char := range []byte(matchingString) {
		if char < 0x20 || char > 0x7E { // non ASCII
			replacedString.WriteString("\\x")
			replacedString.WriteString(hex.EncodeToString([]byte{char}))
		} else {
			replacedString.WriteByte(char)
		}
	}
	matchingString = replacedString.String()
	matchingString = fmt.Sprintf("\"%s\"", matchingString)
	return matchingString
}

type MatchString struct {
	Match      MatchData          `json:"data"`
	Context    *MatchData         `json:"context,omitempty"`
	Offset     *uint64            `json:"offset,omitempty"`
	Field      *jsonlog.Reference `json:"field,omitempty"`
	HideOffset bool               `json:"-"`
}

var needsQuoting = regexp.MustCompile(`[^\x21\x23-\x7E]`)

func (f MatchString) String() string {
	var matchString string
	if needsQuoting.MatchString(f.Match.String()) && !f.Match.FullHex {
		matchString += f.Match.QuotedString()
	} else {
		matchString += f.Match.String()
	}
	if f.Context != nil {
		matchString += " in "
		if needsQuoting.MatchString(f.Context.String()) && !f.Context.FullHex {
			matchString += f.Context.QuotedString()
		} else {
			matchString += f.Context.String()
		}
	}
	if f.Offset != nil {
		// Only show the offset if this match does not encompass the full field and it's not explicitly hidden
		var showOffset = !f.HideOffset
		if f.Field != nil && *f.Offset == 0 {
			if targetString, isString := f.Field.Value().(string); isString {
				if targetString == string(f.Match.Data) {
					showOffset = false
				}
			}
		}
		if showOffset {
			matchString += fmt.Sprintf(" at %#x", *f.Offset)
		}
	}
	if f.Field != nil {
		matchString += " in " + f.Field.String()
	}

	return matchString
}

type MatchStrings []MatchString

const maxMatchStrings = 30

func (f MatchStrings) String() string {
	if len(f) == 0 {
		return "(none)"
	}
	if len(f) == 1 {
		return f[0].String()
	}
	var formattedStrings = make([]string, len(f))
	for i := range f {
		formattedStrings[i] = fmt.Sprintf("Str%d: %s", i+1, f[i].String())
		if i == maxMatchStrings {
			formattedStrings = append(formattedStrings, "... (strings truncated)")
			break
		}
	}
	return strings.Join(formattedStrings, " ")
}
