package main

import (
	"net/http"

	"github.com/sirupsen/logrus"

	_taskHttpDelivery "github.com/dheerajgopi/todo-api/task/delivery/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	_taskHttpDelivery.New(router)

	logrus.Info("Starting server at port 8080")
	logrus.Fatal(http.ListenAndServe(":8080", router))
}
