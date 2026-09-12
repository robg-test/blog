# Bob Productions Blog

A personal blog built with Go, templ, and Tailwind CSS.

## Tech Stack

- **Backend**: Go 1.24 with Chi router
- **Templating**: [templ](https://templ.guide) - type-safe HTML components for Go
- **Styling**: Tailwind CSS v4 + [DaisyUI](https://daisyui.com) components
- **Database**: Turso (libsql/SQLite-compatible)
- **Session Management**: Redis-backed SCS sessions
- **Deployment**: None currently — the site is offline; CI builds the Docker image only

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

The blog is not currently deployed. The VPS it ran on has been decommissioned
and the deploy pipeline removed.

GitHub Actions (`.github/workflows/checks.yml`) runs on every push and pull
request and does two things only:

1. Verifies all images are WebP
2. Builds the Docker image

Nothing is published to a registry and nothing is deployed anywhere. The
`Dockerfile` is kept, and kept building, so the site can be revived on another
host without archaeology.

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
