package common

import (
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/mod/semver"
)

type Version string

const (
	JsonV1 = "v1"
	JsonV2 = "v2"
	JsonV3 = "v3"
)

func (v *Version) UnmarshalJSON(data []byte) error {
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	switch vt := value.(type) {
	case string:
		if !semver.IsValid(vt) {
			return errors.New("invalid version")
		}
		*v = Version(semver.Canonical(vt))
	case int:
		*v = Version(fmt.Sprintf("v%d.0.0", vt))
	case float64:
		*v = Version(fmt.Sprintf("v%d.0.0", int(vt)))
	default:
		return errors.New("invalid version type")
	}
	return nil
}

func (j Version) Compare(otherVersion Version) int {
	return semver.Compare(string(j), string(otherVersion))
}

func (j Version) Major() string {
	return semver.Major(string(j))
}
