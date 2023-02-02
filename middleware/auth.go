package middleware

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/DmitryVesenniy/go-rest-framework/common"
	"github.com/DmitryVesenniy/go-rest-framework/config"
	"github.com/DmitryVesenniy/go-rest-framework/framework/resterrors"
)

var safePath []string = []string{
	"/auth",
}

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// skipping safe paths
		for _, p := range safePath {
			matched, err := regexp.Match(p, []byte(r.URL.Path))
			if err != nil {
				continue
			}
			if matched {
				next.ServeHTTP(w, r)
				return
			}
		}

		conf := config.Get()
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" { // Токен отсутствует, возвращаем  403 http-код Unauthorized
			permissionDeniedErr := &resterrors.UnauthorizedError{Err: "Missing auth token"}
			resterrors.RestErrorResponce(w, permissionDeniedErr)
			return
		}

		splitted := strings.Split(tokenHeader, " ") // Токен обычно поставляется в формате `Bearer {token-body}`, мы проверяем, соответствует ли полученный токен этому требованию
		if len(splitted) != 2 {
			permissionDeniedErr := &resterrors.UnauthorizedError{Err: "Invalid/malformed auth token"}
			resterrors.RestErrorResponce(w, permissionDeniedErr)
			return
		}
		tokenPart := splitted[1]

		var mySigningKey = []byte(conf.SecretKey)

		token, err := jwt.Parse(tokenPart, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			permissionDeniedErr := &resterrors.UnauthorizedError{Err: "Your Token has been expired"}
			resterrors.RestErrorResponce(w, permissionDeniedErr)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), common.ContextUserKey, claims)
			r = r.WithContext(ctx)

			if claims["is_superuser"] == true {
				r.Header.Set("User", "superuser")
				next.ServeHTTP(w, r)
				return
			} else {
				r.Header.Set("User", "user")
				next.ServeHTTP(w, r)
				return
			}
		}
		permissionDeniedErr := &resterrors.UnauthorizedError{Err: "Not authorized"}
		resterrors.RestErrorResponce(w, permissionDeniedErr)
	})
}
