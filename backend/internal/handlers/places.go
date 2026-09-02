package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"spotslog/internal/middleware"
	"spotslog/internal/models"
	"spotslog/internal/storage"
)

type PlacesHandler struct {
	DB      *pgxpool.Pool
	Storage *storage.Storage
}

// errForbidden is returned by the ownership guards so callers can bail out.
// It must never be nil on the deny path — a nil return there would let the
// handler continue and mutate a record it just refused access to.
var errForbidden = errors.New("forbidden")

type createPlaceRequest struct {
	Name         string         `json:"name" binding:"required"`
	Category     string         `json:"category" binding:"required,oneof=restaurant cafe museum library attraction"`
	Address      string         `json:"address" binding:"required"`
	District     *string        `json:"district"`
	Lat          float64        `json:"lat" binding:"required"`
	Lng          float64        `json:"lng" binding:"required"`
	Description  *string        `json:"description"`
	PriceRange   *string        `json:"price_range"`
	OpeningHours map[string]any `json:"opening_hours"`
	Menu         map[string]any `json:"menu"`
	Source       string         `json:"source" binding:"required,oneof=curated user"`
	Visibility   string         `json:"visibility" binding:"omitempty,oneof=public private"`
}

func (h *PlacesHandler) List(c *gin.Context) {
	category := c.Query("category")
	area := c.Query("area")
	search := c.Query("q")

	query := `SELECT id, name, category, address, district, lat, lng, description,
		price_range, opening_hours, menu, source, visibility, created_by, created_at, updated_at
		FROM places WHERE (visibility = 'public')`
	args := []any{}
	argN := 1

	if category != "" {
		query += " AND category = $" + itoa(argN)
		args = append(args, category)
		argN++
	}
	if area != "" {
		query += " AND district ILIKE $" + itoa(argN)
		args = append(args, "%"+area+"%")
		argN++
	}
	if search != "" {
		query += " AND name ILIKE $" + itoa(argN)
		args = append(args, "%"+search+"%")
		argN++
	}
	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := h.DB.Query(c.Request.Context(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list places"})
		return
	}
	defer rows.Close()

	places, err := scanPlaces(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read places"})
		return
	}

	if err := h.attachCoverPhotos(c.Request.Context(), places); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load place photos"})
		return
	}

	c.JSON(http.StatusOK, places)
}

func (h *PlacesHandler) Homepage(c *gin.Context) {
	rows, err := h.DB.Query(c.Request.Context(),
		`SELECT id, name, category, address, district, lat, lng, description,
		 price_range, opening_hours, menu, source, visibility, created_by, created_at, updated_at
		 FROM places
		 WHERE source = 'curated' OR (source = 'user' AND visibility = 'public')
		 ORDER BY created_at DESC LIMIT 50`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load homepage"})
		return
	}
	defer rows.Close()

	places, err := scanPlaces(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read places"})
		return
	}

	if err := h.attachCoverPhotos(c.Request.Context(), places); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load place photos"})
		return
	}

	c.JSON(http.StatusOK, places)
}

func (h *PlacesHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	place, err := h.getPlaceByID(c, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "place not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load place"})
		return
	}

	rows, err := h.DB.Query(c.Request.Context(),
		`SELECT id, place_id, url, uploaded_by, caption, created_at FROM place_photos WHERE place_id = $1 ORDER BY created_at`,
		id,
	)
	if err == nil {
		defer rows.Close()
		var photos []models.PlacePhoto
		for rows.Next() {
			var p models.PlacePhoto
			if err := rows.Scan(&p.ID, &p.PlaceID, &p.URL, &p.UploadedBy, &p.Caption, &p.CreatedAt); err == nil {
				photos = append(photos, p)
			}
		}
		c.JSON(http.StatusOK, gin.H{"place": place, "photos": photos})
		return
	}

	c.JSON(http.StatusOK, gin.H{"place": place, "photos": []models.PlacePhoto{}})
}

func (h *PlacesHandler) Create(c *gin.Context) {
	var req createPlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt(middleware.ContextUserIDKey)
	role := c.GetString(middleware.ContextRoleKey)

	if req.Source == "curated" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admins can create curated places"})
		return
	}

	visibility := req.Visibility
	if req.Source == "curated" {
		visibility = "public"
	} else if visibility == "" {
		visibility = "public"
	}

	openingHours, _ := json.Marshal(req.OpeningHours)
	menu, _ := json.Marshal(req.Menu)

	var place models.Place
	err := h.DB.QueryRow(c.Request.Context(),
		`INSERT INTO places (name, category, address, district, lat, lng, description,
		 price_range, opening_hours, menu, source, visibility, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		 RETURNING id, name, category, address, district, lat, lng, description,
		 price_range, opening_hours, menu, source, visibility, created_by, created_at, updated_at`,
		req.Name, req.Category, req.Address, req.District, req.Lat, req.Lng, req.Description,
		req.PriceRange, openingHours, menu, req.Source, visibility, userID,
	).Scan(&place.ID, &place.Name, &place.Category, &place.Address, &place.District, &place.Lat, &place.Lng,
		&place.Description, &place.PriceRange, &place.OpeningHours, &place.Menu, &place.Source, &place.Visibility,
		&place.CreatedBy, &place.CreatedAt, &place.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create place"})
		return
	}

	c.JSON(http.StatusCreated, place)
}

func (h *PlacesHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.requireOwnerOrAdmin(c, id); err != nil {
		return
	}

	var req createPlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	openingHours, _ := json.Marshal(req.OpeningHours)
	menu, _ := json.Marshal(req.Menu)

	_, err = h.DB.Exec(c.Request.Context(),
		`UPDATE places SET name=$1, category=$2, address=$3, district=$4, lat=$5, lng=$6,
		 description=$7, price_range=$8, opening_hours=$9, menu=$10, updated_at=now()
		 WHERE id=$11`,
		req.Name, req.Category, req.Address, req.District, req.Lat, req.Lng,
		req.Description, req.PriceRange, openingHours, menu, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update place"})
		return
	}

	place, err := h.getPlaceByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload place"})
		return
	}
	c.JSON(http.StatusOK, place)
}

func (h *PlacesHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.requireOwnerOrAdmin(c, id); err != nil {
		return
	}

	if _, err := h.DB.Exec(c.Request.Context(), `DELETE FROM places WHERE id = $1`, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete place"})
		return
	}
	c.Status(http.StatusNoContent)
}

type patchVisibilityRequest struct {
	Visibility string `json:"visibility" binding:"required,oneof=public private"`
}

func (h *PlacesHandler) PatchVisibility(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID := c.GetInt(middleware.ContextUserIDKey)

	var ownerID *int
	var source string
	err = h.DB.QueryRow(c.Request.Context(), `SELECT created_by, source FROM places WHERE id = $1`, id).Scan(&ownerID, &source)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "place not found"})
		return
	}
	if source == "curated" {
		c.JSON(http.StatusForbidden, gin.H{"error": "curated places are always public"})
		return
	}
	if ownerID == nil || *ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the owner can change visibility"})
		return
	}

	var req patchVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.DB.Exec(c.Request.Context(),
		`UPDATE places SET visibility = $1, updated_at = now() WHERE id = $2`, req.Visibility, id,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update visibility"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "visibility": req.Visibility})
}

func (h *PlacesHandler) UploadPhoto(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo file is required"})
		return
	}

	userID := c.GetInt(middleware.ContextUserIDKey)

	url, err := h.Storage.UploadPhoto(c.Request.Context(), "places/"+strconv.Itoa(id), fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload photo"})
		return
	}

	caption := c.PostForm("caption")
	var photo models.PlacePhoto
	err = h.DB.QueryRow(c.Request.Context(),
		`INSERT INTO place_photos (place_id, url, uploaded_by, caption) VALUES ($1,$2,$3,$4)
		 RETURNING id, place_id, url, uploaded_by, caption, created_at`,
		id, url, userID, nullIfEmpty(caption),
	).Scan(&photo.ID, &photo.PlaceID, &photo.URL, &photo.UploadedBy, &photo.Caption, &photo.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save photo record"})
		return
	}

	c.JSON(http.StatusCreated, photo)
}

func (h *PlacesHandler) DeletePhoto(c *gin.Context) {
	photoID, err := strconv.Atoi(c.Param("photoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid photo id"})
		return
	}

	userID := c.GetInt(middleware.ContextUserIDKey)
	role := c.GetString(middleware.ContextRoleKey)

	var uploadedBy *int
	if err := h.DB.QueryRow(c.Request.Context(), `SELECT uploaded_by FROM place_photos WHERE id = $1`, photoID).Scan(&uploadedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}
	if role != "admin" && (uploadedBy == nil || *uploadedBy != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to delete this photo"})
		return
	}

	if _, err := h.DB.Exec(c.Request.Context(), `DELETE FROM place_photos WHERE id = $1`, photoID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete photo"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PlacesHandler) requireOwnerOrAdmin(c *gin.Context, placeID int) error {
	userID := c.GetInt(middleware.ContextUserIDKey)
	role := c.GetString(middleware.ContextRoleKey)

	var ownerID *int
	err := h.DB.QueryRow(c.Request.Context(), `SELECT created_by FROM places WHERE id = $1`, placeID).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "place not found"})
		return err
	}
	if role == "admin" {
		return nil
	}
	if ownerID == nil || *ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to modify this place"})
		return errForbidden
	}
	return nil
}

func (h *PlacesHandler) getPlaceByID(c *gin.Context, id int) (models.Place, error) {
	var place models.Place
	err := h.DB.QueryRow(c.Request.Context(),
		`SELECT id, name, category, address, district, lat, lng, description,
		 price_range, opening_hours, menu, source, visibility, created_by, created_at, updated_at
		 FROM places WHERE id = $1`,
		id,
	).Scan(&place.ID, &place.Name, &place.Category, &place.Address, &place.District, &place.Lat, &place.Lng,
		&place.Description, &place.PriceRange, &place.OpeningHours, &place.Menu, &place.Source, &place.Visibility,
		&place.CreatedBy, &place.CreatedAt, &place.UpdatedAt)
	return place, err
}

func scanPlaces(rows pgx.Rows) ([]models.Place, error) {
	places := []models.Place{}
	for rows.Next() {
		var p models.Place
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Address, &p.District, &p.Lat, &p.Lng,
			&p.Description, &p.PriceRange, &p.OpeningHours, &p.Menu, &p.Source, &p.Visibility,
			&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		places = append(places, p)
	}
	return places, rows.Err()
}

// attachCoverPhotos fetches one representative photo per place (the
// earliest uploaded) and fills in each place's CoverPhotoURL — one batched
// query rather than one per place, same reasoning as the visit photos.
func (h *PlacesHandler) attachCoverPhotos(ctx context.Context, places []models.Place) error {
	if len(places) == 0 {
		return nil
	}

	ids := make([]int, len(places))
	byID := map[int]int{}
	for i, p := range places {
		ids[i] = p.ID
		byID[p.ID] = i
	}

	rows, err := h.DB.Query(ctx,
		`SELECT place_id, url FROM place_photos WHERE place_id = ANY($1) ORDER BY place_id, created_at ASC`,
		ids,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	seen := map[int]bool{}
	for rows.Next() {
		var placeID int
		var url string
		if err := rows.Scan(&placeID, &url); err != nil {
			return err
		}
		if seen[placeID] {
			continue
		}
		seen[placeID] = true
		if idx, ok := byID[placeID]; ok {
			places[idx].CoverPhotoURL = &url
		}
	}
	return rows.Err()
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
