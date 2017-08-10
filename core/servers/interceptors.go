package servers

import (
	"fmt"

	"google.golang.org/grpc"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ellcrys/elldb/core/servers/common"
	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/elldb/core/servers/oauth"
	"github.com/ellcrys/util"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

// CtxKey is used as key when defining string value in a context
type CtxKey string

var (
	// CtxAccount represents an authenticated account
	CtxAccount CtxKey = "account"

	// CtxContract represents an authenticated contract
	CtxContract CtxKey = "contract"

	// CtxTokenClaims represents claims in an auth token
	CtxTokenClaims CtxKey = "token_claims"

	// ErrInvalidToken represents a error about an invalid token
	ErrInvalidToken = fmt.Errorf("permission denied. Invalid token")

	// methodsNotRequiringAuth includes the full method name of methods that
	// must not be processed by the auth interceptor
	methodsNotRequiringAuth = []string{
		"/proto_rpc.API/Login",
		"/proto_rpc.API/GetAccount",
	}

	// methodsRequiringAppToken includes the full method name of methods that
	// require app token
	methodsRequiringAppToken = []string{
		"/proto_rpc.API/CreateBucket",
		"/proto_rpc.API/CreateObjects",
		"/proto_rpc.API/GetObjects",
		"/proto_rpc.API/CreateSession",
		"/proto_rpc.API/GetSession",
		"/proto_rpc.API/DeleteSession",
		"/proto_rpc.API/CountObjects",
		"/proto_rpc.API/CreateMapping",
		"/proto_rpc.API/GetMapping",
		"/proto_rpc.API/GetAllMapping",
		"/proto_rpc.API/CommitSession",
		"/proto_rpc.API/RollbackSession",
		"/proto_rpc.API/UpdateObjects",
		"/proto_rpc.API/DeleteObjects",
	}

	// methodsRequiringSessionToken includes the full method name of methods that
	// require session token
	methodsRequiringSessionToken = []string{
		"/proto_rpc.API/CreateContract",
	}
)

// Interceptors returns the API interceptors
func (s *RPC) Interceptors() grpc.UnaryServerInterceptor {
	return middleware.ChainUnaryServer(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logRPC.Debugf("New request [method=%s]", info.FullMethod)
		return handler(ctx, req)
	}, s.authInterceptor, s.requiresAppToken, s.requiresSessionToken)
}

// authInterceptor checks whether the request has valid access token
func (s *RPC) authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if util.InStringSlice(methodsNotRequiringAuth, info.FullMethod) {
		return handler(ctx, req)
	}

	tokenStr, err := util.GetAuthToken(ctx, "bearer")
	if err != nil {
		return nil, common.Error(401, common.CodeAuthorizationError, "", err.Error())
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		tokenType := token.Claims.(jwt.MapClaims)["type"]
		if tokenType == oauth.TokenTypeApp || tokenType == oauth.TokenTypeSession {
			return []byte(oauth.SigningSecret), nil
		}

		return nil, fmt.Errorf("unknown token type")
	})
	if err != nil {
		logRPC.Debugf("%+v", err)
		return nil, common.Error(401, common.CodeAuthorizationError, "", ErrInvalidToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	if err = claims.Valid(); err != nil {
		return nil, common.Error(401, common.CodeAuthorizationError, "", ErrInvalidToken.Error())
	}

	ctx = context.WithValue(ctx, CtxTokenClaims, claims)

	return handler(ctx, req)
}

// requiresAppToken checks if the claim is a valid app token.
// This will only apply to methods that require an app token. It immediately
// calls the next handler if the method does not require app token.
func (s *RPC) requiresAppToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var err error
	var contract db.Contract

	// verify claims for methods that require an app token
	if !util.InStringSlice(methodsRequiringAppToken, info.FullMethod) {
		return handler(ctx, req)
	}

	claims := ctx.Value(CtxTokenClaims).(jwt.MapClaims)
	if claims["type"] != oauth.TokenTypeApp {
		return nil, common.Error(401, "", "", "endpoint requires an app token")
	}

	// get contract with matching client id
	if err = s.db.Where("client_id = ?", claims["client_id"]).First(&contract).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx, common.Error(401, common.CodeAuthorizationError, "", ErrInvalidToken.Error())
		}
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	ctx = context.WithValue(ctx, CtxContract, &contract)
	return handler(ctx, req)
}

// requiresSessionToken checks if the claim is a valid session token.
// This will only apply to methods that require a session token. It immediately
// calls the next handler if not.
func (s *RPC) requiresSessionToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	var err error
	var account db.Account

	// verify claims for methods that require an app token
	if !util.InStringSlice(methodsRequiringSessionToken, info.FullMethod) {
		return handler(ctx, req)
	}

	claims := ctx.Value(CtxTokenClaims).(jwt.MapClaims)
	if claims["type"] != oauth.TokenTypeSession {
		return nil, common.Error(401, "", "", "endpoint requires a session token")
	}

	// get account
	err = s.db.Where("id = ?", claims["id"]).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx, common.Error(401, common.CodeAuthorizationError, "", ErrInvalidToken.Error())
		}
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	ctx = context.WithValue(ctx, CtxAccount, &account)

	return handler(ctx, req)
}
