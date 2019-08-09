package middlewares

import (
	"github.com/kataras/iris"
	"irisProject/config"

	"github.com/dgrijalva/jwt-go"
)

// middleware configs

var I18nConf = I18nConfig{
	Default:      "en-US",
	URLParameter: "lang",
	Languages: map[string]string{
		// 相对于main.go
		"en-US": "./locales/en-US.ini",
		"zh-CN": "./locales/zh-CN.ini",
	},
}

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

var GetJwtParams = NewJWTMiddleware(JWTConf).Get

func BeforeRequest(ctx iris.Context) {
	defer ctx.Next()
	ctx.Application().Logger().Info("=============== REQUEST START ===============")
}

func AfterRequest(ctx iris.Context) {
	defer ctx.Next()
	ctx.Application().Logger().Info("=============== REQUEST END ===============")
}