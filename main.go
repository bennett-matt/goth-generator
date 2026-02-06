package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/bennett-matt/goth-generator/generator"
)

func main() {
	var (
		version     = flag.Bool("version", false, "Print version and exit")
		name        = flag.String("name", "", "Project name (required)")
		module      = flag.String("module", "", "Go module path (e.g., github.com/user/project)")
		output      = flag.String("output", ".", "Output directory for generated project")
		dbDriver    = flag.String("db", "postgres", "Database driver (postgres or sqlite)")
		port        = flag.String("port", "8080", "Server port")
		withAuth    = flag.Bool("auth", true, "Include authentication")
		withUsers   = flag.Bool("users", true, "Include user management")
		withSessions = flag.Bool("sessions", true, "Include session management")
	)
	flag.Parse()

	if *version {
		fmt.Println("goth-generate 1.1.0")
		return
	}

	if *name == "" {
		fmt.Fprintf(os.Stderr, "Error: -name is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *module == "" {
		*module = strings.ToLower(*name)
	}
	*module = normalizeModulePath(*module)

	config := &generator.Config{
		Name:         *name,
		Module:       *module,
		OutputDir:    *output,
		DBDriver:     *dbDriver,
		Port:         *port,
		WithAuth:     *withAuth,
		WithUsers:    *withUsers,
		WithSessions: *withSessions,
		GoVersion:    goVersionMinor(),
	}

	gen := generator.New(config)
	if err := gen.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully generated project '%s' in %s\n", *name, *output)
	fmt.Printf("📦 Module: %s\n", *module)
	fmt.Printf("🚀 Next steps:\n")
	fmt.Printf("   1. cd %s\n", *output)
	fmt.Printf("   2. make setup\n")
	fmt.Printf("   3. start db (docker-compose up -d db)\n")
	fmt.Printf("   4. make dev\n")
}

// normalizeModulePath strips URL schemes (https://, http://) from module paths.
// Go expects github.com/user/repo, not https://github.com/user/repo.
func normalizeModulePath(module string) string {
	for _, prefix := range []string{"https://", "http://"} {
		if strings.HasPrefix(module, prefix) {
			return strings.TrimPrefix(module, prefix)
		}
	}
	return module
}

// goVersionMinor returns the current Go version as major.minor (e.g. "1.24")
// for use in generated Dockerfiles and docs. Falls back to "1.23" if detection fails.
func goVersionMinor() string {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return "1.23"
	}
	// "go version go1.24.3 linux/amd64" -> "1.24"
	re := regexp.MustCompile(`go(\d+\.\d+)`)
	if m := re.FindStringSubmatch(string(out)); len(m) >= 2 {
		return m[1]
	}
	return "1.23"
}
