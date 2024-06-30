package kemono

import (
	"bytes"
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

type Time time.Time

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

		*d = Time(parsed.UTC())
		return nil
	case json.Number:
		val, err := t.Int64()
		if err != nil {
			return err
		}

		parsed := time.Unix(val, 0).UTC()
		*d = Time(parsed)
		return nil
	case nil:
		*d = Time(time.Time{}.UTC())
		return nil
	}
	return ErrInvalidTimeType
}
