# Stage 1: Build
FROM golang:1.24 AS builder

# Install Node.js for Tailwind CSS
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && apt-get install -y nodejs

WORKDIR /app

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy package.json for Node dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy source code
COPY . .

# Generate templ files and build CSS
RUN go get -tool github.com/a-h/templ/cmd/templ@latest && \
    go tool templ generate && \
    npx @tailwindcss/cli -i ./web/static/css/input.css -o ./web/static/css/output.css && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o blog

# Stage 2: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy only the binary and necessary static files from builder
COPY --from=builder /app/blog /app/blog
COPY --from=builder /app/web/static /app/web/static

ENV ENV=production

EXPOSE 443

CMD ["./blog"]
