# ---------- Frontend (Vue / Vite) ----------
FROM node:24-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
# The API is served by the same container, so the SPA calls it on the same origin.
ARG VITE_API_BASE_URL=/api
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
RUN npm run build-only

# ---------- Backend (Go) ----------
FROM golang:1.25-alpine AS backend
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/spotslog-api ./cmd/api

# ---------- Runtime ----------
FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata && adduser -D -H app
WORKDIR /app
COPY --from=backend /out/spotslog-api ./spotslog-api
COPY --from=frontend /app/dist ./public

ENV GIN_MODE=release \
    PORT=8080 \
    STATIC_DIR=/app/public

USER app
EXPOSE 8080
CMD ["/app/spotslog-api"]
