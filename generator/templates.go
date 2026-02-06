package generator

const baseTemplTemplate = `package templates

import "github.com/a-h/templ"

templ Base(title string) {
	<!DOCTYPE html>
	<html lang="en" data-theme="light">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>{ title }</title>
			<link href="/static/css/output.css" rel="stylesheet" type="text/css"/>
			<script src="/static/js/htmx.min.js" defer></script>
		</head>
		<body hx-boost="true">
			<div class="navbar bg-base-100 shadow-lg">
				<div class="flex-1">
					<a class="btn btn-ghost text-xl" href="/">` + "{{.Name}}" + `</a>
				</div>
				<div class="flex-none">
					<ul class="menu menu-horizontal px-1">
						<li><a href="/">Home</a></li>
						<li><a href="/login">Login</a></li>
					</ul>
				</div>
			</div>
			<main class="container mx-auto p-4">
				{ children... }
			</main>
		</body>
	</html>
}
`

const loginTemplTemplate = `package templates

import "github.com/a-h/templ"

templ Login() {
	<div class="card bg-base-100 shadow-xl max-w-md mx-auto">
		<div class="card-body">
			<h2 class="card-title">Login</h2>
			<form method="POST" action="/login" class="form-control" hx-post="/login" hx-swap="innerHTML" hx-target="closest .card">
				<label class="label">
					<span class="label-text">Email</span>
				</label>
				<input type="email" name="email" class="input input-bordered" required />
				<label class="label">
					<span class="label-text">Password</span>
				</label>
				<input type="password" name="password" class="input input-bordered" required />
				<button type="submit" class="btn btn-primary mt-4">Login</button>
			</form>
		</div>
	</div>
}
`

const registerTemplTemplate = `package templates

import "github.com/a-h/templ"

templ Register() {
	<div class="card bg-base-100 shadow-xl max-w-md mx-auto">
		<div class="card-body">
			<h2 class="card-title">Register</h2>
			<form method="POST" action="/register" class="form-control" hx-post="/register" hx-swap="innerHTML" hx-target="closest .card">
				<label class="label">
					<span class="label-text">Name</span>
				</label>
				<input type="text" name="name" class="input input-bordered" required />
				<label class="label">
					<span class="label-text">Email</span>
				</label>
				<input type="email" name="email" class="input input-bordered" required />
				<label class="label">
					<span class="label-text">Password</span>
				</label>
				<input type="password" name="password" class="input input-bordered" required />
				<button type="submit" class="btn btn-primary mt-4">Register</button>
			</form>
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
