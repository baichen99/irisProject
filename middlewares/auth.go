package middlewares

// clone from https://github.com/iris-contrib/middleware/
import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"irisProject/config"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

const (
	//DefaultContextKey jwt
	DefaultContextKey = "jwt"
)

// JWTConfig is a struct for specifying configuration options for the jwt middleware.
type JWTConfig struct {
	// The function that will return the Key to validate the JWT.
	// It can be either a shared secret or a public key.
	// Default value: nil
	ValidationKeyGetter jwt.Keyfunc
	// The name of the property in the request where the user (&token) information
	// from the JWT will be stored.
	// Default value: "jwt"
	ContextKey string
	// The function that will be called when there's an error validating the token
	// Default value:
	ErrorHandler errorHandler
	// A boolean indicating if the credentials are required or not
	// Default value: false
	CredentialsOptional bool
	// A function that extracts the token from the request
	// Default: FromAuthHeader (i.e., from Authorization header as bearer token)
	Extractor TokenExtractor
	// When set, all requests with the OPTIONS method will use authentication
	// if you enable this option you should register your route with iris.Options(...) also
	// Default: false
	EnableAuthOnOptions bool
	// When set, the middelware verifies that tokens are signed with the specific signing algorithm
	// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
	// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
	// Default: nil
	SigningMethod jwt.SigningMethod
	// When set, the expiration time of token will be check every time
	// if the token was expired, expiration error will be returned
	// Default: false
	Expiration bool
}

type (
	// Token for JWT. Different fields will be used depending on whether you're
	// creating or parsing/verifying a token.
	//
	// A type alias for jwt.Token.
	Token = jwt.Token
	// MapClaims type that uses the map[string]interface{} for JSON decoding
	// This is the default claims type if you don't supply one
	//
	// A type alias for jwt.MapClaims.
	MapClaims = jwt.MapClaims
	// Claims must just have a Valid method that determines
	// if the token is invalid for any supported reason.
	//
	// A type alias for jwt.Claims.
	Claims = jwt.Claims
)

// Shortcuts to create a new Token.
var (
	NewToken           = jwt.New
	NewTokenWithClaims = jwt.NewWithClaims
)

// HS256 and company.
var (
	SigningMethodHS256 = jwt.SigningMethodHS256
	SigningMethodHS384 = jwt.SigningMethodHS384
	SigningMethodHS512 = jwt.SigningMethodHS512
)

// ECDSA - EC256 and company.
var (
	SigningMethodES256 = jwt.SigningMethodES256
	SigningMethodES384 = jwt.SigningMethodES384
	SigningMethodES512 = jwt.SigningMethodES512
)

// A function called whenever an error is encountered
type errorHandler func(context.Context, error)

// TokenExtractor is a function that takes a context as input and returns
// either a token or an error.  An error should only be returned if an attempt
// to specify a token was found, but the information was somehow incorrectly
// formed.  In the case where a token is simply not present, this should not
// be treated as an error.  An empty string should be returned in that case.
type TokenExtractor func(context.Context) (string, error)

// JWTMiddleware the middleware for JSON Web tokens authentication method
type JWTMiddleware struct {
	Config JWTConfig
}

// OnError is the default error handler.
// Use it to change the behavior for each error.
// See `JWTConfig.ErrorHandler`.
func OnError(ctx context.Context, err error) {
	if err == nil {
		return
	}

	ctx.StopExecution()
	ctx.StatusCode(iris.StatusUnauthorized)
	ctx.WriteString(err.Error())
}

// NewJWTMiddleware constructs a new Secure instance with supplied options.
func NewJWTMiddleware(cfg ...JWTConfig) *JWTMiddleware {

	var c JWTConfig
	if len(cfg) == 0 {
		c = JWTConfig{}
	} else {
		c = cfg[0]
	}

	if c.ContextKey == "" {
		c.ContextKey = DefaultContextKey
	}

	if c.ErrorHandler == nil {
		c.ErrorHandler = OnError
	}

	if c.Extractor == nil {
		c.Extractor = FromAuthHeader
	}

	return &JWTMiddleware{Config: c}
}

func logf(ctx iris.Context, format string, args ...interface{}) {
	ctx.Application().Logger().Debugf(format, args...)
}

// Get returns the user (&token) information for this client/request
func (m *JWTMiddleware) Get(ctx context.Context) (username, role string) {
	parsedToken := ctx.Values().Get(m.Config.ContextKey)
	if parsedToken == nil {
		ctx.StopExecution()
	} else {
		payload, _ := parsedToken.(*jwt.Token)
		params, _ := payload.Claims.(jwt.MapClaims)
		username = params["user"].(string)
		role = params["role"].(string)
	}
	return

}

// Serve the middleware's action
func (m *JWTMiddleware) Serve(ctx context.Context) {
	if err := m.CheckJWT(ctx); err != nil {
		m.Config.ErrorHandler(ctx, err)
		return
	}
	// If everything ok then call next.
	ctx.Next()
}

// FromAuthHeader is a "TokenExtractor" that takes a give context and extracts
// the JWT token from the Authorization header.
func FromAuthHeader(ctx context.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// FromParameter returns a function that extracts the token from the specified
// query string parameter
func FromParameter(param string) TokenExtractor {
	return func(ctx context.Context) (string, error) {
		return ctx.URLParam(param), nil
	}
}

// FromFirst returns a function that runs multiple token extractors and takes the
// first token it finds
func FromFirst(extractors ...TokenExtractor) TokenExtractor {
	return func(ctx context.Context) (string, error) {
		for _, ex := range extractors {
			token, err := ex(ctx)
			if err != nil {
				return "", err
			}
			if token != "" {
				return token, nil
			}
		}
		return "", nil
	}
}

var (
	// ErrTokenMissing is the error value that it's returned when
	// a token is not found based on the token extractor.
	ErrTokenMissing = errors.New("required authorization token not found")

	// ErrTokenInvalid is the error value that it's returned when
	// a token is not valid.
	ErrTokenInvalid = errors.New("token is invalid")

	// ErrTokenExpired is the error value that it's returned when
	// a token value is found and it's valid but it's expired.
	ErrTokenExpired = errors.New("token is expired")
)

// SignJWTToken is used for sign a JWT Token for clients
func SignJWTToken(userID uuid.UUID, role string) (string, error) {
	payload := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{
		"iss":  config.Conf.App.Name,                                                          	// Issuer
		"iat":  time.Now().Unix(),                                                              // Issued At
		"nbf":  time.Now().Unix(),                                                              // Not Before 		// JWT Token ID
		"exp":  time.Now().Add(time.Hour * time.Duration(config.Conf.JWT.ExpireHours)).Unix(), 	// Expiration Time
		"user": userID,                                                               			// Username 		// Role of User: Admin/User/etc...
		"role": role,																			// Role
	})
	key, _ := jwt.ParseECPrivateKeyFromPEM([]byte(config.Conf.JWT.PrivateBytes))
	token, err := payload.SignedString(key)
	return token, err
}

// CheckJWT the main functionality, checks for token
func (m *JWTMiddleware) CheckJWT(ctx context.Context) error {
	if !m.Config.EnableAuthOnOptions {
		if ctx.Method() == iris.MethodOptions {
			return nil
		}
	}

	// Use the specified token extractor to extract a token from the request
	token, err := m.Config.Extractor(ctx)

	// If debugging is turned on, log the outcome
	if err != nil {
		logf(ctx, "Error extracting JWT: %v", err)
		return err
	}

	logf(ctx, "Token extracted: %s", token)

	// If the token is empty...
	if token == "" {
		// Check if it was required
		if m.Config.CredentialsOptional {
			logf(ctx, "No credentials found (CredentialsOptional=true)")
			// No error, just no token (and that is ok given that CredentialsOptional is true)
			return nil
		}

		// If we get here, the required token is missing
		logf(ctx, "Error: No credentials found (CredentialsOptional=false)")
		return ErrTokenMissing
	}

	// Now parse the token

	parsedToken, err := jwt.Parse(token, m.Config.ValidationKeyGetter)
	// Check if there was an error in parsing...
	if err != nil {
		logf(ctx, "Error parsing token: %v", err)
		return err
	}

	if m.Config.SigningMethod != nil && m.Config.SigningMethod.Alg() != parsedToken.Header["alg"] {
		err := fmt.Errorf("Expected %s signing method but token specified %s",
			m.Config.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		logf(ctx, "Error validating token algorithm: %v", err)
		return err
	}

	// Check if the parsed token is valid...
	if !parsedToken.Valid {
		logf(ctx, "Token is invalid")
		m.Config.ErrorHandler(ctx, ErrTokenInvalid)
		return ErrTokenInvalid
	}

	if m.Config.Expiration {
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			if expired := claims.VerifyExpiresAt(time.Now().Unix(), true); !expired {
				logf(ctx, "Token is expired")
				return ErrTokenExpired
			}
		}
	}

	logf(ctx, "JWT: %v", parsedToken)

	// If we get here, everything worked and we can set the
	// user property in context.
	ctx.Values().Set(m.Config.ContextKey, parsedToken)

	return nil
}
