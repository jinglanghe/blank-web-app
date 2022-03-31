package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func unauthorized(c *gin.Context, code int, err error) {
	var customErr Error
	_ = json.Unmarshal([]byte(err.Error()), &customErr)

	c.Abort()
	c.JSON(code, gin.H{
		"code": customErr.Code,
		"msg":  customErr.Message,
		"data": struct {}{},
	})
}

type AuthN struct {
	opts Options

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is TokenTime + MaxRefresh.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName string

	// TimeFunc provides the current time. You can override it to use another time value.
	// This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc func() time.Time

	// Public key
	pubKey *rsa.PublicKey

	// SendAuthorization allow return authorization header for every request
	SendAuthorization bool
}

func NewJwtAuthN(opts ...Option) *AuthN {
	authn := &AuthN{
		MaxRefresh:        time.Hour * 24 * 90,
		TokenLookup:       "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:     "Bearer",
		TimeFunc:          time.Now,
		SendAuthorization: true,
	}

	for _, o := range opts {
		o(&authn.opts)
	}

	if authn.usingPublicKeyAlgo() {
		_ = authn.readKeys()
	} else if authn.opts.Key == nil {
		return nil
	}
	return authn
}

func (a *AuthN) readKeys() error {
	err := a.publicKey()
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthN) publicKey() error {
	keyData, err := ioutil.ReadFile(a.opts.PubKeyFile)
	if err != nil {
		return ErrNoPubKeyFile
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPubKey
	}
	a.pubKey = key
	return nil
}

func (a *AuthN) usingPublicKeyAlgo() bool {
	switch a.opts.SigningAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
}

// GetClaimsFromJWT get claims from JWT token
func (a *AuthN) GetClaimsFromJWT(c *gin.Context) (jwt.MapClaims, error) {
	token, err := a.ParseToken(c)

	if err != nil {
		return nil, err
	}

	if a.SendAuthorization {
		if v, ok := c.Get("JWT_TOKEN"); ok {
			c.Header("Authorization", a.TokenHeadName+" "+v.(string))
		}
	}

	//toke
	claims := jwt.MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims, nil
}

func (a *AuthN) jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == a.TokenHeadName) {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}

func (a *AuthN) jwtFromQuery(c *gin.Context, key string) (string, error) {
	token := c.Query(key)

	if token == "" {
		return "", ErrEmptyQueryToken
	}

	return token, nil
}

func (a *AuthN) jwtFromCookie(c *gin.Context, key string) (string, error) {
	cookie, _ := c.Cookie(key)

	if cookie == "" {
		return "", ErrEmptyCookieToken
	}

	return cookie, nil
}

func (a *AuthN) jwtFromParam(c *gin.Context, key string) (string, error) {
	token := c.Param(key)

	if token == "" {
		return "", ErrEmptyParamToken
	}

	return token, nil
}

// ParseToken parse jwt token from gin context
func (a *AuthN) ParseToken(c *gin.Context) (*jwt.Token, error) {
	var token string
	var err error

	methods := strings.Split(a.TokenLookup, ",")
	for _, method := range methods {
		if len(token) > 0 {
			break
		}
		parts := strings.Split(strings.TrimSpace(method), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "header":
			token, err = a.jwtFromHeader(c, v)
		case "query":
			token, err = a.jwtFromQuery(c, v)
		case "cookie":
			token, err = a.jwtFromCookie(c, v)
		case "param":
			token, err = a.jwtFromParam(c, v)
		}
	}

	if err != nil {
		return nil, err
	}

	t, tErr := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(a.opts.SigningAlgorithm) != t.Method {
			return nil, errors.New("invalid signing algorithm!")
		}
		if a.usingPublicKeyAlgo() {
			return a.pubKey, nil
		}

		// save token string if valid
		c.Set("JWT_TOKEN", token)

		return a.opts.Key, nil
	})
	if tErr != nil {
		ErrTokenParseError.Message = tErr.Error()
		return nil, ErrTokenParseError
	}

	return t, nil
}

func (a *AuthN) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := a.GetClaimsFromJWT(c)
		if err != nil {
			unauthorized(c, http.StatusUnauthorized, err)
			return
		}

		if claims["exp"] == nil {
			unauthorized(c, http.StatusBadRequest, ErrMissingExpField)
			return
		}

		if _, ok := claims["exp"].(float64); !ok {
			unauthorized(c, http.StatusBadRequest, ErrWrongFormatOfExp)
			return
		}

		if int64(claims["exp"].(float64)) < a.TimeFunc().Unix() {
			unauthorized(c, http.StatusUnauthorized, ErrExpiredToken)
			return
		}

		c.Set("JWT_PAYLOAD", claims)

		c.Set(userIdCtxKey, int64(claims["user_id"].(float64)))
		c.Set(userNameCtxKey, claims["user_name"])
		c.Set(groupIdCtxKey, int64(claims["group_id"].(float64)))
		c.Set(groupAccountCtxKey, claims["group_account"])
		c.Set(orgIdCtxKey, int64(claims["organization_id"].(float64)))

		c.Next()
	}
}
