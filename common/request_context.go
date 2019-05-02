package common

import (
	"github.com/sirupsen/logrus"
)

// RequestContext stores request scoped data
type RequestContext struct {
	RequestID string
	Response  *APIResponse
	LogEntry  *logrus.Entry
}

// AddLogFields will add the specified fields to the LogEntry
func (reqCtx *RequestContext) AddLogFields(fields logrus.Fields) {
	reqCtx.LogEntry = reqCtx.LogEntry.WithFields(fields)
}

// AddLogMessage will add the message to the LogEntry
func (reqCtx *RequestContext) AddLogMessage(message string) {
	reqCtx.LogEntry.Message = message
}

// LogInfo will write info level log
func (reqCtx *RequestContext) LogInfo() {
	reqCtx.LogEntry.Info(reqCtx.LogEntry.Message)
}

// LogWarn will write warn level log
func (reqCtx *RequestContext) LogWarn() {
	reqCtx.LogEntry.Warn(reqCtx.LogEntry.Message)
}

// LogError will write error level log
func (reqCtx *RequestContext) LogError() {
	reqCtx.LogEntry.Error(reqCtx.LogEntry.Message)
}
