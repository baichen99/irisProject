package middlewares

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
