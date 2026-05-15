# AGENTS.md — Bob Productions Blog

## Project Overview

Personal blog built with **Go 1.24 + Chi router v5**, server-rendered HTML via **templ** components, styled with **Tailwind CSS v4 + DaisyUI v5**. Deployed via Docker to Digital Ocean, CI/CD on GitHub Actions.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.24 (`github.com/robgtest/blog`) |
| Router | Chi v5 (replaced Gorilla Mux) |
| HTML | templ v0.3.898 (type-safe Go components) |
| Styling | Tailwind CSS v4 + DaisyUI v5 + @tailwindcss/typography |
| Database | Turso (libsql, SQLite-compatible) |
| Sessions | Redis-backed SCS v2 |
| Interactivity | htmx 1.9.11 (from CDN) |
| Hosting | Digital Ocean VPS, HTTPS on :443 |
| CI/CD | GitHub Actions → Docker → GHCR → SSH deploy |

## Commands

```bash
# Development (with live reload)
make dev                  # tailwind watch + air

# Build & run once
make run                  # templ generate + go build + ./blog

# Manual steps
go mod download           # Go dependencies
npm install               # Tailwind + DaisyUI deps
go tool templ generate    # Generate _templ.go from .templ files
npx @tailwindcss/cli -i ./web/static/css/input.css -o ./web/static/css/output.css
ENV=development TURSO_DATABASE="file:blog.db" go run .
```

Air watches `.go`, `.tpl`, `.tmpl`, `.templ`, `.html` files and rebuilds automatically.

## Project Structure

```
main.go                   # Entry point, routing (all routes defined here)
internal/
  database.go             # Turso/libsql DB init (global DB *sql.DB)
  cache.go                # Redis session manager + GetMessage/PutMessage/ClearSession
  models/blog_meta.go     # BlogMeta struct (Title, Url, Description, ImageUri, Published)
  static/blog_data.go     # 8 BlogMeta vars with post metadata + generated URIs
web/
  pages/
    index.templ           # Homepage with hero, blog card grid, nav, footer
    blog-page.templ       # Blog post wrapper with breadcrumbs, PrismJS, copy button
  blogs/
    intro_blog.templ      # Each post = one .templ file + generated _templ.go
    aws_serverless.templ
    opinion_on_copilot.templ
    performance_workshop.templ
    quiet_skills.templ
    automation_fail.templ # NEWEST: "Grug Guide to Why Test Automation Fails"
    stoicism/             # Stoicism-themed posts in subdirectory
      control_and_choice.templ
      to_be_steady.templ
  components/
    navigation.templ      # Top navbar (hamburger menu + theme toggle + htmx)
    footer.templ          # Fox SVG logo + LinkedIn/GitHub links
    card.templ            # Blog post card (image, badge, read time)
    blog-library.templ    # Grid rendering all post cards
    blog-badge.templ      # Colored category badge
    blog-time-caption.templ # Clock icon + read time
    breadcrumb.templ      # Home > Blog breadcrumb
    blog-meta.templ       # OG + Twitter meta tags
    code-block.templ      # Syntax-highlighted code block
    multi-line-code-block.templ # Code block with copy button
    article-bottom.templ  # Stub for article footer
  static/
    css/input.css         # Tailwind source (12 lines: imports Tailwind + DaisyUI + Typography)
    css/output.css        # Compiled (gitignored)
    css/prism.css         # PrismJS theme
    js/prism.js           # PrismJS for syntax highlighting
    images/               # All WebP images, organised by topic subdirectories
```

## Key Patterns & Conventions

### Blog Post Creation
Every new blog post requires 4 changes:

1. **Create the templ component** at `web/blogs/<name>.templ`:
   - Package is `blogs` (or `stoicism` for stoicism posts)
   - Define `meta<Name>()` private component that renders `components.MetaData(static.<Name>Data)`
   - Define exported component `Blog<Name>(theme string)` that wraps content in `@pages.BlogPage(theme, meta<Name>()) { ... }`
   - Content goes inside the `BlogPage` block, styled with Tailwind/`prose`

2. **Add metadata** to `internal/static/blog_data.go`:
   - Create a `BlogMeta` var: `<Name>Data = models.BlogMeta{Title, Url, Description, ImageUri, Published}`
   - `Url` uses the dynamic `uri` prefix + `blog/<slug>`
   - Use the `date()` helper for `Published` field

3. **Register the route** in `main.go:setupBlogHandler()`:
   - Add a `case "<slug>":` that maps to the new blog component

4. **Add to RSS feed** in `main.go:serveRSS()`:
   - Include the `BlogMeta` var in the `posts` slice

### Naming Conventions

- **Files**: snake_case (`quiet_skills.templ`, `intro_blog.templ`)
- **Go types/functions**: PascalCase exported, camelCase private
- **Blog components**: `Blog<Name>(theme string)` for exported func
- **Blog meta functions**: `meta<Name>()` private templ component
- **Meta vars in blog_data.go**: `<Name>Data` (e.g. `QuietSkillsData`, `GrugAutomationData`)
- **URL slugs**: kebab-case matching the `case` in the route switch
- **Images**: WebP format only (CI enforces this — PNG/JPG cause build failure)

### Styling

- All styling via Tailwind utility classes + DaisyUI component classes (`navbar`, `card`, `badge`, `hero`, `btn`, `breadcrumbs`, `prose`, etc.)
- Two DaisyUI themes: `retro` (default) and `synthwave` (dark)
- Blog content wrapped in `prose` class for typography plugin
- Code blocks use PrismJS with prism-tomorrow theme

### State Management

- No client-side state (no React/Redux)
- Theme preference stored in Redis-backed SCS session (24h lifetime)
- Theme passed as `theme string` parameter to every page/blog component
- Theme toggled via htmx `hx-put="/theme"` on checkbox

### Session Helpers (`internal/cache.go`)

```go
internal.GetMessage("theme", r)    // Get theme from session (returns "retro" or "synthwave" or "")
internal.PutMessage("theme", val, r) // Set theme in session
internal.ClearSession(r)           // Destroy session
```

### Routing

All routes in `main.go` using Chi v5:
- `GET /` → IndexPage
- `PUT /theme` → Toggle theme
- `GET /blog/{id}` → Switch on `id` for individual posts
- `GET /styles.css`, `/prism.css`, `/js/prism.js` → Static files
- `GET /rss.xml` → RSS 2.0 feed
- `GET /images/*` → Whitelisted image serving (validated against startup-scanned file list)

## Gotchas & Important Notes

- **Missing route bug**: `static.GrugAutomationData.Url` points to `/blog/grug-automation` but there is NO matching `case "grug-automation":` in `setupBlogHandler()` in `main.go`. The component `blogs.GrugAutomationBlog` exists but is unreachable. Fix this if adding the route.
- **No tests** exist in the codebase. CI builds and deploys but does not run tests.
- **Generated files are gitignored**: `*_templ.go` (from templ), `output.css` (from Tailwind), compiled binaries.
- **Images must be WebP**. The CI pipeline (`go.yml`) scans for `.png`, `.jpg`, `.jpeg` and fails if any found.
- **Theme defaults to "retro"** if session is empty. The toggle swaps between "retro" and "synthwave".
- **Database** (Turso) is initialized at startup but currently NOT queried in any request handler. It's wired up and ready for use (e.g., view counts, comments).
- **Stoicism posts** use descriptive filenames (`control_and_choice.templ`, `to_be_steady.templ`) and live in a separate `stoicism` package.
- **Commit style**: Imperative mood, descriptive but varied in formatting. No squash merges. No conventional commit standard.

## CI/CD Pipeline

Push to `main`:
1. Verify all images are WebP
2. Docker build (multi-stage: Go 1.24 builder + Node 20 for Tailwind → Alpine runtime)
3. Push to `ghcr.io`
4. SCP `docker-compose.yaml` to VPS
5. SSH: pull new image, compose down/up, health check `curl -f -k https://localhost:443`
6. Prune old images

## Build Pipeline (Docker)

1. `go mod download` + `npm install`
2. `go tool templ generate`
3. `npx @tailwindcss/cli -i input.css -o output.css`
4. `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o blog`
5. Copy binary + `web/static/` to Alpine runtime image
