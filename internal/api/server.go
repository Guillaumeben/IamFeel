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

// GetSportIcon returns the emoji icon for a given sport
func GetSportIcon(sport string) string {
    switch sport {
    case "boxing":
        return "🥊"
    case "fitness":
        return "💪"
    case "running":
        return "🏃"
    case "bjj":
        return "🥋"
    case "cycling":
        return "🚴"
    case "swimming":
        return "🏊"
    case "crossfit":
        return "🏋️"
    default:
        return "🥊" // Default to boxing
    }
}

// GetSportIconForUser returns the sport icon for a user
func (s *Server) GetSportIconForUser(userID int) string {
    // Get user's primary sport from database
    sports, err := s.db.GetUserSports(userID)
    if err != nil || len(sports) == 0 {
        return "🥊" // Default to boxing
    }

    // Find primary sport
    for _, sport := range sports {
        if sport.IsPrimary {
            return GetSportIcon(sport.SportName)
        }
    }

    // If no primary, use first sport
    return GetSportIcon(sports[0].SportName)
}
