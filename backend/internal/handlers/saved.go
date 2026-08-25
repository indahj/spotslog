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
		`SELECT s.id, s.user_id, s.place_id, s.created_at,
		 p.id, p.name, p.category, p.address, p.district, p.lat, p.lng
		 FROM saved_places s JOIN places p ON p.id = s.place_id
		 WHERE s.user_id = $1 ORDER BY s.created_at DESC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list saved places"})
		return
	}
	defer rows.Close()

	saved := []models.SavedPlace{}
	for rows.Next() {
		var s models.SavedPlace
		var p models.Place
		if err := rows.Scan(&s.ID, &s.UserID, &s.PlaceID, &s.CreatedAt,
			&p.ID, &p.Name, &p.Category, &p.Address, &p.District, &p.Lat, &p.Lng); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read saved places"})
			return
		}
		s.Place = &p
		saved = append(saved, s)
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
