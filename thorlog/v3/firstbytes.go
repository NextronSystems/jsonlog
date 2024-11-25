package thorlog

import (
	"encoding/hex"
	"encoding/json"
	"strings"
)

type FirstBytes []byte

func trimNonPrintableChars(data []byte) string {
	var builder strings.Builder
	builder.Grow(len(data))
	for _, char := range data {
		if char >= 0x20 && char <= 0x7E {
			builder.WriteByte(char)
		}
	}
	return builder.String()
}

func (f FirstBytes) String() string {
	return hex.EncodeToString(f) + " / " + trimNonPrintableChars(f)
}

type firstBytesJson struct {
	Hex   string `json:"hex"`
	Ascii string `json:"ascii"`
}

func (f FirstBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(firstBytesJson{hex.EncodeToString(f), trimNonPrintableChars(f)})
}

func (f *FirstBytes) UnmarshalJSON(data []byte) error {
	var jsonStruct firstBytesJson
	err := json.Unmarshal(data, &jsonStruct)
	if err != nil {
		return err
	}
	unhexedData, err := hex.DecodeString(jsonStruct.Hex)
	if err != nil {
		return err
	}
	*f = unhexedData
	return nil
}

func (f FirstBytes) JSONSchemaAlias() any { return firstBytesJson{} }
