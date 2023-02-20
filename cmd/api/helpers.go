package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type response map[string]interface{}

func (app *application) sendResponse(w http.ResponseWriter, data response, statusCode int) error {
	response, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(response))
	return nil
}

func (app *application) sendInternalServerErrorResponse(w http.ResponseWriter) {
	response, _ := json.Marshal(response{"message": "An error occured, please try again"})

	w.Header().Set("Contetn-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(response))
}

func (app *application) convertToInt(value string, defaultValue int) int {
	val, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return val
}
