package common

import (
	"encoding/json"
	"fmt"
	"net/http"

	todoErr "github.com/dheerajgopi/todo-api/common/error"
	"github.com/dheerajgopi/todo-api/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// App stores app config
type App struct {
	Logger *logrus.Logger
	Config *config.Config
}

// CreateHandler creates a new HandlerFunc with a new RequestContext per request
func (app *App) CreateHandler(fn handlerFunc) func(res http.ResponseWriter, req *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		reqCtx := app.newRequestContext()

		// add log details
		reqCtx.AddLogFields(logrus.Fields{
			"requestId": reqCtx.RequestID,
			"request":   req.RequestURI,
			"method":    req.Method,
		})

		status, data, apiError := fn(res, req, reqCtx)

		reqCtx.AddLogFields(logrus.Fields{
			"status": status,
		})

		switch status {
		case http.StatusOK:
			reqCtx.Response.Status = 200
			reqCtx.Response.Data = data
			reqCtx.LogInfo()
		case http.StatusCreated:
			reqCtx.Response.Status = 201
			reqCtx.Response.Data = data
			reqCtx.LogInfo()
		case http.StatusBadRequest:
			reqCtx.Response.Status = 400
			reqCtx.Response.Errors = apiError.Body
			fmt.Println(reqCtx.LogEntry.Message)
			reqCtx.LogInfo()
		case http.StatusNotFound:
			reqCtx.Response.Status = 404
			reqCtx.Response.Errors = apiError.Body
			reqCtx.LogWarn()
		case http.StatusConflict:
			reqCtx.Response.Status = 409
			reqCtx.Response.Errors = apiError.Body
			reqCtx.LogEntry = reqCtx.LogEntry.WithError(apiError)
			reqCtx.LogWarn()
		default:
			reqCtx.Response.Status = 500
			reqCtx.Response.Errors = apiError.Body
			reqCtx.LogEntry = reqCtx.LogEntry.WithError(apiError)
			reqCtx.LogError()
		}

		response, _ := json.Marshal(reqCtx.Response)

		res.Header().Set("Content-Type", "application/json")
		res.Header().Set("X-Request-ID", reqCtx.RequestID)
		res.WriteHeader(reqCtx.Response.Status)
		res.Write(response)
	}
}

// NewRequestContext creates new struct to store request scoped data
func (app *App) newRequestContext() *RequestContext {
	requestID, _ := uuid.NewUUID()

	return &RequestContext{
		RequestID: requestID.String(),
		Response:  &APIResponse{},
		LogEntry: app.Logger.WithFields(
			logrus.Fields{},
		),
	}
}

// handler inserts the request scope
type handlerFunc func(http.ResponseWriter, *http.Request, *RequestContext) (int, interface{}, *todoErr.APIError)
