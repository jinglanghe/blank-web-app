package jwt

import "encoding/json"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (e *Error) Error() string {
	str, _ := json.Marshal(e)
	return string(str)
}

func ModuleError(baseCode int) int {
	return Module*100000 + baseCode
}

var (
	Module = 5000

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrExpiredToken = &Error{ModuleError(30007), "token is expired"}

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = &Error{ModuleError(30008), "auth header is empty"}

	// ErrMissingExpField missing exp field in token
	ErrMissingExpField = &Error{ModuleError(30009), "missing exp field"}

	// ErrWrongFormatOfExp field must be float64 format
	ErrWrongFormatOfExp = &Error{ModuleError(30010), "exp must be float64 format"}

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = &Error{ModuleError(30011), "auth header is invalid"}

	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	ErrEmptyQueryToken = &Error{ModuleError(30012), "query token is empty"}

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cookie is empty
	ErrEmptyCookieToken = &Error{ModuleError(30013), "cookie token is empty"}

	// ErrEmptyParamToken can be thrown if authing with parameter in path, the parameter in path is empty
	ErrEmptyParamToken = &Error{ModuleError(30014), "parameter token is empty"}

	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = &Error{ModuleError(30015), "invalid signing algorithm!!"}

	// ErrNoPubKeyFile indicates that the given public key is unreadable
	ErrNoPubKeyFile = &Error{ModuleError(30017), "public key file unreadable"}


	// ErrInvalidPubKey indicates the the given public key is invalid
	ErrInvalidPubKey = &Error{ModuleError(30019), "public key invalid"}

	// ErrTokenParseError indicates the the given token is invalid
	ErrTokenParseError = &Error{ModuleError(30021), "token parse error"}
)
