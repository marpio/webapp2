package middleware

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/gorilla/context"
	gor_handlers "github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/marpio/webapp2/server/httputils"
	"net/http"
	"time"
)

type MiddlewareContext interface {
	GetLogger() *log.Logger
	GetName() string
	GetCookieStore() *sessions.CookieStore
}

type middlewareCtx struct {
	MiddlewareContext
}

func CustomResponseWriterMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		crw := httputils.NewCustomResponseWriter(w)
		h.ServeHTTP(crw, r)
	})
}

func (ctx middlewareCtx) LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Try to get the real IP
		remoteAddr := r.RemoteAddr
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			remoteAddr = realIP
		}

		entry := ctx.MiddlewareContext.GetLogger().WithFields(log.Fields{
			"request": r.RequestURI,
			"method":  r.Method,
			"remote":  remoteAddr,
		})

		if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
			entry = entry.WithField("request_id", reqID)
		}
		entry.Info("started handling request")

		h.ServeHTTP(rw, r)

		latency := time.Since(start)
		res := rw.(httputils.CustomResponseWriter)
		ctx.GetLogger().WithFields(log.Fields{
			"status":      res.Status(),
			"text_status": http.StatusText(res.Status()),
			"size":        res.Size(),
			"took":        latency,
			fmt.Sprintf("measure#%s.latency", ctx.MiddlewareContext.GetName()): latency.Nanoseconds(),
		}).Info("completed handling request")
	})
}

func (ctx middlewareCtx) RecoveryMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)

				ctx.MiddlewareContext.GetLogger().WithFields(log.Fields{
					"error": err,
					"stack": errors.Wrap(err, 2).ErrorStack(),
				}).Info("Panic occured - recovering")
			}
		}()
		h.ServeHTTP(rw, r)
	})
}

func (ctx middlewareCtx) SetCookieStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		context.Set(req, "cookieStore", ctx.MiddlewareContext.GetCookieStore())

		next.ServeHTTP(res, req)
	})
}

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookieStore := context.Get(req, "cookieStore").(*sessions.CookieStore)
		session, _ := cookieStore.Get(req, "holzrepublic-session")
		user := session.Values["user"]

		if user == nil {
			http.Redirect(res, req, "/login", 302)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func NewDefaultMiddlewareChain(env MiddlewareContext) alice.Chain {
	var mCtx = &middlewareCtx{env}
	var chain = alice.New(
		mCtx.RecoveryMiddleware,
		CustomResponseWriterMiddleware,
		mCtx.LoggerMiddleware,
		mCtx.SetCookieStore,
		gor_handlers.CompressHandler)
	return chain
}
