package thorlog

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/NextronSystems/jsonlog"
)

type StringWithEncoding struct {
	EncodedData string         `json:"data"`
	Encoding    StringEncoding `json:"encoding"`
}

type StringEncoding string

const (
	Plain StringEncoding = "plain"
	Hex   StringEncoding = "hex"
)

// Encode encodes the given data into a StringWithEncoding,
// choosing the most appropriate encoding based on its content.
func Encode(s []byte) StringWithEncoding {
	if utf8.Valid(s) {
		return StringWithEncoding{
			EncodedData: string(s),
			Encoding:    Plain,
		}
	} else {
		return StringWithEncoding{
			EncodedData: hex.EncodeToString(s),
			Encoding:    Hex,
		}
	}
}

// EncodeString encodes the given data into a StringWithEncoding,
// choosing the most appropriate encoding based on its content.
func EncodeString(s string) StringWithEncoding {
	if utf8.ValidString(s) {
		return StringWithEncoding{
			EncodedData: s,
			Encoding:    Plain,
		}
	} else {
		return StringWithEncoding{
			EncodedData: hex.EncodeToString([]byte(s)),
			Encoding:    Hex,
		}
	}
}

// Data returns the raw byte sequence represented by the StringWithEncoding.
func (s StringWithEncoding) Data() []byte {
	switch s.Encoding {
	case Plain:
		return []byte(s.EncodedData)
	case Hex:
		data, err := hex.DecodeString(s.EncodedData)
		if err != nil {
			return []byte("<invalid hex data: " + err.Error() + ">")
		}
		return data
	default:
		return []byte(fmt.Sprintf("<unknown encoding %s> %s", s.Encoding, s.EncodedData))
	}
}

var notOnlyASCII = regexp.MustCompile(`[^\x20-\x7E\x0d\x0a\x09]+`) // printable chars + \r,\n,\t

// String returns a human-readable representation of the encoded string.
// The representation is guaranteed to be valid UTF-8.
func (s StringWithEncoding) String() string {
	data := s.decode()
	if needsQuoting.MatchString(data) {
		return quote(data)
	}
	return data
}

// decode returns the plain text, after decoding it from UTF-16, if applicable.
func (s StringWithEncoding) decode() string {
	plaintext := s.Data()

	if decoded, ok := attemptDecodeUTF16(plaintext); ok {
		return decoded
	}

	return string(plaintext)
}

// attemptDecodeUTF16 tries to decode the given byte slice as UTF-16 and checks
// whether the decoded string contains non-ASCII characters.
// It returns the decoded string and a boolean indicating whether the decoding was successful.
func attemptDecodeUTF16(b []byte) (string, bool) {
	// Try UTF16 encoding
	if len(b) > 1 && b[0] == 0xFF && b[1] == 0xFE {
		// Remove byte order mark
		b = b[2:]
	}
	if len(b) > 0 && b[0] == 0 {
		// Might be UTF16 shifted by one byte
		b = b[1:]
	}
	decodedUtf16, _ := decodeUTF16(b)
	if !notOnlyASCII.MatchString(decodedUtf16) && len(decodedUtf16) > 0 {
		// Can cleanly be rendered as UTF-16
		return decodedUtf16, true
	}
	return "", false
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

func quote(s string) string {
	s = escaper.Replace(s)
	var quotedString bytes.Buffer
	quotedString.WriteString(`"`)
	for _, char := range []byte(s) {
		if char < 0x20 || char > 0x7E { // non ASCII
			quotedString.WriteString("\\x")
			quotedString.WriteString(hex.EncodeToString([]byte{char}))
		} else {
			quotedString.WriteByte(char)
		}
	}
	quotedString.WriteString(`"`)
	return quotedString.String()
}

// MatchString describes a sequence of bytes in an object
// that was matched on by a signature.
type MatchString struct {
	// Match contains the bytes that were matched.
	Match StringWithEncoding `json:"data"`
	// Context contains the bytes surrounding the matched bytes.
	// This may be missing if no context is available.
	Context *StringWithEncoding `json:"context,omitempty"`
	// Offset contains the Match's offset within the Field
	// where the data was matched.
	Offset *uint64 `json:"offset,omitempty"`
	// Field points to the field within the object that was matched on.
	Field      *jsonlog.Reference `json:"field,omitempty"`
	HideOffset bool               `json:"-"`
}

var needsQuoting = regexp.MustCompile(`[^\x21\x23-\x7E]`)

func (f MatchString) String() string {
	matchString := f.Match.String()
	if f.Context != nil {
		matchString += " in " + f.Context.String()
	}
	if f.Offset != nil {
		// Only show the offset if this match does not encompass the full field and it's not explicitly hidden
		var showOffset = !f.HideOffset
		if f.Field != nil && *f.Offset == 0 {
			if targetString, isString := f.Field.Value().(string); isString {
				if targetString == f.Match.EncodedData {
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

// MatchStrings is a list of matching byte sequences that explains
// why a specific signature matched on an object.
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
