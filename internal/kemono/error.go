package kemono

import (
	"encoding/json"
	"net/http"
)

func NewUpstreamResponseError(resp *http.Response) UpstreamResponseError {
	type errResponse struct {
		Error string `json:"error"`
	}
	var errText errResponse
	_ = json.NewDecoder(resp.Body).Decode(&errText)
	_ = resp.Body.Close()
	return UpstreamResponseError{
		Response:  resp,
		errorText: errText.Error,
	}
}

type UpstreamResponseError struct {
	Response  *http.Response
	errorText string
}

func (u UpstreamResponseError) Body() string {
	if len(u.errorText) == 0 {
		return u.Response.Status
	}
	return u.errorText
}

func (u UpstreamResponseError) Error() string {
	return "upstream request failed with status: " + u.Response.Status
}
