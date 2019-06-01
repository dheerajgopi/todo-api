package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dheerajgopi/todo-api/common"
	todoErr "github.com/dheerajgopi/todo-api/common/error"
	"github.com/dheerajgopi/todo-api/common/middlewares"
	"github.com/dheerajgopi/todo-api/models"
	"github.com/dheerajgopi/todo-api/task"
	"github.com/gorilla/mux"
)

// TaskHandler represents HTTP handler for tasks
type TaskHandler struct {
	TaskService task.Service
	App         *common.App
}

// New creates new HTTP handler for task
func New(router *mux.Router, service task.Service, app *common.App) {
	handler := &TaskHandler{
		TaskService: service,
		App:         app,
	}

	supportedFields := make(map[string]common.FieldType)
	supportedFields["title"] = common.STRING
	supportedFields["created_at"] = common.UNIXTIME

	jwtMiddleware := middlewares.JwtValidator(app.Config.Auth.Jwt.Secret)
	paginator := middlewares.Paginator(supportedFields)

	router.HandleFunc("/tasks", app.CreateHandler(jwtMiddleware(handler.Create))).Methods("POST")
	router.HandleFunc("/tasks", app.CreateHandler(jwtMiddleware(paginator(handler.List)))).Methods("GET")
}

// Create will store new task
func (handler *TaskHandler) Create(res http.ResponseWriter, req *http.Request, reqCtx *common.RequestContext) (int, interface{}, *todoErr.APIError) {
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	var createTaskReqBody CreateTaskRequest
	err := decoder.Decode(&createTaskReqBody)

	if err != nil {
		apiError := todoErr.NewAPIError("", &todoErr.APIErrorBody{
			Message: "Invalid request body",
		})

		return http.StatusBadRequest, nil, apiError
	}

	validationErrors := createTaskReqBody.ValidateAndBuild()

	if len(validationErrors) > 0 {
		apiError := todoErr.NewAPIError("", validationErrors...)

		return http.StatusBadRequest, nil, apiError
	}

	now := time.Now()

	newTask := &models.Task{
		Title:       createTaskReqBody.Title,
		Description: createTaskReqBody.Description,
		CreatedBy: &models.User{
			ID: reqCtx.UserID,
		},
		IsComplete: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	switch err = handler.TaskService.Create(context.TODO(), newTask); err.(type) {
	case nil:
		break
	default:
		apiError := todoErr.NewAPIError(err.Error(), &todoErr.APIErrorBody{
			Message: "Internal server error",
		})

		return http.StatusInternalServerError, nil, apiError
	}

	taskData := &TaskData{
		ID:          newTask.ID,
		Title:       newTask.Title,
		Description: newTask.Description,
		IsComplete:  newTask.IsComplete,
		CreatedAt:   newTask.CreatedAt,
		UpdatedAt:   newTask.UpdatedAt,
	}

	responseData := &CreateTaskResponse{
		Task: taskData,
	}

	return http.StatusCreated, responseData, nil
}

// List will return all tasks
func (handler *TaskHandler) List(res http.ResponseWriter, req *http.Request, reqCtx *common.RequestContext) (int, interface{}, *todoErr.APIError) {
	// sortFields := [...]string{
	// 	"id",
	// 	"title",
	// 	"created_at"
	// }

	// reqCtx.Page.Sort = append(reqCtx.Page.Sort, &common.Sort{
	// 	Field:     "id",
	// 	Direction: "desc",
	// })

	taskList := make([]*TaskData, 0)

	tasksByUserID, err := handler.TaskService.List(context.TODO(), reqCtx.UserID, reqCtx.Page)

	switch err {
	case nil:
		break
	default:
		apiError := todoErr.NewAPIError(err.Error(), &todoErr.APIErrorBody{
			Message: "Internal server error",
		})

		return http.StatusInternalServerError, nil, apiError
	}

	for _, task := range tasksByUserID {
		taskList = append(taskList, &TaskData{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			IsComplete:  task.IsComplete,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		})
	}

	var pageLen int
	pageInfo := Cursor{}

	if pageLen = len(taskList); pageLen == 0 {
		pageInfo.HasNext = false
	} else {
		lastRow := taskList[pageLen-1]
		pageInfo.HasNext = true

		for _, sort := range reqCtx.Page.Sort {
			switch sort.Field {
			case "title":
				sort.LastVal = lastRow.Title
			case "id":
				sort.LastVal = strconv.FormatInt(lastRow.ID, 10)
			}
		}

		cursorJSON, _ := json.Marshal(reqCtx.Page.Sort)
		pageInfo.Next = base64.URLEncoding.EncodeToString(cursorJSON)
	}

	responseData := &ListTaskResponse{
		Tasks:    taskList,
		PageInfo: pageInfo,
	}

	return http.StatusOK, responseData, nil
}
