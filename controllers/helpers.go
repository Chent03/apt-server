package controllers

import (
	"encoding/json"
	"net/http"
)

func parseResponse(r *http.Request, dst interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&dst)
	if err != nil {
		return err
	}
	return nil
}
