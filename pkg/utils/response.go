package utils

import (
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type MessageResponse struct {
	Message string `json:"message"`
}

type CreatedResponse struct {
	Id int64 `json:"id"`
}

func ResponseInternalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(MessageResponse{
		Message: "Internal error!",
	})
	log.Printf(err.Error() + "\n" + string(debug.Stack()))
}

func ResponseMessage(w http.ResponseWriter, httpCode int, message string) {
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(MessageResponse{
		Message: message,
	})
}

func ResponseCreated(w http.ResponseWriter, id int64) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreatedResponse{
		Id: id,
	})
}

func Response(w http.ResponseWriter, httpCode int, data interface{}) {
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(data)
}