package generator

// Config holds the configuration for project generation
type Config struct {
	Name        string
	Module      string
	OutputDir   string
	DBDriver    string
	Port        string
	WithAuth    bool
	WithUsers   bool
	WithSessions bool
}
