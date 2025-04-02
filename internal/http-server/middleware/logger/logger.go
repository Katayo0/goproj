package logger

import (
	"log/slog"
	"github.com/gin-gonic/gin"
	"github.com/sumit-tembe/gin-requestid"
)

/*
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler{
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", requestid.GetRequestIDFromContext(r.Context())),
			)
			ww := gin.WrapH(w, r.ProtoMajor)
		}
	}
}
	*/

	func New(log *slog.Logger) gin.HandlerFunc {
		return func (c *gin.Context) {
			c.Next()
			log.Debug("Method: %s, Path: %s, Status: %d, requestID: %d",
					  c.Request.Method,
					  c.Request.URL.Path,
					  c.Request.Response.Status,
					  requestid.GetRequestIDFromContext(c),
					)
		}
	}