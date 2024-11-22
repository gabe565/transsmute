// Code generated by "enumer -type Format -trimprefix Format -transform lower"; DO NOT EDIT.

package feed

import (
	"fmt"
	"strings"
)

const _FormatName = "unknownatomrssjson"

var _FormatIndex = [...]uint8{0, 7, 11, 14, 18}

const _FormatLowerName = "unknownatomrssjson"

func (i Format) String() string {
	if i >= Format(len(_FormatIndex)-1) {
		return fmt.Sprintf("Format(%d)", i)
	}
	return _FormatName[_FormatIndex[i]:_FormatIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _FormatNoOp() {
	var x [1]struct{}
	_ = x[FormatUnknown-(0)]
	_ = x[FormatAtom-(1)]
	_ = x[FormatRSS-(2)]
	_ = x[FormatJSON-(3)]
}

var _FormatValues = []Format{FormatUnknown, FormatAtom, FormatRSS, FormatJSON}

var _FormatNameToValueMap = map[string]Format{
	_FormatName[0:7]:        FormatUnknown,
	_FormatLowerName[0:7]:   FormatUnknown,
	_FormatName[7:11]:       FormatAtom,
	_FormatLowerName[7:11]:  FormatAtom,
	_FormatName[11:14]:      FormatRSS,
	_FormatLowerName[11:14]: FormatRSS,
	_FormatName[14:18]:      FormatJSON,
	_FormatLowerName[14:18]: FormatJSON,
}

var _FormatNames = []string{
	_FormatName[0:7],
	_FormatName[7:11],
	_FormatName[11:14],
	_FormatName[14:18],
}

// FormatString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func FormatString(s string) (Format, error) {
	if val, ok := _FormatNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _FormatNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Format values", s)
}

// FormatValues returns all values of the enum
func FormatValues() []Format {
	return _FormatValues
}

// FormatStrings returns a slice of all String values of the enum
func FormatStrings() []string {
	strs := make([]string, len(_FormatNames))
	copy(strs, _FormatNames)
	return strs
}

// IsAFormat returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Format) IsAFormat() bool {
	for _, v := range _FormatValues {
		if i == v {
			return true
		}
	}
	return false
}
