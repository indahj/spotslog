package handlers

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// SPA serves the built frontend from dir. Existing files are served as-is;
// any other GET falls back to index.html so client-side routes (e.g. /places/3)
// survive a page reload. Unknown /api paths still get a JSON 404.
func SPA(dir string) gin.HandlerFunc {
	root := http.Dir(dir)
	fileServer := http.FileServer(root)
	index := filepath.Join(dir, "index.html")

	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api/") || (c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		if f, err := root.Open(path.Clean(p)); err == nil {
			stat, statErr := f.Stat()
			f.Close()
			if statErr == nil && !stat.IsDir() {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		c.File(index)
	}
}
