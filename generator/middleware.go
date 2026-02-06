package generator

const middlewareGoTemplate = `package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	{{if .WithSessions}}"{{.Module}}/internal/session"{{end}}
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

{{if .WithSessions}}
func Session(sessionStore *session.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err == nil && sessionID != nil {
			sess, err := sessionStore.Get(r.Context(), sessionID.Value)
			if err == nil && sess != nil {
				ctx := context.WithValue(r.Context(), "session", sess)
				ctx = context.WithValue(ctx, "userID", sess.UserID)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}
{{end}}

{{if .WithAuth}}
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public routes
		publicRoutes := map[string]bool{
			"/":        true,
			"/health":  true,
			"/login":   true,
			"/register": true,
			"/logout":  true,
			"/static/": true,
		}

		for route := range publicRoutes {
			if r.URL.Path == route {
				next.ServeHTTP(w, r)
				return
			}
			// Prefix match for routes like /static/ (not for "/" alone)
			if len(route) > 1 && strings.HasSuffix(route, "/") && len(r.URL.Path) >= len(route) && r.URL.Path[:len(route)] == route {
				next.ServeHTTP(w, r)
				return
			}
		}

		userID := r.Context().Value("userID")
		if userID == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
{{end}}
`

func (g *Generator) generateMiddleware() error {
	if err := g.writeTemplate(g.projectPath("internal/middleware/middleware.go"), middlewareGoTemplate, g.config); err != nil {
		return err
	}
	return g.writeTemplate(g.projectPath("internal/middleware/middleware_test.go"), middlewareTestTemplate, g.config)
}

const middlewareTestTemplate = `package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})
	handler := Logging(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Error("Logging: next handler was not called")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Logging: got status %d", rec.Code)
	}
}

func TestRecovery(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	handler := Recovery(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Recovery: got status %d, want 500", rec.Code)
	}
}

{{if .WithAuth}}
func TestAuth_NoUserID_RedirectsToLogin(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := Auth(next)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("Auth: got status %d, want 303", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/login" {
		t.Errorf("Auth: Location = %q, want /login", loc)
	}
}

func TestAuth_PublicRoute_Allowed(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})
	handler := Auth(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Error("Auth: next handler was not called for public route")
	}
}
{{end}}
`
