package middlewares

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dheerajgopi/todo-api/common"
	todoErr "github.com/dheerajgopi/todo-api/common/error"
)

// JwtValidator middleware validates the token in the Authorization header.
// It responds with 403 error in case of invalid or missing token.
func JwtValidator(secret string) MiddlewareFunc {
	return func(f common.HandlerFunc) common.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request, reqCtx *common.RequestContext) (int, interface{}, *todoErr.APIError) {
			authHeader := req.Header["Authorization"]

			if authHeader == nil {
				err := todoErr.UnauthorizedError{}
				apiError := todoErr.NewAPIError(err.Error(), &todoErr.APIErrorBody{
					Message: "Access denied",
				})

				return 403, nil, apiError
			}

			token, err := jwt.Parse(authHeader[0], func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if !token.Valid || err != nil {
				err := todoErr.UnauthorizedError{}
				apiError := todoErr.NewAPIError(err.Error(), &todoErr.APIErrorBody{
					Message: "Access denied",
				})

				return 403, nil, apiError
			}

			claims, _ := token.Claims.(jwt.MapClaims)
			reqCtx.UserID = int64(claims["userId"].(float64))

			return f(res, req, reqCtx)
		}
	}
}
