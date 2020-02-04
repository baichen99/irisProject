package middlewares

import (
	"irisProject/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
)

// middleware configs

// JWTConf is config for JWTMiddleware
var JWTConf = JWTConfig{
	ValidationKeyGetter: func(*Token) (interface{}, error) {
		key, err := jwt.ParseECPublicKeyFromPEM([]byte(config.Conf.JWT.PublicBytes))
		if err != nil {
			return nil, err
		}
		return key, nil
	},
	Extractor:     FromAuthHeader,
	SigningMethod: SigningMethodES512,
}

// CheckJWTToken is a user authorization middleware
var CheckJWTToken = NewJWTMiddleware(JWTConf).Serve

// CorsAllowAll is a cors middleware that allow all methods/origins
var CorsAllowAll = AllowAll()

var GetJWTParams = NewJWTMiddleware(JWTConf).Get
var CheckSuper = NewRoleMiddleware(roleConfig{
	Role: "Super",
}).Serve

func BeforeHandleRequest(ctx iris.Context) {
	defer ctx.Next()
	requestID := uuid.NewV4().String()
	ctx.Values().Set("RequestID", requestID)
	ctx.Application().Logger().Infof("=====================REQUEST " + requestID + " START====================")
}

// AfterHandleRequest is a global middleware after handle a request
func AfterHandleRequest(ctx iris.Context) {
	defer ctx.Next()
	requestID := ctx.Values().Get("RequestID").(string)
	ctx.Application().Logger().Infof("=====================REQUEST " + requestID + " END====================")
}
