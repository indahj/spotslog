-- Spotslog initial schema

CREATE TYPE user_role AS ENUM ('user', 'admin');
CREATE TYPE place_category AS ENUM ('restaurant', 'cafe', 'museum', 'library', 'attraction');
CREATE TYPE place_source AS ENUM ('curated', 'user');
CREATE TYPE place_visibility AS ENUM ('public', 'private');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    avatar_url TEXT,
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE places (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    category place_category NOT NULL,
    address TEXT NOT NULL,
    district TEXT,
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL,
    description TEXT,
    price_range TEXT,
    opening_hours JSONB,
    menu JSONB,
    source place_source NOT NULL DEFAULT 'user',
    visibility place_visibility NOT NULL DEFAULT 'public',
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    -- Curated recommendations are always publicly visible; only user
    -- submissions get to choose between public and private.
    CONSTRAINT curated_is_always_public CHECK (
        source = 'user' OR visibility = 'public'
    )
);

CREATE INDEX idx_places_category ON places (category);
CREATE INDEX idx_places_lat_lng ON places (lat, lng);
CREATE INDEX idx_places_source_visibility ON places (source, visibility);

CREATE TABLE place_photos (
    id SERIAL PRIMARY KEY,
    place_id INTEGER NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    uploaded_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    caption TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_place_photos_place_id ON place_photos (place_id);

CREATE TABLE visits (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    place_id INTEGER NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    visited_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_visits_user_id ON visits (user_id);
CREATE INDEX idx_visits_place_id ON visits (place_id);

CREATE TABLE visit_photos (
    id SERIAL PRIMARY KEY,
    visit_id INTEGER NOT NULL REFERENCES visits(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    caption TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_visit_photos_visit_id ON visit_photos (visit_id);

CREATE TABLE saved_places (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    place_id INTEGER NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, place_id)
);

CREATE INDEX idx_saved_places_user_id ON saved_places (user_id);
