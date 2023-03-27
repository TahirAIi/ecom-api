package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
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

func (app *application) convertToInt(value string) (int, error) {
	val, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (app *application) uploadFile(key string, r *http.Request) (*string, error) {
	file, fileHeader, err := r.FormFile("main_picture")

	if err == nil {
		defer file.Close()
		err = os.Mkdir("../../uploads", os.ModePerm)
		if err != nil {
			return nil, err
		}

		fileName := strconv.Itoa(int(time.Now().UnixNano())) + fileHeader.Filename
		destinationFile, err := os.Create(fmt.Sprintf("../uploads/%s", fileName))
		if err != nil {
			return nil, err
		}
		defer destinationFile.Close()

		_, err = io.Copy(destinationFile, file)
		if err != nil {
			return nil, err
		}
		return &fileName, nil
	}
	return nil, err
}

func (app *application) GenerateFileUrl(fileName string) string {
	return os.Getenv("FILES_BASE_URL") + "/files/" + fileName
}
