package api

import (
    "html/template"
    "path/filepath"

    "github.com/tuxnam/iamfeel/internal/config"
    "github.com/tuxnam/iamfeel/internal/db"
)

// Server holds the application state
type Server struct {
    db        *db.DB
    templates *template.Template
}

// NewServer creates a new server instance
func NewServer(database *db.DB) *Server {
    // Load templates
    templates := template.Must(template.ParseGlob(filepath.Join("web", "templates", "*.html")))

    return &Server{
        db:        database,
        templates: templates,
    }
}

// GetUserConfig loads the config for a specific user
func (s *Server) GetUserConfig(userID int) (*config.UserConfig, error) {
    return config.LoadUserConfigByID(userID)
}
