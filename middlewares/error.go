package middlewares

import (
	"github.com/kataras/iris/v12"
)

const (
	message string = "message"
	errinfo string = "error"
)

const (
	// BadRequest 400 Error
	BadRequest string = "BadRequest"
	// Unauthorized 401 Error
	Unauthorized string = "Unauthorized"
	// Forbidden 403 Error
	Forbidden string = "Forbidden"
	// MethodError 405 Error
	MethodError string = "MethodError"
	// NotFound 404 Error
	NotFound string = "NotFound"
	// ConfictError 409 Error
	ConfictError string = "ConfictError"
	// UnprocessableEntity 422 Error
	UnprocessableEntity string = "UnprocessableEntity"
	// TooManyRequests 429 Error
	TooManyRequests string = "TooManyRequests"
	// InternalError 500 Error
	InternalError string = "InternalError"
	// URLNotFound 404 Error
	URLNotFound = "URLNotFound"
	// UnknownError other error code
	UnknownError string = "UnknownError"
)

// ErrorHandler handles request error
func ErrorHandler(ctx iris.Context) {
	errMessage := ctx.Values().GetString(errinfo)
	var msg string
	switch ctx.GetStatusCode() {
	case iris.StatusBadRequest: // 400 Error
		msg = BadRequest
	case iris.StatusUnauthorized: //401 Error
		msg = Unauthorized
	case iris.StatusForbidden: // 403 Error
		msg = Forbidden
	case iris.StatusNotFound: // 404 Error
		msg = NotFound
	case iris.StatusMethodNotAllowed: // 405 Error
		msg = MethodError
	case iris.StatusConflict: // 409 Error
		msg = ConfictError
	case iris.StatusUnprocessableEntity: // 422 Error
		msg = UnprocessableEntity
	case iris.StatusTooManyRequests: // 429 Error
		msg = TooManyRequests
	case iris.StatusInternalServerError: // 500 Error
		msg = InternalError
	default:
		msg = UnknownError
	}
	msg = ctx.Tr(msg)
	err := ctx.Tr(errMessage)
	if errMessage != "" {
		ctx.JSON(iris.Map{
			message: msg,
			errinfo: err,
		})
		return
	}
	ctx.JSON(iris.Map{
		message: ctx.Tr(URLNotFound),
	})
}
