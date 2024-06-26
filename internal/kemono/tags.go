package kemono

import (
	"encoding/csv"
	"encoding/json"
	"strings"
)

type Tags []string

func (t *Tags) UnmarshalJSON(bytes []byte) error {
	var val string
	if err := json.Unmarshal(bytes, &val); err != nil {
		return err
	}
	if val == "" {
		*t = nil
		return nil
	}
	val = strings.TrimPrefix(val, "{")
	val = strings.TrimSuffix(val, "}")
	c := csv.NewReader(strings.NewReader(val))
	var err error
	*t, err = c.Read()
	return err
}
