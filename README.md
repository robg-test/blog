# Bob Productions Blog

A personal blog built with Go, templ, and Tailwind CSS.

## Tech Stack

- **Backend**: Go 1.24 with Chi router
- **Templating**: [templ](https://templ.guide) - type-safe HTML components for Go
- **Styling**: Tailwind CSS v4 + [DaisyUI](https://daisyui.com) components
- **Database**: Turso (libsql/SQLite-compatible)
- **Session Management**: Redis-backed SCS sessions
- **Deployment**: Docker + GitHub Actions to Digital Ocean VPS

## Features

- Theme switching (retro/synthwave) with session persistence
- Blog posts written in templ components
- Performance testing guides and stoic philosophy posts
- Responsive design with DaisyUI components

## Project Structure

```
├── main.go                 # Entry point, routing, server setup
├── internal/
│   ├── database.go        # Database initialization
│   ├── static/           # Blog metadata and content definitions
│   └── models/          # Data models
├── web/
│   ├── blogs/            # Blog post templ components
│   ├── components/       # Reusable UI components
│   ├── pages/            # Page templates
│   └── static/          # CSS, JS, and images
└── Dockerfile            # Multi-stage Docker build
```

## Development

```bash
# Install dependencies
go mod download
npm install

# Generate templ files
go tool templ generate

# Build CSS
npx @tailwindcss/cli -i ./web/static/css/input.css -o ./web/static/css/output.css

# Run the server
ENV=development TURSO_DATABASE="file:blog.db" go run .
```

Server runs at `http://localhost:8080`

## Deployment

Pushed to `main` branch triggers GitHub Actions:
1. Builds Docker image
2. Pushes to GitHub Container Registry
3. Deploys to Digital Ocean VPS via SSH
4. Health check verifies deployment

## Blog Posts

- An Introduction To Bob Productions
- Software Performance Guide: AWS Lambdas
- Weekly Stoic: Control & Choice
- Weekly Stoic: To Be Steady & Unsteady
- The Code Suggestion Crisis
- AD Performance Workshop
- Grug Guide to Why Test Automation Fails

## License

MIT
