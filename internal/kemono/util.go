package kemono

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//nolint:gochecknoglobals
var serviceNameReplacer = strings.NewReplacer(
	"fans", "Fans",
	"star", "Star",
)

func formatServiceName(name string) string {
	caser := cases.Title(language.English)
	name = caser.String(name)
	name = serviceNameReplacer.Replace(name)
	return name
}

var ErrInvalidTimeType = errors.New("invalid time type")

type Time struct {
	time.Time
}

func (d *Time) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	t, err := decoder.Token()
	if err != nil {
		return err
	}

	switch t := t.(type) {
	case string:
		parsed, err := time.Parse("2006-01-02T15:04:05", t)
		if err != nil {
			return err
		}

		d.Time = parsed.UTC()
		return nil
	case json.Number:
		val, err := t.Int64()
		if err != nil {
			return err
		}

		d.Time = time.Unix(val, 0).UTC()
		return nil
	case nil:
		d.Time = time.Time{}.UTC()
		return nil
	}
	return ErrInvalidTimeType
}

type Tags []string

func (t *Tags) UnmarshalJSON(bytes []byte) error {
	var list []string
	if err := json.Unmarshal(bytes, &list); err == nil {
		*t = list
		return nil
	}

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
