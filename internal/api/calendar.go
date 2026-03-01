package api

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "strings"
    "time"

    "github.com/tuxnam/iamfeel/internal/db"
)

// HandleCalendarExport exports the current week's plan to calendar format
func (s *Server) HandleCalendarExport(w http.ResponseWriter, r *http.Request) {
    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }

    calType := r.URL.Query().Get("type")
    if calType == "" {
        http.Error(w, "Calendar type required", http.StatusBadRequest)
        return
    }

    // Get current week's planned sessions
    now := time.Now()
    weekStart := now.AddDate(0, 0, -int(now.Weekday()))
    weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
    weekEnd := weekStart.AddDate(0, 0, 6)
    weekEnd = time.Date(weekEnd.Year(), weekEnd.Month(), weekEnd.Day(), 23, 59, 59, 0, weekEnd.Location())

    plannedSessions, err := s.db.GetPlannedSessions(user.ID, weekStart, weekEnd)
    if err != nil {
        log.Printf("Error getting planned sessions: %v", err)
        http.Error(w, "Failed to get planned sessions", http.StatusInternalServerError)
        return
    }

    switch calType {
    case "google":
        s.handleGoogleCalendar(w, r, plannedSessions)
    case "apple", "microsoft":
        s.handleICSExport(w, r, plannedSessions)
    default:
        http.Error(w, "Invalid calendar type", http.StatusBadRequest)
    }
}

// handleGoogleCalendar creates a Google Calendar URL
func (s *Server) handleGoogleCalendar(w http.ResponseWriter, r *http.Request, sessions []*db.TrainingSession) {
    if len(sessions) == 0 {
        http.Error(w, "No sessions to export", http.StatusBadRequest)
        return
    }

    // Google Calendar only supports one event at a time, so we'll create a URL for the first session
    // In practice, you might want to open multiple tabs or provide a list
    session := sessions[0]

    // Format: YYYYMMDDTHHMMSS
    startTime := time.Date(session.SessionDate.Year(), session.SessionDate.Month(), session.SessionDate.Day(), 9, 0, 0, 0, time.Local)
    endTime := startTime.Add(time.Duration(session.DurationMinutes) * time.Minute)

    dateStart := startTime.Format("20060102T150405")
    dateEnd := endTime.Format("20060102T150405")

    title := fmt.Sprintf("IamFeel: %s", session.SessionType)
    description := session.Notes

    googleURL := fmt.Sprintf("https://calendar.google.com/calendar/render?action=TEMPLATE&text=%s&dates=%s/%s&details=%s",
        url.QueryEscape(title),
        dateStart,
        dateEnd,
        url.QueryEscape(description),
    )

    http.Redirect(w, r, googleURL, http.StatusSeeOther)
}

// handleICSExport generates an ICS file for Apple Calendar and Outlook
func (s *Server) handleICSExport(w http.ResponseWriter, r *http.Request, sessions []*db.TrainingSession) {
    if len(sessions) == 0 {
        http.Error(w, "No sessions to export", http.StatusBadRequest)
        return
    }

    // Build ICS file content
    var icsBuilder strings.Builder
    icsBuilder.WriteString("BEGIN:VCALENDAR\r\n")
    icsBuilder.WriteString("VERSION:2.0\r\n")
    icsBuilder.WriteString("PRODID:-//IamFeel//Training Plan//EN\r\n")
    icsBuilder.WriteString("CALSCALE:GREGORIAN\r\n")
    icsBuilder.WriteString("METHOD:PUBLISH\r\n")

    for _, session := range sessions {
        // Set start time to 9 AM on the session date
        startTime := time.Date(session.SessionDate.Year(), session.SessionDate.Month(), session.SessionDate.Day(), 9, 0, 0, 0, time.Local)
        endTime := startTime.Add(time.Duration(session.DurationMinutes) * time.Minute)

        // Format timestamps in UTC
        dtStart := startTime.UTC().Format("20060102T150405Z")
        dtEnd := endTime.UTC().Format("20060102T150405Z")
        dtStamp := time.Now().UTC().Format("20060102T150405Z")

        // Generate unique ID
        uid := fmt.Sprintf("%d-%s@iamfeel.app", session.ID, dtStamp)

        icsBuilder.WriteString("BEGIN:VEVENT\r\n")
        icsBuilder.WriteString(fmt.Sprintf("UID:%s\r\n", uid))
        icsBuilder.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", dtStamp))
        icsBuilder.WriteString(fmt.Sprintf("DTSTART:%s\r\n", dtStart))
        icsBuilder.WriteString(fmt.Sprintf("DTEND:%s\r\n", dtEnd))
        icsBuilder.WriteString(fmt.Sprintf("SUMMARY:IamFeel: %s\r\n", session.SessionType))

        if session.Notes != "" {
            // Escape special characters in notes
            description := strings.ReplaceAll(session.Notes, "\n", "\\n")
            description = strings.ReplaceAll(description, ",", "\\,")
            icsBuilder.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", description))
        }

        icsBuilder.WriteString("END:VEVENT\r\n")
    }

    icsBuilder.WriteString("END:VCALENDAR\r\n")

    // Set headers for file download
    w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
    w.Header().Set("Content-Disposition", "attachment; filename=iamfeel-training-plan.ics")
    w.Write([]byte(icsBuilder.String()))
}

// HandleSessionCalendarExport exports a single session to ICS format
func (s *Server) HandleSessionCalendarExport(w http.ResponseWriter, r *http.Request) {
    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }

    sessionIDStr := r.URL.Query().Get("id")
    if sessionIDStr == "" {
        http.Error(w, "Session ID required", http.StatusBadRequest)
        return
    }

    // Parse session ID
    var sessionID int64
    if _, err := fmt.Sscanf(sessionIDStr, "%d", &sessionID); err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    // Get all planned sessions for the user to find the specific one
    // We use a wide date range to capture the session
    now := time.Now()
    startDate := now.AddDate(0, -1, 0) // 1 month ago
    endDate := now.AddDate(0, 2, 0)     // 2 months ahead

    sessions, err := s.db.GetPlannedSessions(user.ID, startDate, endDate)
    if err != nil {
        log.Printf("Error getting planned sessions: %v", err)
        http.Error(w, "Failed to get session", http.StatusInternalServerError)
        return
    }

    // Find the specific session by ID
    var session *db.TrainingSession
    for _, s := range sessions {
        if int64(s.ID) == sessionID {
            session = s
            break
        }
    }

    if session == nil {
        http.Error(w, "Session not found", http.StatusNotFound)
        return
    }

    // Export as ICS
    s.handleICSExport(w, r, []*db.TrainingSession{session})
}
