package api

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"

    "github.com/tuxnam/iamfeel/internal/db"
)

type contextKey string

const (
    userContextKey contextKey = "user"
    userIDCookie   string     = "iamfeel_user_id"
)

// UserMiddleware adds the current user to the request context
func (s *Server) UserMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Try to get user ID from cookie
        cookie, err := r.Cookie(userIDCookie)
        var user *db.User

        if err == nil && cookie.Value != "" {
            // Parse user ID from cookie
            userID, err := strconv.Atoi(cookie.Value)
            if err == nil {
                user, err = s.db.GetUser(userID)
                if err != nil {
                    log.Printf("Failed to get user from cookie: %v", err)
                }
            }
        }

        // If no user from cookie, get first user as fallback
        if user == nil {
            user, err = s.db.GetFirstUser()
            if err != nil {
                http.Error(w, "No users found. Please run onboarding first.", http.StatusInternalServerError)
                return
            }
            // Set cookie for first user
            http.SetCookie(w, &http.Cookie{
                Name:     userIDCookie,
                Value:    fmt.Sprintf("%d", user.ID),
                Path:     "/",
                MaxAge:   86400 * 365, // 1 year
                HttpOnly: true,
                SameSite: http.SameSiteLaxMode,
            })
        }

        // Add user to context
        ctx := context.WithValue(r.Context(), userContextKey, user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// GetCurrentUser retrieves the current user from the request context
func GetCurrentUser(r *http.Request) (*db.User, error) {
    user, ok := r.Context().Value(userContextKey).(*db.User)
    if !ok || user == nil {
        return nil, fmt.Errorf("no user in context")
    }
    return user, nil
}

// SetCurrentUserCookie sets the current user cookie
func SetCurrentUserCookie(w http.ResponseWriter, userID int) {
    http.SetCookie(w, &http.Cookie{
        Name:     userIDCookie,
        Value:    fmt.Sprintf("%d", userID),
        Path:     "/",
        MaxAge:   86400 * 365, // 1 year
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
    })
}

// OnboardingMiddleware checks if user has a config and redirects to onboarding if missing
func (s *Server) OnboardingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip onboarding check for certain routes
        if r.URL.Path == "/onboard" || strings.HasPrefix(r.URL.Path, "/static/") {
            next.ServeHTTP(w, r)
            return
        }

        user, err := GetCurrentUser(r)
        if err != nil {
            next.ServeHTTP(w, r)
            return
        }

        // Check if user has basic profile in database
        _, err = s.GetUserConfig(user.ID)
        if err != nil {
            http.Redirect(w, r, "/onboard", http.StatusSeeOther)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// GetThemeClass returns the CSS theme class based on user's primary sport
func (s *Server) GetThemeClass(userID int) string {
    primarySport, err := s.db.GetPrimarySport(userID)
    if err != nil {
        return "theme-boxing" // Default theme
    }
    return fmt.Sprintf("theme-%s", primarySport.SportName)
}
