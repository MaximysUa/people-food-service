package mwlogger

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	logging "people-food-service/pkg/client/logger"
	"time"
)

func New(logger *logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		t1 := time.Now()
		logger.Debugln("logger middleware enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger.Debugf("\nmethod: %s\n path: %s\n remote_adr: %s\n user_agent: %s\n request_id: %s\n",
				r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), middleware.GetReqID(r.Context()))
			mw := middleware.NewWrapResponseWriter(w, r.ProtoMinor)

			defer func() {
				logger.Debugf("request completed.\n Status: %d, bytes: %d, duration: %s",
					mw.Status(), mw.BytesWritten(), time.Since(t1))
			}()
			next.ServeHTTP(mw, r)
		}
		return http.HandlerFunc(fn)
	}
}
