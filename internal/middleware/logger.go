package middleware

import (
	"github.com/invinciblewest/gophermart/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggerMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rd := &responseData{
			status: 200,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   rd,
		}

		next.ServeHTTP(&lw, r)

		logger.Log.Info("got incoming HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", rd.status),
			zap.Int("size", rd.size),
			zap.Duration("duration", time.Since(start)),
		)
	}

	return http.HandlerFunc(fn)
}
