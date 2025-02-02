package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func MLog(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			next = xs[i](next)
		}
		return next
	}
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func LogResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &ResponseLogger{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		if !rw.wroteHeader {
			rw.WriteHeader(http.StatusOK)
		}

		log.Printf("%s - %s %s %s - Status: %d, Size: %d, Duration: %s\n",
			r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI(),
			rw.StatusCode, rw.Size, time.Since(start))
	})
}

type ResponseLogger struct {
	http.ResponseWriter
	StatusCode  int
	Size        int
	wroteHeader bool
}

func (rl *ResponseLogger) WriteHeader(statusCode int) {
	if rl.wroteHeader {
		return
	}
	rl.wroteHeader = true
	rl.StatusCode = statusCode
	rl.ResponseWriter.WriteHeader(statusCode)
}

func (rl *ResponseLogger) Write(data []byte) (int, error) {
	if !rl.wroteHeader {
		rl.WriteHeader(http.StatusOK)
	}
	size, err := rl.ResponseWriter.Write(data)
	rl.Size += size
	return size, err
}
