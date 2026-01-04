package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(p []byte) (int, error) {
	if sr.status == 0 {
		sr.status = http.StatusOK
	}
	n, err := sr.ResponseWriter.Write(p)
	sr.bytes += n
	return n, err
}

func RequestLogger(base *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqId := r.Header.Get("X-Request-Id")
			if reqId == "" {
				reqId = newRequestId()
			}

			l := base.With(
				zap.String("requestId", reqId),
				zap.String("httpMethod", r.Method),
				zap.String("httpPath", r.URL.Path),
			)

			ctx := utils.WithLogger(r.Context(), l)
			r = r.WithContext(ctx)

			rec := &statusRecorder{ResponseWriter: w}
			start := time.Now()

			next.ServeHTTP(rec, r)

			l.Info("request",
				zap.Int("status", rec.status),
				zap.Int("bytes", rec.bytes),
				zap.Duration("duration", time.Since(start)),
				zap.String("remoteAddr", r.RemoteAddr),
				zap.String("userAgent", r.UserAgent()),
			)
		})
	}
}

func newRequestId() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
