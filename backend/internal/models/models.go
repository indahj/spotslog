package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Place struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	Category      string         `json:"category"`
	Address       string         `json:"address"`
	District      *string        `json:"district,omitempty"`
	Lat           float64        `json:"lat"`
	Lng           float64        `json:"lng"`
	Description   *string        `json:"description,omitempty"`
	PriceRange    *string        `json:"price_range,omitempty"`
	OpeningHours  map[string]any `json:"opening_hours,omitempty"`
	Menu          map[string]any `json:"menu,omitempty"`
	Source        string         `json:"source"`
	Visibility    string         `json:"visibility"`
	CreatedBy     *int           `json:"created_by,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CoverPhotoURL *string        `json:"cover_photo_url,omitempty"`
}

type PlacePhoto struct {
	ID         int       `json:"id"`
	PlaceID    int       `json:"place_id"`
	URL        string    `json:"url"`
	UploadedBy *int      `json:"uploaded_by,omitempty"`
	Caption    *string   `json:"caption,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type Visit struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PlaceID   int       `json:"place_id"`
	VisitedAt time.Time `json:"visited_at"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	Place  *Place       `json:"place,omitempty"`
	Photos []VisitPhoto `json:"photos,omitempty"`
}

type VisitPhoto struct {
	ID        int       `json:"id"`
	VisitID   int       `json:"visit_id"`
	URL       string    `json:"url"`
	Caption   *string   `json:"caption,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type SavedPlace struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PlaceID   int       `json:"place_id"`
	CreatedAt time.Time `json:"created_at"`

	Place *Place `json:"place,omitempty"`
}
