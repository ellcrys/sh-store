package oauth

import (
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/jinzhu/gorm"
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

// OAuth implements all the OAuth2 requirements of the application
type OAuth struct {
	db *gorm.DB
}

// NewOAuth creates a new OAuth instance
func NewOAuth(db *gorm.DB) *OAuth {
	return &OAuth{
		db: db,
	}
}

// GetToken returns a token
func (o *OAuth) GetToken(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var grantType = r.URL.Query().Get("grant_type")
	switch grantType {
	case GrantTypeClientCredentials:
		return o.getAppToken(w, r)
	case "":
		return common.NewSingleAPIErr(400, "", "", "Grant type is required", nil), 400
	case GrantTypeAuthorizationCode:
		return common.NewSingleAPIErr(400, "", "", "Not implemented", nil), 400
	default:
		return common.NewSingleAPIErr(400, "", "", "Invalid grant type", nil), 400
	}
}

// getAppToken process a token request using the `client_credentials` grant type
// and returns an app token that never expires.
func (o *OAuth) getAppToken(w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var clientID = r.URL.Query().Get("client_id")
	var clientSecret = r.URL.Query().Get("client_secret")

	if clientID == "" {
		return common.NewSingleAPIErr(400, common.CodeInvalidParam, "", "Client id is required", nil), 400
	}
	if clientSecret == "" {
		return common.NewSingleAPIErr(400, common.CodeInvalidParam, "", "Client secret is required", nil), 400
	}

	// get the account
	var account db.Account
	err := o.db.Where("client_id = ?", clientID).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return common.NewSingleAPIErr(400, common.CodeInvalidParam, "", "Client id or secret are invalid", nil), 401
		}
		return common.ServerError, 500
	}

	// check client secret
	if account.ClientSecret != clientSecret {
		return common.NewSingleAPIErr(400, common.CodeInvalidParam, "", "Client id or secret are invalid", nil), 401
	}

	token, _ := MakeToken(SigningSecret, map[string]interface{}{
		"id":   clientID,
		"type": TokenTypeApp,
		"iat":  time.Now().Unix(),
	})

	return map[string]string{
		"access_token": token,
	}, 201
}
