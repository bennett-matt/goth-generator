package generator

const handlersGoTemplate = `package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/julienschmidt/httprouter"
	"{{.Module}}/internal/models"
	"{{.Module}}/web/templates"
	{{if .WithSessions}}"{{.Module}}/internal/session"{{end}}
	{{if .WithAuth}}"{{.Module}}/internal/user"{{end}}
)

type Handler struct {
	AppName string
	DB      *sql.DB
	{{if .WithSessions}}SessionStore *session.Store{{end}}
	{{if .WithAuth}}UserService *user.Service{{end}}
}

func NewHandler(appName string, db *sql.DB{{if .WithSessions}}, sessionStore *session.Store{{end}}{{if .WithAuth}}, userService *user.Service{{end}}) *Handler {
	return &Handler{
		AppName: appName,
		DB:      db,
		{{if .WithSessions}}SessionStore: sessionStore,{{end}}
		{{if .WithAuth}}UserService: userService,{{end}}
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
		if user, err := h.UserService.GetByID(r.Context(), userID.(int)); err == nil {
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
	ctx := templ.WithChildren(r.Context(), templates.Login())
	templates.Base("Login", h.AppName, false).Render(ctx, w)
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Implementation for login
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := templ.WithChildren(r.Context(), templates.Register())
	templates.Base("Register", h.AppName, false).Render(ctx, w)
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Implementation for registration
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Implementation for logout
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
{{end}}

{{if .WithUsers}}
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rows, err := h.DB.Query("SELECT id, email, name, created_at, updated_at FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
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

	var u models.User
	err = h.DB.QueryRow("SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
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
