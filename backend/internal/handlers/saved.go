package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"spotslog/internal/middleware"
	"spotslog/internal/models"
)

type SavedHandler struct {
	DB *pgxpool.Pool
}

type createSavedRequest struct {
	PlaceID int `json:"place_id" binding:"required"`
}

func (h *SavedHandler) List(c *gin.Context) {
	userID := c.GetInt(middleware.ContextUserIDKey)

	rows, err := h.DB.Query(c.Request.Context(),
		`SELECT id, user_id, place_id, created_at FROM saved_places WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list saved places"})
		return
	}
	defer rows.Close()

	saved := []models.SavedPlace{}
	placeIDs := []int{}
	for rows.Next() {
		var s models.SavedPlace
		if err := rows.Scan(&s.ID, &s.UserID, &s.PlaceID, &s.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read saved places"})
			return
		}
		saved = append(saved, s)
		placeIDs = append(placeIDs, s.PlaceID)
	}

	if len(placeIDs) > 0 {
		placeRows, err := h.DB.Query(c.Request.Context(),
			`SELECT id, name, category, address, district, lat, lng, description,
			 price_range, opening_hours, menu, source, visibility, created_by, created_at, updated_at
			 FROM places WHERE id = ANY($1)`,
			placeIDs,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load saved place details"})
			return
		}
		defer placeRows.Close()

		places, err := scanPlaces(placeRows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read place details"})
			return
		}

		if err := attachCoverPhotos(c.Request.Context(), places, h.DB); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load place photos"})
			return
		}
		if err := attachRatings(c.Request.Context(), places, h.DB); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load ratings"})
			return
		}

		byID := map[int]models.Place{}
		for _, p := range places {
			byID[p.ID] = p
		}
		for i := range saved {
			if p, ok := byID[saved[i].PlaceID]; ok {
				saved[i].Place = &p
			}
		}
	}

	c.JSON(http.StatusOK, saved)
}

func (h *SavedHandler) Create(c *gin.Context) {
	var req createSavedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt(middleware.ContextUserIDKey)

	var saved models.SavedPlace
	err := h.DB.QueryRow(c.Request.Context(),
		`INSERT INTO saved_places (user_id, place_id) VALUES ($1, $2)
		 ON CONFLICT (user_id, place_id) DO UPDATE SET user_id = EXCLUDED.user_id
		 RETURNING id, user_id, place_id, created_at`,
		userID, req.PlaceID,
	).Scan(&saved.ID, &saved.UserID, &saved.PlaceID, &saved.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save place (does it exist?)"})
		return
	}

	c.JSON(http.StatusCreated, saved)
}

func (h *SavedHandler) Delete(c *gin.Context) {
	placeID, err := strconv.Atoi(c.Param("placeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place id"})
		return
	}

	userID := c.GetInt(middleware.ContextUserIDKey)

	if _, err := h.DB.Exec(c.Request.Context(),
		`DELETE FROM saved_places WHERE user_id = $1 AND place_id = $2`, userID, placeID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove saved place"})
		return
	}
	c.Status(http.StatusNoContent)
}
