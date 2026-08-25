package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"spotslog/internal/middleware"
	"spotslog/internal/models"
	"spotslog/internal/storage"
)

type VisitsHandler struct {
	DB      *pgxpool.Pool
	Storage *storage.Storage
}

type createVisitRequest struct {
	PlaceID   int     `json:"place_id" binding:"required"`
	Notes     *string `json:"notes"`
	VisitedAt *string `json:"visited_at"`
}

func (h *VisitsHandler) List(c *gin.Context) {
	userID := c.GetInt(middleware.ContextUserIDKey)

	rows, err := h.DB.Query(c.Request.Context(),
		`SELECT v.id, v.user_id, v.place_id, v.visited_at, v.notes, v.created_at,
		 p.id, p.name, p.category, p.address, p.district, p.lat, p.lng
		 FROM visits v JOIN places p ON p.id = v.place_id
		 WHERE v.user_id = $1 ORDER BY v.visited_at DESC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list visits"})
		return
	}
	defer rows.Close()

	visits := []models.Visit{}
	byID := map[int]int{} // visit id -> index in visits
	for rows.Next() {
		var v models.Visit
		var p models.Place
		if err := rows.Scan(&v.ID, &v.UserID, &v.PlaceID, &v.VisitedAt, &v.Notes, &v.CreatedAt,
			&p.ID, &p.Name, &p.Category, &p.Address, &p.District, &p.Lat, &p.Lng); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read visits"})
			return
		}
		v.Place = &p
		v.Photos = []models.VisitPhoto{}
		byID[v.ID] = len(visits)
		visits = append(visits, v)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read visits"})
		return
	}

	if len(visits) == 0 {
		c.JSON(http.StatusOK, visits)
		return
	}

	// Attach photos in one extra query rather than one per visit — the UI
	// renders them inline, so a visit without its photos looks like the
	// upload silently failed.
	photoRows, err := h.DB.Query(c.Request.Context(),
		`SELECT vp.id, vp.visit_id, vp.url, vp.caption, vp.created_at
		 FROM visit_photos vp JOIN visits v ON v.id = vp.visit_id
		 WHERE v.user_id = $1 ORDER BY vp.created_at`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load visit photos"})
		return
	}
	defer photoRows.Close()

	for photoRows.Next() {
		var p models.VisitPhoto
		if err := photoRows.Scan(&p.ID, &p.VisitID, &p.URL, &p.Caption, &p.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read visit photos"})
			return
		}
		if idx, ok := byID[p.VisitID]; ok {
			visits[idx].Photos = append(visits[idx].Photos, p)
		}
	}

	c.JSON(http.StatusOK, visits)
}

func (h *VisitsHandler) Create(c *gin.Context) {
	var req createVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt(middleware.ContextUserIDKey)

	var visit models.Visit
	var err error
	if req.VisitedAt != nil {
		err = h.DB.QueryRow(c.Request.Context(),
			`INSERT INTO visits (user_id, place_id, notes, visited_at) VALUES ($1,$2,$3,$4)
			 RETURNING id, user_id, place_id, visited_at, notes, created_at`,
			userID, req.PlaceID, req.Notes, *req.VisitedAt,
		).Scan(&visit.ID, &visit.UserID, &visit.PlaceID, &visit.VisitedAt, &visit.Notes, &visit.CreatedAt)
	} else {
		err = h.DB.QueryRow(c.Request.Context(),
			`INSERT INTO visits (user_id, place_id, notes) VALUES ($1,$2,$3)
			 RETURNING id, user_id, place_id, visited_at, notes, created_at`,
			userID, req.PlaceID, req.Notes,
		).Scan(&visit.ID, &visit.UserID, &visit.PlaceID, &visit.VisitedAt, &visit.Notes, &visit.CreatedAt)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create visit (does the place exist?)"})
		return
	}

	c.JSON(http.StatusCreated, visit)
}

func (h *VisitsHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.requireOwner(c, id); err != nil {
		return
	}

	var req struct {
		Notes     *string `json:"notes"`
		VisitedAt *string `json:"visited_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.DB.Exec(c.Request.Context(),
		`UPDATE visits SET notes = COALESCE($1, notes), visited_at = COALESCE($2, visited_at) WHERE id = $3`,
		req.Notes, req.VisitedAt, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update visit"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *VisitsHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.requireOwner(c, id); err != nil {
		return
	}

	if _, err := h.DB.Exec(c.Request.Context(), `DELETE FROM visits WHERE id = $1`, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete visit"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *VisitsHandler) UploadPhoto(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.requireOwner(c, id); err != nil {
		return
	}

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo file is required"})
		return
	}

	url, err := h.Storage.UploadPhoto(c.Request.Context(), "visits/"+strconv.Itoa(id), fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload photo"})
		return
	}

	caption := c.PostForm("caption")
	var photo models.VisitPhoto
	err = h.DB.QueryRow(c.Request.Context(),
		`INSERT INTO visit_photos (visit_id, url, caption) VALUES ($1,$2,$3)
		 RETURNING id, visit_id, url, caption, created_at`,
		id, url, nullIfEmpty(caption),
	).Scan(&photo.ID, &photo.VisitID, &photo.URL, &photo.Caption, &photo.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save photo record"})
		return
	}

	c.JSON(http.StatusCreated, photo)
}

func (h *VisitsHandler) DeletePhoto(c *gin.Context) {
	photoID, err := strconv.Atoi(c.Param("photoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid photo id"})
		return
	}

	userID := c.GetInt(middleware.ContextUserIDKey)

	var visitOwnerID int
	err = h.DB.QueryRow(c.Request.Context(),
		`SELECT v.user_id FROM visit_photos vp JOIN visits v ON v.id = vp.visit_id WHERE vp.id = $1`,
		photoID,
	).Scan(&visitOwnerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}
	if visitOwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to delete this photo"})
		return
	}

	if _, err := h.DB.Exec(c.Request.Context(), `DELETE FROM visit_photos WHERE id = $1`, photoID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete photo"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *VisitsHandler) requireOwner(c *gin.Context, visitID int) error {
	userID := c.GetInt(middleware.ContextUserIDKey)

	var ownerID int
	err := h.DB.QueryRow(c.Request.Context(), `SELECT user_id FROM visits WHERE id = $1`, visitID).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "visit not found"})
		return err
	}
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to modify this visit"})
		return errForbidden
	}
	return nil
}
