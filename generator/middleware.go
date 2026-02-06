package generator

const middlewareGoTemplate = `package middleware

import (
	"context"
	"log"
	"net/http"
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
			"/health":  true,
			"/login":   true,
			"/register": true,
			"/static/": true,
		}

		for route := range publicRoutes {
			if r.URL.Path == route || len(r.URL.Path) > len(route) && r.URL.Path[:len(route)] == route {
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
	return g.writeTemplate(g.projectPath("internal/middleware/middleware.go"), middlewareGoTemplate, g.config)
}
