![Version](https://img.shields.io/badge/version-1.0.0-blue)
# Spotslog

A places discovery and visit-tracking app for Jakarta — explore recommendations for cafés, restaurants, museums, and libraries, mark places
as visited, save favorites to a wishlist, and rate the places you've been.
Built as a full-stack portfolio project.

![DEMO](./frontend/src/assets/demo.gif)

## Features

- **Browse & filter** curated and community-submitted places by category and area
- **Auto-rotating slideshow** on the homepage highlighting the newest recommendations
- **Account registration** — create an account to unlock visit tracking, wishlist, and more.
- **Visit history** — mark places visited, attach notes and photos
- **Wishlist** — save places to visit later
- **Add place** — submit places you've visited that aren't in the recommendations yet, choosing whether to keep them private or share them publicly
- **Ratings** — rate places you've actually visited (1-5 stars), with an
  aggregate average shown on each place card
- **Admin role** — curate official recommendations, edit or remove any place
- **Interactive map** view of places using Leaflet
- **Auth** — register/login with JWT-based sessions

## Tech stack

**Backend**
- Go, [Gin](https://gin-gonic.com/) for routing
- PostgreSQL via [pgx](https://github.com/jackc/pgx)
- JWT auth ([golang-jwt](https://github.com/golang-jwt/jwt)), bcrypt password hashing
- S3-compatible object storage ([MinIO](https://min.io/)) for photo uploads
- Plain SQL migrations, no ORM

**Frontend**
- Vue 3 (Composition API) + TypeScript + Vite
- [Pinia](https://pinia.vuejs.org/) for state management
- [Leaflet](https://leafletjs.com/) + [leaflet-geosearch](https://github.com/smeijer/leaflet-geosearch) for maps and address search
- [Headless UI](https://headlessui.com/) for accessible dialogs/modals
- [Font Awesome](https://fontawesome.com/) for icons

## Getting started

### Prerequisites

- Go 1.25+
- Node.js 22.18+ (or 24.12+)
- PostgreSQL
- An S3-compatible storage service (e.g. local MinIO) for photo uploads

### Backend

```bash
cd backend
cp .env.example .env   # fill in DATABASE_URL, JWT_SECRET, S3 credentials, etc.
# apply migrations in backend/migrations/ against your database
go run ./cmd/api
```

The API runs on `http://localhost:8080` by default.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

The app runs on `http://localhost:5173` by default (Vite will pick the next
available port if that one's taken).

## Project structure

```
backend/
  cmd/api/            server entrypoint
  internal/
    config/            env/config loading
    db/                database connection pool
    handlers/           HTTP handlers (auth, places, visits, saved, ratings)
    middleware/          JWT auth, CORS
    models/               shared data structs
    storage/               S3/MinIO client for photo uploads
  migrations/          versioned SQL migrations

frontend/
  src/
    api/            typed API client
    components/        reusable UI (PlaceCard, PlaceMap, HeroSlideshow, ...)
    router/               route definitions and auth guard
    stores/               Pinia stores (auth, places, visits, saved, ratings)
    views/                 page-level components
```

## Author

Built by [indahjul](https://github.com/indahj) — [LinkedIn](https://www.linkedin.com/in/indahjuliani/)
