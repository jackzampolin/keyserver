package api

import "encoding/json"

type restError struct {
	Error string `json:"error"`
}

func newError(err error) restError {
	return restError{err.Error()}
}

func (e restError) marshal() []byte {
	out, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return out
}
