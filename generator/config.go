package generator

// Config holds the configuration for project generation
type Config struct {
	Name         string
	Module       string
	OutputDir    string
	DBDriver     string
	Port         string
	WithAuth     bool
	WithUsers    bool
	WithSessions bool
	GoVersion    string // e.g. "1.24" - populated from `go version` at generation time
}
