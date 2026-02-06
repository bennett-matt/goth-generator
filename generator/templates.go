package generator

const baseTemplTemplate = `package templates

templ Base(title string, appName string, loggedIn bool) {
	<!DOCTYPE html>
	<html lang="en" data-theme="corporate">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>{ title }</title>
			<link href="/static/css/output.css" rel="stylesheet" type="text/css"/>
			<script src="/static/js/htmx.min.js" defer></script>
			<script>
				function setTheme(theme) {
					document.documentElement.setAttribute('data-theme', theme);
					localStorage.setItem('theme', theme);
				}
				document.addEventListener('DOMContentLoaded', function() {
					const saved = localStorage.getItem('theme');
					if (saved) document.documentElement.setAttribute('data-theme', saved);
				});
			</script>
		</head>
		<body hx-boost="true" class="min-h-screen bg-base-200">
			<div class="navbar bg-base-100 shadow-lg">
				<div class="flex-1">
					<a class="btn btn-ghost text-xl font-semibold" href="/">{ appName }</a>
				</div>
				<div class="flex-none gap-2">
					<div class="dropdown dropdown-end">
						<div tabindex="0" role="button" class="btn btn-ghost btn-sm gap-1">
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16"><path d="M8 11a3 3 0 1 1 0-6 3 3 0 0 1 0 6m0 1a4 4 0 1 0 0-8 4 4 0 0 0 0 8"/><path d="M8 0a1 1 0 0 1 1 1v1a1 1 0 0 1-2 0V1a1 1 0 0 1 1-1m0 4a1 1 0 0 1 1 1v1a1 1 0 0 1-2 0V5a1 1 0 0 1 1-1m0 4a1 1 0 0 1 1 1v1a1 1 0 0 1-2 0V9a1 1 0 0 1 1-1"/></svg>
							Theme
						</div>
						<ul tabindex="0" class="dropdown-content menu bg-base-100 rounded-box z-50 mt-2 w-40 p-2 shadow-lg">
							<li><button type="button" onclick="setTheme('light')">Light</button></li>
							<li><button type="button" onclick="setTheme('dark')">Dark</button></li>
							<li><button type="button" onclick="setTheme('corporate')">Corporate</button></li>
							<li><button type="button" onclick="setTheme('business')">Business</button></li>
							<li><button type="button" onclick="setTheme('cupcake')">Cupcake</button></li>
							<li><button type="button" onclick="setTheme('retro')">Retro</button></li>
							<li><button type="button" onclick="setTheme('forest')">Forest</button></li>
							<li><button type="button" onclick="setTheme('night')">Night</button></li>
						</ul>
					</div>
					<ul class="menu menu-horizontal px-1">
						<li><a href="/">Home</a></li>
						if loggedIn {
							<li><a href="/logout">Logout</a></li>
						} else {
							<li><a href="/login">Login</a></li>
						}
					</ul>
				</div>
			</div>
			<main class="container mx-auto p-6 max-w-6xl">
				{ children... }
			</main>
		</body>
	</html>
}
`

const homeTemplTemplate = `package templates

templ Home(loggedIn bool, userName string) {
	if loggedIn {
		<div class="space-y-6">
			<div class="hero bg-base-100 rounded-box">
				<div class="hero-content text-center">
					<div class="max-w-md">
						<h1 class="text-4xl font-bold">Welcome back, { userName }</h1>
						<p class="py-4">You're logged in. Get started building your application.</p>
						<div class="flex gap-4 justify-center flex-wrap">
							<a href="/users" class="btn btn-primary">View Users</a>
						</div>
					</div>
				</div>
			</div>
		</div>
	} else {
		<div class="space-y-16">
			<div class="hero min-h-[60vh] bg-base-100 rounded-box">
				<div class="hero-content text-center">
					<div class="max-w-2xl">
						<h1 class="text-5xl font-bold">Build something great</h1>
						<p class="py-6 text-lg opacity-80">A modern web application with authentication, sessions, and a clean foundation. Get started in minutes.</p>
						<div class="flex gap-4 justify-center flex-wrap">
							<a href="/register" class="btn btn-primary btn-lg">Get Started</a>
							<a href="/login" class="btn btn-outline btn-lg">Sign In</a>
						</div>
					</div>
				</div>
			</div>
			<div class="grid md:grid-cols-3 gap-6">
				<div class="card bg-base-100 shadow-xl">
					<div class="card-body">
						<h2 class="card-title text-primary">Authentication</h2>
						<p>Built-in user registration, login, and session management. Secure by default.</p>
					</div>
				</div>
				<div class="card bg-base-100 shadow-xl">
					<div class="card-body">
						<h2 class="card-title text-secondary">Modern Stack</h2>
						<p>Go, Tailwind, DaisyUI, HTMX, and Templ. Fast, simple, and maintainable.</p>
					</div>
				</div>
				<div class="card bg-base-100 shadow-xl">
					<div class="card-body">
						<h2 class="card-title text-accent">Database Ready</h2>
						<p>PostgreSQL or SQLite with migrations. SQLC for type-safe queries.</p>
					</div>
				</div>
			</div>
		</div>
	}
}
`

const loginTemplTemplate = `package templates

templ Login(errorMsg string, csrfToken string) {
	<div class="flex justify-center items-center min-h-[50vh]">
		<div class="card bg-base-100 shadow-xl w-full max-w-md">
			<div class="card-body">
				<h2 class="card-title text-2xl mb-4">Sign in</h2>
				if len(errorMsg) > 0 {
					<div class="alert alert-error mb-4">
						<span>{ errorMsg }</span>
					</div>
				}
				<form method="POST" action="/login" class="form-control gap-4">
					<input type="hidden" name="csrf_token" value={ csrfToken }/>
					<label class="form-control">
						<span class="label-text font-medium">Email</span>
						<input type="email" name="email" class="input input-bordered" placeholder="you@example.com" required />
					</label>
					<label class="form-control">
						<span class="label-text font-medium">Password</span>
						<input type="password" name="password" class="input input-bordered" required />
					</label>
					<button type="submit" class="btn btn-primary mt-2">Login</button>
				</form>
				<p class="text-sm text-center mt-4 opacity-70">Don't have an account? <a href="/register" class="link link-primary">Register</a></p>
			</div>
		</div>
	</div>
}
`

const registerTemplTemplate = `package templates

templ Register(errorMsg string, csrfToken string) {
	<div class="flex justify-center items-center min-h-[50vh]">
		<div class="card bg-base-100 shadow-xl w-full max-w-md">
			<div class="card-body">
				<h2 class="card-title text-2xl mb-4">Create account</h2>
				if len(errorMsg) > 0 {
					<div class="alert alert-error mb-4">
						<span>{ errorMsg }</span>
					</div>
				}
				<form method="POST" action="/register" class="form-control gap-4">
					<input type="hidden" name="csrf_token" value={ csrfToken }/>
					<label class="form-control">
						<span class="label-text font-medium">Name</span>
						<input type="text" name="name" class="input input-bordered" placeholder="Your name" required />
					</label>
					<label class="form-control">
						<span class="label-text font-medium">Email</span>
						<input type="email" name="email" class="input input-bordered" placeholder="you@example.com" required />
					</label>
					<label class="form-control">
						<span class="label-text font-medium">Password</span>
						<input type="password" name="password" class="input input-bordered" required />
					</label>
					<button type="submit" class="btn btn-primary mt-2">Register</button>
				</form>
				<p class="text-sm text-center mt-4 opacity-70">Already have an account? <a href="/login" class="link link-primary">Sign in</a></p>
			</div>
		</div>
	</div>
}
`

func (g *Generator) generateTemplates() error {
	// Base template
	basePath := g.projectPath("web/templates/base.templ")
	if err := g.writeFile(basePath, baseTemplTemplate); err != nil {
		return err
	}

	// Home template (always generated - landing page + dashboard)
	homePath := g.projectPath("web/templates/home.templ")
	if err := g.writeFile(homePath, homeTemplTemplate); err != nil {
		return err
	}

	// Login template
	if g.config.WithAuth {
		loginPath := g.projectPath("web/templates/login.templ")
		if err := g.writeFile(loginPath, loginTemplTemplate); err != nil {
			return err
		}

		// Register template
		registerPath := g.projectPath("web/templates/register.templ")
		if err := g.writeFile(registerPath, registerTemplTemplate); err != nil {
			return err
		}
	}

	return nil
}
