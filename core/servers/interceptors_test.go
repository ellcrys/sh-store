package servers

import (
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/elldb/core/servers/oauth"
	"github.com/ellcrys/util"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	SigningSecret = "secret"
)

func TestAuthInterceptor(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("authInterceptor", t, func() {

			var noopHandler = func(context.Context, interface{}) (interface{}, error) {
				return nil, nil
			}

			Convey("Should call handle if function does not require authorization", func() {
				methodsNotRequiringAuth = append(methodsNotRequiringAuth, "/proto_rpc.API/Method")
				called := false
				rpc.authInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method"}, func(context.Context, interface{}) (interface{}, error) {
					called = true
					return nil, nil
				})
				So(called, ShouldBeTrue)
			})

			Convey("Should return error if metadata is not included in context", func() {
				_, err := rpc.authInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method2"}, noopHandler)
				So(err, ShouldNotBeNil)
				m, _ := util.JSONToMap(err.Error())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"detail": "no metadata in context",
					"status": "401",
					"code":   "auth_error",
				})
			})

			Convey("Should return error if authorization header has no value", func() {
				ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", ""))
				_, err := rpc.authInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method2"}, noopHandler)
				So(err, ShouldNotBeNil)
				m, _ := util.JSONToMap(err.Error())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"detail": "authorization format is invalid",
					"status": "401",
					"code":   "auth_error",
				})
			})

			Convey("Should return error if authorization header has invalid bearer", func() {
				ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "invalid format here"))
				_, err := rpc.authInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method2"}, noopHandler)
				So(err, ShouldNotBeNil)
				m, _ := util.JSONToMap(err.Error())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"detail": "request unauthenticated with bearer",
					"status": "401",
					"code":   "auth_error",
				})
			})

			Convey("Should return error if authorization header has invalid token", func() {
				ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer invalid"))
				_, err := rpc.authInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method2"}, noopHandler)
				So(err, ShouldNotBeNil)
				m, _ := util.JSONToMap(err.Error())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"detail": "permission denied. Invalid token",
					"status": "401",
					"code":   "auth_error",
				})
			})

			Convey("Should return error if token type is unknown", func() {

				token, _ := oauth.MakeToken(SigningSecret, map[string]interface{}{
					"type": "unknown",
				})

				ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer "+token))
				_, err := rpc.authInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method2"}, noopHandler)
				So(err, ShouldNotBeNil)
				m, _ := util.JSONToMap(err.Error())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"detail": "permission denied. Invalid token",
					"status": "401",
					"code":   "auth_error",
				})
			})

			Convey("Should return error if token type is expired", func() {

				token, _ := oauth.MakeToken(SigningSecret, map[string]interface{}{
					"type": oauth.TokenTypeApp,
					"exp":  time.Now().AddDate(0, 0, -2).Unix(),
				})

				ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer "+token))
				_, err := rpc.authInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method2"}, noopHandler)
				So(err, ShouldNotBeNil)
				m, _ := util.JSONToMap(err.Error())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"detail": "permission denied. Invalid token",
					"status": "401",
					"code":   "auth_error",
				})
			})

			Convey("Should successfully validate and include authorization claims in handler context", func() {

				var handlerCtx context.Context
				var handler = func(ctx context.Context, req interface{}) (interface{}, error) {
					handlerCtx = ctx
					return nil, nil
				}

				token, _ := oauth.MakeToken(SigningSecret, map[string]interface{}{
					"type": oauth.TokenTypeApp,
				})

				ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer "+token))
				_, err := rpc.authInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method2"}, handler)
				So(err, ShouldBeNil)
				So(handlerCtx, ShouldNotBeEmpty)
				claims := handlerCtx.Value(CtxTokenClaims).(jwt.MapClaims)
				So(claims, ShouldResemble, jwt.MapClaims{
					"type": "app",
				})
			})
		})
	})
}

func TestRequiresAppToken(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("requiresAppToken", t, func() {

			var noopHandler = func(context.Context, interface{}) (interface{}, error) {
				return nil, nil
			}

			Convey("Should call handle if function does not require app token", func() {
				called := false
				rpc.requiresAppToken(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method"}, func(context.Context, interface{}) (interface{}, error) {
					called = true
					return nil, nil
				})
				So(called, ShouldBeTrue)
			})

			Convey("Should return error if token type is not an app token", func() {
				methodsRequiringAppToken = append(methodsRequiringAppToken, "/proto_rpc.API/Method")
				ctx := context.WithValue(context.Background(), CtxTokenClaims, jwt.MapClaims{"type": "session"})
				_, err := rpc.requiresAppToken(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method"}, noopHandler)
				So(err, ShouldNotBeNil)
				m, _ := util.JSONToMap(err.Error())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"detail": "endpoint requires an app token",
					"status": "401",
				})
			})

			Convey("Should return error if client is not found", func() {
				methodsRequiringAppToken = append(methodsRequiringAppToken, "/proto_rpc.API/Method")
				ctx := context.WithValue(context.Background(), CtxTokenClaims, jwt.MapClaims{"type": oauth.TokenTypeApp, "client_id": "unknown"})
				_, err := rpc.requiresAppToken(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method"}, noopHandler)
				So(err, ShouldNotBeNil)
				m, _ := util.JSONToMap(err.Error())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"detail": "permission denied. Invalid token",
					"status": "401",
					"code":   "auth_error",
				})
			})

			Convey("Should successfully include authenticated client/contract in context passed to handler", func() {

				var contract = db.Contract{ID: util.UUID4(), ClientID: util.RandString(5), ClientSecret: util.RandString(5)}
				err := rpc.db.Create(&contract).Error
				So(err, ShouldBeNil)

				var handlerCtx context.Context
				var handler = func(ctx context.Context, req interface{}) (interface{}, error) {
					handlerCtx = ctx
					return nil, nil
				}

				methodsRequiringAppToken = append(methodsRequiringAppToken, "/proto_rpc.API/Method")
				ctx := context.WithValue(context.Background(), CtxTokenClaims, jwt.MapClaims{"type": oauth.TokenTypeApp, "client_id": contract.ClientID})
				_, err = rpc.requiresAppToken(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method"}, handler)
				So(err, ShouldBeNil)
				So(handlerCtx, ShouldNotBeEmpty)
				So(handlerCtx.Value(CtxContract), ShouldHaveSameTypeAs, &contract)
				So(handlerCtx.Value(CtxContract).(*db.Contract).ID, ShouldEqual, contract.ID)
			})

		})
	})
}

func TestRequiresSessionToken(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("requiresSessionToken", t, func() {

			var noopHandler = func(context.Context, interface{}) (interface{}, error) {
				return nil, nil
			}

			Convey("Should call handle if function does not require session token", func() {
				called := false
				rpc.requiresSessionToken(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method"}, func(context.Context, interface{}) (interface{}, error) {
					called = true
					return nil, nil
				})
				So(called, ShouldBeTrue)
			})

			Convey("Should return error if token type is not an session token", func() {
				methodsRequiringSessionToken = append(methodsRequiringSessionToken, "/proto_rpc.API/Method")
				ctx := context.WithValue(context.Background(), CtxTokenClaims, jwt.MapClaims{"type": oauth.TokenTypeApp})
				_, err := rpc.requiresSessionToken(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method"}, noopHandler)
				So(err, ShouldNotBeNil)
				m, _ := util.JSONToMap(err.Error())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"detail": "endpoint requires a session token",
					"status": "401",
				})
			})

			Convey("Should return error authenticated session account does not exist", func() {
				methodsRequiringSessionToken = append(methodsRequiringSessionToken, "/proto_rpc.API/Method")
				ctx := context.WithValue(context.Background(), CtxTokenClaims, jwt.MapClaims{"type": oauth.TokenTypeSession, "id": "unknown"})
				_, err := rpc.requiresSessionToken(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method"}, noopHandler)
				So(err, ShouldNotBeNil)
				m, _ := util.JSONToMap(err.Error())
				errs := m["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"detail": "permission denied. Invalid token",
					"status": "401",
					"code":   "auth_error",
				})
			})

			Convey("Should successfully include authenticated account in context passed to handler", func() {

				var account = db.Account{ID: util.UUID4()}
				err := rpc.db.Create(&account).Error
				So(err, ShouldBeNil)

				var handlerCtx context.Context
				var handler = func(ctx context.Context, req interface{}) (interface{}, error) {
					handlerCtx = ctx
					return nil, nil
				}

				methodsRequiringSessionToken = append(methodsRequiringSessionToken, "/proto_rpc.API/Method")
				ctx := context.WithValue(context.Background(), CtxTokenClaims, jwt.MapClaims{"type": oauth.TokenTypeSession, "id": account.ID})
				_, err = rpc.requiresSessionToken(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proto_rpc.API/Method"}, handler)
				So(err, ShouldBeNil)
				So(handlerCtx, ShouldNotBeEmpty)
				So(handlerCtx.Value(CtxAccount), ShouldHaveSameTypeAs, &account)
				So(handlerCtx.Value(CtxAccount).(*db.Account).ID, ShouldEqual, account.ID)
			})
		})
	})
}
