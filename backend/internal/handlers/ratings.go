package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"spotslog/internal/middleware"
)

type RatingsHandler struct {
	DB *pgxpool.Pool
}

type rateRequest struct {
	PlaceID int `json:"place_id" binding:"required"`
	Rating  int `json:"rating" binding:"required,min=1,max=5"`
}

func (h *RatingsHandler) Rate(c *gin.Context) {
	var req rateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt(middleware.ContextUserIDKey)

	var visited bool
	err := h.DB.QueryRow(c.Request.Context(),
		`SELECT  EXISTS(SELECT 1 FROM visits WHERE user_id = $1 AND place_id =$2)`,
		userID, req.PlaceID,
	).Scan(&visited)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check visit history"})
		return
	}
	if !visited {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only rate places you've visit"})
		return
	}

	_, err = h.DB.Exec(c.Request.Context(),
		`INSERT INTO place_ratings (user_id, place_id, rating) VALUES ($1, $2, $3) ON CONFLICT (user_id, place_id) DO UPDATE SET rating = EXCLUDED.rating, updated_at = now()`,
		userID, req.PlaceID, req.Rating,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"place_id": req.PlaceID, "rating": req.Rating})
}

func (h *RatingsHandler) GetMine(c *gin.Context) {
	placeID, err := strconv.Atoi(c.Param("placeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place id"})
		return
	}

	userID := c.GetInt(middleware.ContextUserIDKey)

	var rating int
	err = h.DB.QueryRow(c.Request.Context(),
		`SELECT rating FROM place_ratings WHERE user_id = $1 AND place_id = $2`, userID, placeID,
	).Scan(&rating)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"rating": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rating": rating})
}
