package kemono

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
)

type Creator struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Service string `json:"service"`
	Indexed uint   `json:"indexed"`
	Updated uint   `json:"updated"`
}

var ErrCreatorNotFound = errors.New("creator not found")

func getCreatorInfo(ctx context.Context, host, name, service string) (Creator, error) {
	var creator Creator

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+host+"/api/v1/creators", nil)
	if err != nil {
		return creator, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return creator, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	decoder := json.NewDecoder(resp.Body)

	if t, err := decoder.Token(); err != nil {
		return creator, err
	} else if t != json.Delim('[') {
		return creator, &json.UnmarshalTypeError{Value: "object", Type: reflect.TypeOf([]Creator{})}
	}

	for decoder.More() {
		if err := decoder.Decode(&creator); err != nil {
			return creator, err
		}

		if creator.Name == name && creator.Service == service {
			cancel()
			return creator, nil
		}
	}

	return creator, ErrCreatorNotFound
}
