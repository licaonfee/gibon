package gibon_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/licaonfee/gibon"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "default")
}
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u, p, ok := r.BasicAuth(); ok {
			if u == "admin" && p == "admin" {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Access Denied", http.StatusForbidden)
	})
}

func redirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Redirect", http.StatusPermanentRedirect)
	})
}
func TestApplyMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		middlewares []gibon.Middleware
		request     func() *http.Request
		code        int
		response    string
	}{
		{
			name:        "Empty Pipeline",
			middlewares: nil,
			code:        http.StatusOK,
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				return req
			},
			response: "default",
		},
		{
			name:        "Authentication OK",
			middlewares: []gibon.Middleware{authMiddleware},
			code:        http.StatusOK,
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.SetBasicAuth("admin", "admin")
				return req
			},
			response: "default",
		},
		{
			name:        "Authentication Forbidden",
			middlewares: []gibon.Middleware{authMiddleware},
			code:        http.StatusForbidden,
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.SetBasicAuth("admin", "no valid")
				return req
			},
			response: "Access Denied\n",
		},
		{
			name:        "Rewrite Auth Forbidden",
			middlewares: []gibon.Middleware{authMiddleware, redirectMiddleware},
			code:        http.StatusForbidden,
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.SetBasicAuth("admin", "no valid")
				return req
			},
			response: "Access Denied\n",
		},
		{
			name:        "Rewrite Auth Success",
			middlewares: []gibon.Middleware{authMiddleware, redirectMiddleware},
			code:        http.StatusPermanentRedirect,
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.SetBasicAuth("admin", "admin")
				return req
			},
			response: "Redirect\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := tt.request()
			pipeline := gibon.New()
			for i := 0; i < len(tt.middlewares); i++ {
				pipeline = pipeline.With(tt.middlewares[i])
			}
			rr := httptest.NewRecorder()
			handler := pipeline.BuildFunc(defaultHandler)
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tt.code {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.code)
			}
			if rr.Body.String() != tt.response {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.response)
			}
		})
	}

}
