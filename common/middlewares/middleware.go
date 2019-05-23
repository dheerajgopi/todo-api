package middlewares

import (
	"github.com/dheerajgopi/todo-api/common"
)

// MiddlewareFunc is a middleware adapter to allow use of ordinary functions
// as middlewares.
type MiddlewareFunc func(common.HandlerFunc) common.HandlerFunc
