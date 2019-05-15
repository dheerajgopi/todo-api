package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	common "github.com/dheerajgopi/todo-api/common"
	"github.com/dheerajgopi/todo-api/config"
	_userHttpDelivery "github.com/dheerajgopi/todo-api/user/delivery/http"
	_userRepo "github.com/dheerajgopi/todo-api/user/repository"
	_userService "github.com/dheerajgopi/todo-api/user/service"
	"github.com/go-sql-driver/mysql"
)

func main() {
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

	cfg := &config.Config{}
	if err := cfg.Load(); err != nil {
		logger.Errorf("Error loading config: %v", err)
		os.Exit(1)
	}

	dbConfig := mysql.NewConfig()
	dbConfig.User = cfg.Database.User
	dbConfig.Passwd = cfg.Database.Password
	dbConfig.Addr = cfg.Database.Address
	dbConfig.DBName = cfg.Database.Name
	dbConfig.Net = "tcp"
	dbConfig.ParseTime = true

	// initialize DB
	dbConn, err := sql.Open("mysql", dbConfig.FormatDSN())
	err = dbConn.Ping()

	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		os.Exit(1)
	}

	defer dbConn.Close()

	app := &common.App{
		Config: cfg,
		Logger: logger,
	}

	cfgJSON, _ := json.Marshal(app.Config)

	fmt.Println(string(cfgJSON))

	router := mux.NewRouter()
	userRepo := _userRepo.New(dbConn)
	userService := _userService.New(userRepo)
	_userHttpDelivery.New(router, userService, app)

	port := strconv.Itoa(cfg.Application.Port)
	logger.Info(fmt.Sprintf("Starting server at port %s", port))
	logger.Fatal(http.ListenAndServe(":"+port, router))
}
