package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

// PerformanceMiddleware logs request duration and response size
func PerformanceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}
		
		// Process request
		next.ServeHTTP(wrapped, r)
		
		// Calculate duration
		duration := time.Since(start)
		
		// Log performance metrics
		log.Printf("⏱️  [%s] %s - Status: %d - Duration: %v - Size: %d bytes",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
			wrapped.bytes,
		)
		
		// Warn if slow (>100ms)
		if duration > 100*time.Millisecond {
			log.Printf("⚠️  SLOW REQUEST: %s %s took %v", r.Method, r.URL.Path, duration)
		}
	})
}