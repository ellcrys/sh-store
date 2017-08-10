package oauth

import (
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// GrantTypeClientCredentials represents a client credential grant type
	GrantTypeClientCredentials = "client_credentials"

	// GrantTypeAuthorizationCode represets an authorization grant type
	GrantTypeAuthorizationCode = "authorization_code"

	// SigningSecret is the secret used to create oauth tokens
	SigningSecret = os.Getenv("AUTH_SECRET")

	// TokenTypeApp is the name of an app token
	TokenTypeApp = "app"

	// TokenTypeSession is the name of a session token
	TokenTypeSession = "session"
)

// MakeToken creates an HMAC token
func MakeToken(secret string, claims map[string]interface{}) (string, error) {
	_claims := jwt.MapClaims(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, _claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
