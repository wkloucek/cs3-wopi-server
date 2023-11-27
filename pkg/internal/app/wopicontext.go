package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

type key int

const (
	wopiContextKey key = iota
)

type WopiContext struct {
	AccessToken   string
	FileReference providerv1beta1.Reference
	User          *userv1beta1.User
	ViewMode      appproviderv1beta1.OpenInAppRequest_ViewMode
	EditAppUrl    string
	ViewAppUrl    string
}

func WopiContextAuthMiddleware(app *demoApp, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.URL.Query().Get("access_token")
		if accessToken == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(app.Config.WopiSecret), nil
		})

		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if err := claims.Valid(); err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		wopiContextAccessToken, err := DecryptAES([]byte(app.Config.WopiSecret), claims.WopiContext.AccessToken)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		claims.WopiContext.AccessToken = wopiContextAccessToken

		ctx = context.WithValue(ctx, wopiContextKey, claims.WopiContext)
		// authentication for the CS3 api
		ctx = metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, claims.WopiContext.AccessToken)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WopiContextFromCtx(ctx context.Context) (WopiContext, error) {
	if wopiContext, ok := ctx.Value(wopiContextKey).(WopiContext); ok {
		return wopiContext, nil
	}
	return WopiContext{}, errors.New("no wopi context found")
}
