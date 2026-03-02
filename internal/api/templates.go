package api

import (
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/tuxnam/iamfeel/internal/db"
)

// TemplatesData holds data for the templates page
type TemplatesData struct {
    UserName   string
    Templates  []*db.SessionTemplate
    ThemeClass string
    SportIcon  string
}

// HandleTemplates renders the templates page
func (s *Server) HandleTemplates(w http.ResponseWriter, r *http.Request) {
    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    templates, err := s.db.GetSessionTemplates(user.ID)
    if err != nil {
        log.Printf("Error loading templates: %v", err)
        templates = []*db.SessionTemplate{}
    }

    data := TemplatesData{
        UserName:   user.Name,
        Templates:  templates,
        ThemeClass: s.GetThemeClass(user.ID),
        SportIcon:  s.GetSportIconForUser(user.ID),
    }

    if err := s.templates.ExecuteTemplate(w, "templates.html", data); err != nil {
        http.Error(w, "Failed to render template", http.StatusInternalServerError)
        log.Printf("Template error: %v", err)
    }
}

// HandleCreateTemplate creates a new session template
func (s *Server) HandleCreateTemplate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    templateName := r.FormValue("template_name")
    sportName := r.FormValue("sport_name")
    sessionType := r.FormValue("session_type")
    description := r.FormValue("description")

    durationMinutes, _ := strconv.Atoi(r.FormValue("duration_minutes"))
    perceivedEffort, _ := strconv.Atoi(r.FormValue("perceived_effort"))

    template := &db.SessionTemplate{
        UserID:          user.ID,
        TemplateName:    templateName,
        SportName:       sportName,
        SessionType:     sessionType,
        DurationMinutes: durationMinutes,
        PerceivedEffort: perceivedEffort,
        Description:     description,
    }

    if err := s.db.CreateSessionTemplate(template); err != nil {
        http.Error(w, "Failed to create template", http.StatusInternalServerError)
        log.Printf("Error creating template: %v", err)
        return
    }

    http.Redirect(w, r, "/templates", http.StatusSeeOther)
}

// HandleDeleteTemplate deletes a session template
func (s *Server) HandleDeleteTemplate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    templateID, err := strconv.Atoi(r.FormValue("template_id"))
    if err != nil {
        http.Error(w, "Invalid template ID", http.StatusBadRequest)
        return
    }

    // Verify template belongs to user
    template, err := s.db.GetSessionTemplate(templateID)
    if err != nil {
        http.Error(w, "Template not found", http.StatusNotFound)
        return
    }

    if template.UserID != user.ID {
        http.Error(w, "Unauthorized", http.StatusForbidden)
        return
    }

    if err := s.db.DeleteSessionTemplate(templateID); err != nil {
        http.Error(w, "Failed to delete template", http.StatusInternalServerError)
        log.Printf("Error deleting template: %v", err)
        return
    }

    http.Redirect(w, r, "/templates", http.StatusSeeOther)
}

// HandleUseTemplate creates a session from a template (quick-log)
func (s *Server) HandleUseTemplate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    templateID, err := strconv.Atoi(r.FormValue("template_id"))
    if err != nil {
        http.Error(w, "Invalid template ID", http.StatusBadRequest)
        return
    }

    sessionDate := r.FormValue("session_date")
    if sessionDate == "" {
        sessionDate = time.Now().Format("2006-01-02")
    }

    // Get template
    template, err := s.db.GetSessionTemplate(templateID)
    if err != nil {
        http.Error(w, "Template not found", http.StatusNotFound)
        return
    }

    // Verify template belongs to user
    if template.UserID != user.ID {
        http.Error(w, "Unauthorized", http.StatusForbidden)
        return
    }

    // Parse session date
    parsedDate, err := time.Parse("2006-01-02", sessionDate)
    if err != nil {
        http.Error(w, "Invalid date format", http.StatusBadRequest)
        return
    }

    // Create session from template
    session := &db.TrainingSession{
        UserID:          user.ID,
        SessionDate:     parsedDate,
        SessionType:     template.SessionType,
        DurationMinutes: template.DurationMinutes,
        PerceivedEffort: template.PerceivedEffort,
        Notes:           template.Description,
        Completed:       true,
        Skipped:         false,
        Planned:         false,
    }

    if err := s.db.CreateTrainingSession(session); err != nil {
        http.Error(w, "Failed to log session", http.StatusInternalServerError)
        log.Printf("Error creating session from template: %v", err)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}
