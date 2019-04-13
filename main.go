package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	_taskHttpDelivery "github.com/dheerajgopi/todo-api/task/delivery/http"
	_taskRepo "github.com/dheerajgopi/todo-api/task/repository"
	_taskService "github.com/dheerajgopi/todo-api/task/service"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	timeout := 5 * time.Second

	dbConfig := mysql.NewConfig()
	dbConfig.User = "root"
	dbConfig.Passwd = "root"
	dbConfig.Addr = "0.0.0.0:3307"
	dbConfig.DBName = "todo"
	dbConfig.Net = "tcp"

	dbConn, err := sql.Open("mysql", dbConfig.FormatDSN())
	err = dbConn.Ping()

	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}

	defer dbConn.Close()

	router := mux.NewRouter()
	taskRepo := _taskRepo.New(dbConn)
	taskService := _taskService.New(taskRepo)
	_taskHttpDelivery.New(router, taskService, timeout)

	logrus.Info("Starting server at port 8080")
	logrus.Fatal(http.ListenAndServe(":8080", router))
}
