package middleware

import (
	"net/http"
	"runtime"

	"github.com/DmitryVesenniy/go-rest-framework/framework/resterrors"
)

func PanicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				// framework.APP.LogService.Error(string(buf))

				resterrors.RestErrorResponce(w, resterrors.AppErr)
				return
			}
		}()

		h.ServeHTTP(w, r)
	})
}
