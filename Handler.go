package myerrors

import (
	"log/slog"
	"net/http"
	"errors"
	"encoding/json"
	"runtime/debug"
)

type ErrorResponse struct {
	Code      Code   `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId  omitempty"`
}

type Logger interface {
	Error(msg string, args ...any)
}

type HandlerConfig struct {
	Logger Logger
	Expose bool
}

func defaultLogger() Logger {
	return slog.Default()
}

func WriteError(w http.ResponseWriter , r *http.Request , err error , cfg HandlerConfig){
	logger := cfg.Logger
	if logger == nil{
    logger =  defaultLogger();
	}
	var ae *ApiError

	if !errors.As(err , &ae){
		ae = Internal("something went wrong")
		logger.Error("unhandled error" , "err" ,err , "path" , r.URL.Path)
	} else{
     logger.Error("request failed", "code", ae.Code, "op", ae.Op, "err", ae.Err, "path", r.URL.Path)
	}

	resp := ErrorResponse{Code: ae.Code, Message: ae.Message}
	if cfg.Expose && ae.Err != nil{
		resp.Message = ae.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ae.HTTPStatus())
	_ = json.NewEncoder(w).Encode(resp)
}

// RecoverMiddleware is the "global exception handler" for panics: wrap
// your router with this once, at the top, and no single handler panic
// can take the whole server down. Any recovered panic is turned into a
// standard 500 ErrorResponse instead of a dropped connection.
func RecoverMiddleware(cfg HandlerConfig) func(http.Handler) http.Handler {
	logger := cfg.Logger
	if logger == nil {
		logger = defaultLogger()
	}
 
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						"panic", rec,
						"path", r.URL.Path,
						"stack", string(debug.Stack()),
					)
					WriteError(w, r, Internal("internal server error"), cfg)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
 
// HandlerFunc is a handler that returns an error, letting business logic
// `return myerrors.NotFound(...)` instead of writing responses inline.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error
 
// Wrap adapts a HandlerFunc into a standard http.HandlerFunc, routing
// any returned error through WriteError. Use this at the route
// registration site:
//
//	mux.Handle("/users/{id}", myerrors.Wrap(getUser, cfg))
func WrapHandler(h HandlerFunc, cfg HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			WriteError(w, r, err, cfg)
		}
	}
}
 