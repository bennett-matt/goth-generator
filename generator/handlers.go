package generator

const handlersGoTemplate = `package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/nosurf"
	"{{.Module}}/web/templates"
	{{if or .WithSessions .WithAuth}}"{{.Module}}/internal/session"{{end}}
	{{if or .WithAuth .WithUsers}}"{{.Module}}/internal/user"{{end}}
)

type Handler struct {
	AppName string
	{{if or .WithSessions .WithAuth}}SessionStore *session.Store{{end}}
	{{if or .WithAuth .WithUsers}}UserService *user.Service{{end}}
}

func NewHandler(appName string{{if or .WithSessions .WithAuth}}, sessionStore *session.Store{{end}}{{if or .WithAuth .WithUsers}}, userService *user.Service{{end}}) *Handler {
	return &Handler{
		AppName: appName,
		{{if or .WithSessions .WithAuth}}SessionStore: sessionStore,{{end}}
		{{if or .WithAuth .WithUsers}}UserService: userService,{{end}}
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	loggedIn := false
	userName := ""
	{{if .WithAuth}}
	if userID := r.Context().Value("userID"); userID != nil {
		if user, err := h.UserService.GetByID(r.Context(), int(userID.(int64))); err == nil {
			loggedIn = true
			userName = user.Name
		}
	}
	{{end}}
	ctx := templ.WithChildren(r.Context(), templates.Home(loggedIn, userName))
	templates.Base("Home", h.AppName, loggedIn).Render(ctx, w)
}

{{if .WithAuth}}
func (h *Handler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	errorMsg := r.URL.Query().Get("error")
	ctx := templ.WithChildren(r.Context(), templates.Login(errorMsg, nosurf.Token(r)))
	templates.Base("Login", h.AppName, false).Render(ctx, w)
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Invalid form"), http.StatusSeeOther)
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Email and password are required"), http.StatusSeeOther)
		return
	}
	user, err := h.UserService.GetByEmail(r.Context(), email)
	if err != nil {
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Invalid email or password"), http.StatusSeeOther)
		return
	}
	if err := h.UserService.VerifyPassword(user, password); err != nil {
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Invalid email or password"), http.StatusSeeOther)
		return
	}
	if h.SessionStore == nil {
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Sessions not configured"), http.StatusSeeOther)
		return
	}
	sess, err := h.SessionStore.Create(r.Context(), user.ID, time.Now().Add(7*24*time.Hour))
	if err != nil {
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Failed to create session"), http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Path:     "/",
		MaxAge:   7 * 24 * 3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	errorMsg := r.URL.Query().Get("error")
	ctx := templ.WithChildren(r.Context(), templates.Register(errorMsg, nosurf.Token(r)))
	templates.Base("Register", h.AppName, false).Render(ctx, w)
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/register?error="+url.QueryEscape("Invalid form"), http.StatusSeeOther)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	if name == "" || email == "" || password == "" {
		http.Redirect(w, r, "/register?error="+url.QueryEscape("All fields are required"), http.StatusSeeOther)
		return
	}
	if len(password) < 8 {
		http.Redirect(w, r, "/register?error="+url.QueryEscape("Password must be at least 8 characters"), http.StatusSeeOther)
		return
	}
	_, err := h.UserService.Create(r.Context(), email, password, name)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			http.Redirect(w, r, "/register?error="+url.QueryEscape("Email already registered"), http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/register?error="+url.QueryEscape("Registration failed"), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login?registered=1", http.StatusSeeOther)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if cookie, err := r.Cookie("session_id"); err == nil && cookie != nil && h.SessionStore != nil {
		h.SessionStore.Delete(r.Context(), cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
{{end}}

{{if .WithUsers}}
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users, err := h.UserService.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	u, err := h.UserService.GetByID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}
{{end}}
`

func (g *Generator) generateHandlers() error {
	return g.writeTemplate(g.projectPath("internal/handlers/handlers.go"), handlersGoTemplate, g.config)
}
