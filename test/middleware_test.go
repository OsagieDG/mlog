package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OsagieDG/mlog/service/middleware"
)

func TestLogRequest(t *testing.T) {
	handler := middleware.LogRequest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}

	if res.Body.String() != "OK" {
		t.Fatalf("expected body 'OK', got '%s'", res.Body.String())
	}
}

func TestRecoverPanic(t *testing.T) {
	handler := middleware.RecoverPanic(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.Code)
	}

	if !strings.Contains(res.Body.String(), "Internal Server Error") {
		t.Fatalf("expected body to contain 'Internal Server Error', got '%s'", res.Body.String())
	}
}

func TestResponseLogger(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, "Created")
	})

	logger := &middleware.ResponseLogger{ResponseWriter: httptest.NewRecorder(), StatusCode: 0}

	req := httptest.NewRequest(http.MethodPost, "/create", nil)
	handler.ServeHTTP(logger, req)

	if logger.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", logger.StatusCode)
	}
}

func TestLogResponse(t *testing.T) {
	handler := middleware.LogResponse(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}

	if res.Body.String() != "Hello, World!" {
		t.Fatalf("expected body 'Hello, World!', got '%s'", res.Body.String())
	}
}

func TestLogStack(t *testing.T) {
	handler := middleware.MLog(
		middleware.RecoverPanic,
		middleware.LogRequest,
		middleware.LogResponse,
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Stack test"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/stack", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}

	if res.Body.String() != "Stack test" {
		t.Fatalf("expected body 'Stack test', got '%s'", res.Body.String())
	}
}
