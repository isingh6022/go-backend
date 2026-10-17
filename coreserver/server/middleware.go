package server

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/isingh6022/go-fullstack-ai/coreserver/types"
)

func getRecoverMiddleware(l *logger) types.Middleware {
	errHandler := func(w http.ResponseWriter, r *http.Request) {
		if rec := recover(); rec != nil {
			if rec == http.ErrAbortHandler {
				panic(rec)
			}
			message := types.DefaultPanicMessage
			if l.isDebugEnabled() {
				message = fmt.Sprint(rec)
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(message))

			l.Error("panic recovered",
				"error", rec,
				"method", r.Method,
				"path", r.URL.Path,
				"stack", string(debug.Stack()),
			)
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer errHandler(w, r)

			next.ServeHTTP(w, r)
		})
	}
}

func Chain(middlewares ...types.Middleware) types.Middleware {
	return func(h http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			if middlewares[i] == nil {
				continue
			}
			h = middlewares[i](h)
		}
		return h
	}
}
