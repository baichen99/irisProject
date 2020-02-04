package utils

import (
	"github.com/kataras/iris/v12"
	"gopkg.in/go-playground/validator.v9"
	"runtime"
)

func SetResponseError(context iris.Context, statusCode int, logPrefix string, err error) {
	_, fileName, lineNumber, ok := runtime.Caller(1)
	if !ok {
		fileName = "UNKNOWN"
		lineNumber = 0
	}
	context.Application().Logger().Errorf("[%s:%d] %s: %s", fileName, lineNumber, logPrefix, err.Error())
	context.StatusCode(statusCode)
	context.Values().Set("error", err.Error())
}

// ReadValidateForm reads form from JSON and validate it using validator.v9
func ReadValidateForm(c iris.Context, form interface{}) (err error) {
	err = c.ReadJSON(form)
	if err != nil {
		return
	}
	validate := validator.New()
	err = validate.Struct(form)
	return
}
