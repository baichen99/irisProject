package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"irisProject/config"
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

var JWTConf = JWTConfig{
	ValidationKeyGetter: func(*Token) (interface{}, error) {
		key, err := jwt.ParseECPublicKeyFromPEM([]byte(config.Conf.JWT.PublicBytes))
		if err != nil {
			return nil, err
		}
		return key, nil
	},
	Extractor: FromAuthHeader,
	SigningMethod: SigningMethodES512,
}

// CheckJWTToken is a user authorization middleware
var CheckJWTToken = NewJWTMiddleware(JWTConf).Serve
