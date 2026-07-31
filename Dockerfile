# Stage 1: Build SvelteKit static UI
FROM node:22-alpine AS frontend-builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Stage 2: Build Go backend binary with embedded www/
FROM golang:alpine AS backend-builder

WORKDIR /app
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /app/www ./www

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o koala-github ./main.go

# Stage 3: Minimal Alpine production runtime image
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=backend-builder /app/koala-github /app/koala-github

# Environment variables
ENV PORT=8080
ENV DB_PATH=/data/koalagithub.db

# Persistent SQLite database volume
VOLUME ["/data"]

EXPOSE 8080

CMD ["/app/koala-github"]
