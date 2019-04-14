package main

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	common "github.com/dheerajgopi/todo-api/common"
	_userHttpDelivery "github.com/dheerajgopi/todo-api/user/delivery/http"
	_userRepo "github.com/dheerajgopi/todo-api/user/repository"
	_userService "github.com/dheerajgopi/todo-api/user/service"
	"github.com/go-sql-driver/mysql"
)

func main() {
	timeout := 5 * time.Second

	dbConfig := mysql.NewConfig()
	dbConfig.User = "root"
	dbConfig.Passwd = "root"
	dbConfig.Addr = "0.0.0.0:3307"
	dbConfig.DBName = "todo"
	dbConfig.Net = "tcp"
	dbConfig.ParseTime = true

	// initialize logger
	logRotate := &lumberjack.Logger{
		Filename:   "application.log",
		MaxSize:    5,
		MaxBackups: 100,
		MaxAge:     10,
		Compress:   false,
	}

	logWriters := io.MultiWriter(os.Stdout, logRotate)
	defer logRotate.Close()

	logger := logrus.New()
	logger.SetOutput(logWriters)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// initialize DB
	dbConn, err := sql.Open("mysql", dbConfig.FormatDSN())
	err = dbConn.Ping()

	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		os.Exit(1)
	}

	defer dbConn.Close()

	config := &common.AppConfig{
		RequestTimeout: timeout,
	}

	app := &common.App{
		Config: config,
		Logger: logger,
	}

	router := mux.NewRouter()
	userRepo := _userRepo.New(dbConn)
	userService := _userService.New(userRepo)
	_userHttpDelivery.New(router, userService, app)

	logger.Info("Starting server at port 8080")
	logger.Fatal(http.ListenAndServe(":8080", router))
}
