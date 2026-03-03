package api

import (
    "log"
    "net/http"
    "strconv"

    "github.com/tuxnam/iamfeel/internal/db"
)

// UsersData holds data for the users page
type UsersData struct {
    CurrentUser *UserInfo
    AllUsers    []*UserInfo
    ThemeClass  string
    SportIcon   string
}

// UserInfo holds display information for a user
type UserInfo struct {
    ID              int
    Name            string
    Age             int
    ExperienceLevel string
    IsCurrent       bool
}

// HandleUsers renders the users management page
func (s *Server) HandleUsers(w http.ResponseWriter, r *http.Request) {
    // Try to get current user from cookie (don't use context since middleware is bypassed)
    var currentUser *db.User
    cookie, err := r.Cookie(userIDCookie)
    if err == nil && cookie.Value != "" {
        userID, err := strconv.Atoi(cookie.Value)
        if err == nil {
            currentUser, _ = s.db.GetUser(userID)
        }
    }

    allUsers, err := s.db.GetAllUsers()
    if err != nil {
        http.Error(w, "Failed to load users", http.StatusInternalServerError)
        log.Printf("Error loading users: %v", err)
        return
    }

    // Convert to UserInfo
    var usersInfo []*UserInfo
    for _, user := range allUsers {
        isCurrent := false
        if currentUser != nil {
            isCurrent = user.ID == currentUser.ID
        }
        usersInfo = append(usersInfo, &UserInfo{
            ID:              user.ID,
            Name:            user.Name,
            Age:             user.Age,
            ExperienceLevel: user.ExperienceLevel,
            IsCurrent:       isCurrent,
        })
    }

    var currentUserInfo *UserInfo
    themeClass := "theme-boxing" // Default
    sportIcon := "🥊"             // Default

    if currentUser != nil {
        currentUserInfo = &UserInfo{
            ID:              currentUser.ID,
            Name:            currentUser.Name,
            Age:             currentUser.Age,
            ExperienceLevel: currentUser.ExperienceLevel,
            IsCurrent:       true,
        }
        themeClass = s.GetThemeClass(currentUser.ID)
        sportIcon = s.GetSportIconForUser(currentUser.ID)
    }

    data := UsersData{
        CurrentUser: currentUserInfo,
        AllUsers:    usersInfo,
        ThemeClass:  themeClass,
        SportIcon:   sportIcon,
    }

    if err := s.templates.ExecuteTemplate(w, "users.html", data); err != nil {
        http.Error(w, "Failed to render template", http.StatusInternalServerError)
        log.Printf("Template error: %v", err)
    }
}

// HandleCreateUser handles creating a new user
func (s *Server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    name := r.FormValue("name")
    ageStr := r.FormValue("age")
    experienceLevel := r.FormValue("experience_level")
    primarySport := r.FormValue("primary_sport")

    if name == "" || ageStr == "" || experienceLevel == "" || primarySport == "" {
        http.Error(w, "All fields are required", http.StatusBadRequest)
        return
    }

    age, err := strconv.Atoi(ageStr)
    if err != nil {
        http.Error(w, "Invalid age", http.StatusBadRequest)
        return
    }

    // Create user in database
    user, err := s.db.CreateUser(name, age, experienceLevel)
    if err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        log.Printf("Error creating user: %v", err)
        return
    }

    // Create primary sport entry for this user
    err = s.db.CreateUserSport(user.ID, primarySport, true)
    if err != nil {
        log.Printf("Warning: Failed to create primary sport for user %d: %v", user.ID, err)
        // Don't fail the user creation, just log the warning
    }

    log.Printf("Created new user: %s (ID: %d) with primary sport: %s", user.Name, user.ID, primarySport)

    // Redirect to users page
    http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// HandleSwitchUser handles switching the current user
func (s *Server) HandleSwitchUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    userIDStr := r.FormValue("user_id")
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Verify user exists
    user, err := s.db.GetUser(userID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Set cookie for new user
    SetCurrentUserCookie(w, user.ID)

    log.Printf("Switched to user: %s (ID: %d)", user.Name, user.ID)

    // Redirect to dashboard
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleDeleteUser handles deleting a user
func (s *Server) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    userIDStr := r.FormValue("user_id")
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Get user before deleting for logging
    user, err := s.db.GetUser(userID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Check if this is the current user
    cookie, err := r.Cookie(userIDCookie)
    isCurrentUser := false
    if err == nil && cookie.Value != "" {
        currentUserID, err := strconv.Atoi(cookie.Value)
        if err == nil && currentUserID == userID {
            isCurrentUser = true
        }
    }

    // Delete user (CASCADE will handle related records)
    err = s.db.DeleteUser(userID)
    if err != nil {
        http.Error(w, "Failed to delete user", http.StatusInternalServerError)
        log.Printf("Error deleting user: %v", err)
        return
    }

    log.Printf("Deleted user: %s (ID: %d)", user.Name, user.ID)

    // If we deleted the current user, clear the cookie
    if isCurrentUser {
        http.SetCookie(w, &http.Cookie{
            Name:     userIDCookie,
            Value:    "",
            Path:     "/",
            MaxAge:   -1,
            HttpOnly: true,
            SameSite: http.SameSiteLaxMode,
        })
    }

    // Redirect to users page
    http.Redirect(w, r, "/users", http.StatusSeeOther)
}
