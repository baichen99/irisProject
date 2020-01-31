package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"irisProject/config"
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

func BeforeRequest(ctx iris.Context) {
	defer ctx.Next()
	ctx.Application().Logger().Info("=============== REQUEST START ===============")
}

func AfterRequest(ctx iris.Context) {
	defer ctx.Next()
	ctx.Application().Logger().Info("=============== REQUEST END ===============")
}
