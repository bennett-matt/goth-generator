package generator

const mainGoTemplate = `package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/nosurf"
	"{{.Module}}/internal/database"
	"{{.Module}}/internal/handlers"
	"{{.Module}}/internal/middleware"
	{{if .WithSessions}}"{{.Module}}/internal/session"{{end}}
	{{if .WithAuth}}"{{.Module}}/internal/user"{{end}}
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db, err := database.New()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	{{if .WithSessions}}
	sessionStore := session.NewStore(db)
	{{end}}
	{{if .WithAuth}}
	userService := user.NewService(db)
	{{end}}

	h := handlers.NewHandler("{{.Name}}", db{{if .WithSessions}}, sessionStore{{end}}{{if .WithAuth}}, userService{{end}})

	router := httprouter.New()

	// Health check
	router.GET("/health", handlers.HealthCheck)

	// Static files
	router.ServeFiles("/static/*filepath", http.Dir("web/static"))

	// Apply middleware
	handler := middleware.Logging(router)
	handler = middleware.Recovery(handler)
	{{if .WithSessions}}
	handler = middleware.Session(sessionStore, handler)
	{{end}}
	{{if .WithAuth}}
	handler = middleware.Auth(handler)
	{{end}}
	handler = nosurf.New(handler)

	// Routes - home page is always available
	router.GET("/", h.Home)
	{{if .WithAuth}}
	router.GET("/login", h.Login)
	router.POST("/login", h.HandleLogin)
	router.GET("/register", h.Register)
	router.POST("/register", h.HandleRegister)
	router.POST("/logout", h.HandleLogout)
	{{end}}

	{{if .WithUsers}}
	router.GET("/users", h.ListUsers)
	router.GET("/users/:id", h.GetUser)
	{{end}}

	port := os.Getenv("PORT")
	if port == "" {
		port = "{{.Port}}"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
`

func (g *Generator) generateMain() error {
	return g.writeTemplate(g.projectPath("cmd/server/main.go"), mainGoTemplate, g.config)
}
