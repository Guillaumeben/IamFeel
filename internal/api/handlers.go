package api

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/tuxnam/iamfeel/internal/agent"
    "github.com/tuxnam/iamfeel/internal/config"
    "github.com/tuxnam/iamfeel/internal/db"
)

const (
    // DailyAICallLimit is the maximum number of AI calls allowed per user per day
    DailyAICallLimit = 10
)

// DashboardData holds data for the dashboard template
type DashboardData struct {
    UserName          string
    CurrentWeek       string
    CurrentDateTime   string
    Plan              *db.WeeklyPlan
    TodayDate         string
    Sessions          []*db.TrainingSession
    PlannedSessions   []*db.TrainingSession
    Supplements       []config.Supplement
    AIRecommendations string
    ThemeClass        string
    SportIcon         string
    WeeklyStats       *db.WeeklyStats
    StreakStats       *db.StreakStats
    MonthlyStats      *db.MonthlyStats
}

// HandleDashboard renders the main dashboard
func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    now := time.Now()
    weekStart := agent.GetWeekStart(now)

    // Get current week's plan
    plan, err := s.db.GetWeeklyPlan(user.ID, weekStart)
    if err != nil {
        // No plan exists yet, that's okay
        plan = nil
    }

    // Get this week's completed sessions
    weekEnd := weekStart.AddDate(0, 0, 6)
    sessions, err := s.db.GetTrainingSessions(user.ID, weekStart, weekEnd)
    if err != nil {
        log.Printf("Error loading sessions: %v", err)
        sessions = []*db.TrainingSession{}
    }

    // Get planned sessions for this week
    plannedSessions, err := s.db.GetPlannedSessions(user.ID, weekStart, weekEnd)
    if err != nil {
        log.Printf("Error loading planned sessions: %v", err)
        plannedSessions = []*db.TrainingSession{}
    }

    // Load user config
    userConfig, err := s.GetUserConfig(user.ID)
    if err != nil {
        log.Printf("Error loading user config: %v", err)
    }

    supplements := []config.Supplement{}
    if userConfig != nil {
        supplements = userConfig.Supplements
    }

    // Load stats for dashboard summary
    weeklyStats, err := s.db.GetWeeklyStats(user.ID, weekStart)
    if err != nil {
        log.Printf("Error loading weekly stats: %v", err)
        weeklyStats = nil
    }

    streakStats, err := s.db.GetStreakStats(user.ID)
    if err != nil {
        log.Printf("Error loading streak stats: %v", err)
        streakStats = nil
    }

    monthlyStats, err := s.db.GetMonthlyStats(user.ID, 30)
    if err != nil {
        log.Printf("Error loading monthly stats: %v", err)
        monthlyStats = nil
    }

    data := DashboardData{
        UserName:          user.Name,
        CurrentWeek:       fmt.Sprintf("%s - %s", weekStart.Format("Jan 2"), weekEnd.Format("Jan 2, 2006")),
        CurrentDateTime:   now.Format("Monday, January 2, 2006 at 3:04 PM"),
        Plan:              plan,
        TodayDate:         now.Format("2006-01-02"),
        Sessions:          sessions,
        PlannedSessions:   plannedSessions,
        Supplements:       supplements,
        AIRecommendations: "Stay consistent with your training schedule. Focus on recovery between intense sessions.",
        ThemeClass:        s.GetThemeClass(user.ID),
        SportIcon:         s.GetSportIconForUser(user.ID),
        WeeklyStats:       weeklyStats,
        StreakStats:       streakStats,
        MonthlyStats:      monthlyStats,
    }

    if err := s.templates.ExecuteTemplate(w, "dashboard.html", data); err != nil {
        http.Error(w, "Failed to render template", http.StatusInternalServerError)
        log.Printf("Template error: %v", err)
    }
}

// HistoryData holds data for the history page
type HistoryData struct {
    UserName        string
    Sessions        []*db.TrainingSession
    RestDayNotes    []*db.RestDayNote
    WeeklyStats     *db.WeeklyStats
    MonthlyStats    *db.MonthlyStats
    ActivityStats   []*db.ActivityStats
    StreakStats     *db.StreakStats
    VolumeData      []*db.VolumeDataPoint
    EffortData      []*db.EffortDataPoint
    ThemeClass      string
    SportIcon       string
    FilterStart     string
    FilterEnd       string
    FilterType      string
    FilterMinEffort string
    FilterMaxEffort string
    HasFilters      bool
}

// HandleHistory renders the training history page
func (s *Server) HandleHistory(w http.ResponseWriter, r *http.Request) {
    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    // Parse filter parameters from query string
    startDate := r.URL.Query().Get("start_date")
    endDate := r.URL.Query().Get("end_date")
    sessionType := r.URL.Query().Get("session_type")
    minEffortStr := r.URL.Query().Get("min_effort")
    maxEffortStr := r.URL.Query().Get("max_effort")

    // Convert effort strings to integers
    var minEffort, maxEffort int
    if minEffortStr != "" {
        fmt.Sscanf(minEffortStr, "%d", &minEffort)
    }
    if maxEffortStr != "" {
        fmt.Sscanf(maxEffortStr, "%d", &maxEffort)
    }

    // Check if any filters are applied
    hasFilters := startDate != "" || endDate != "" || sessionType != "" || minEffortStr != "" || maxEffortStr != ""

    var sessions []*db.TrainingSession
    if hasFilters {
        // Use filtered query
        filters := db.SessionFilters{
            StartDate:   startDate,
            EndDate:     endDate,
            SessionType: sessionType,
            MinEffort:   minEffort,
            MaxEffort:   maxEffort,
        }
        sessions, err = s.db.GetFilteredTrainingSessions(user.ID, filters, 100)
    } else {
        // Use default query (last 50 past sessions)
        sessions, err = s.db.GetPastTrainingSessions(user.ID, 50)
    }

    if err != nil {
        http.Error(w, "Failed to load sessions", http.StatusInternalServerError)
        log.Printf("Error loading sessions: %v", err)
        return
    }

    // Get rest day notes (last 50 by default, or filtered by date range)
    var restDayNotes []*db.RestDayNote
    if hasFilters && (startDate != "" || endDate != "") {
        // Apply date range filter
        var start, end time.Time
        if startDate != "" {
            start, _ = time.Parse("2006-01-02", startDate)
        } else {
            start = time.Now().AddDate(-1, 0, 0) // 1 year ago
        }
        if endDate != "" {
            end, _ = time.Parse("2006-01-02", endDate)
        } else {
            end = time.Now()
        }
        restDayNotes, err = s.db.GetRestDayNotes(user.ID, start, end)
    } else {
        restDayNotes, err = s.db.GetRecentRestDayNotes(user.ID, 50)
    }
    if err != nil {
        log.Printf("Error loading rest day notes: %v", err)
        restDayNotes = []*db.RestDayNote{}
    }

    // Load stats data
    weekStart := agent.GetWeekStart(time.Now())
    weeklyStats, err := s.db.GetWeeklyStats(user.ID, weekStart)
    if err != nil {
        log.Printf("Error loading weekly stats: %v", err)
        weeklyStats = nil
    }

    monthlyStats, err := s.db.GetMonthlyStats(user.ID, 30)
    if err != nil {
        log.Printf("Error loading monthly stats: %v", err)
        monthlyStats = nil
    }

    activityStats, err := s.db.GetActivityStats(user.ID, 30)
    if err != nil {
        log.Printf("Error loading activity stats: %v", err)
        activityStats = []*db.ActivityStats{}
    }

    streakStats, err := s.db.GetStreakStats(user.ID)
    if err != nil {
        log.Printf("Error loading streak stats: %v", err)
        streakStats = nil
    }

    volumeData, err := s.db.GetVolumeDataPoints(user.ID, 30)
    if err != nil {
        log.Printf("Error loading volume data: %v", err)
        volumeData = []*db.VolumeDataPoint{}
    }

    effortData, err := s.db.GetEffortDataPoints(user.ID, 30)
    if err != nil {
        log.Printf("Error loading effort data: %v", err)
        effortData = []*db.EffortDataPoint{}
    }

    data := HistoryData{
        UserName:        user.Name,
        Sessions:        sessions,
        RestDayNotes:    restDayNotes,
        WeeklyStats:     weeklyStats,
        MonthlyStats:    monthlyStats,
        ActivityStats:   activityStats,
        StreakStats:     streakStats,
        VolumeData:      volumeData,
        EffortData:      effortData,
        ThemeClass:      s.GetThemeClass(user.ID),
        SportIcon:       s.GetSportIconForUser(user.ID),
        FilterStart:     startDate,
        FilterEnd:       endDate,
        FilterType:      sessionType,
        FilterMinEffort: minEffortStr,
        FilterMaxEffort: maxEffortStr,
        HasFilters:      hasFilters,
    }

    if err := s.templates.ExecuteTemplate(w, "history.html", data); err != nil {
        http.Error(w, "Failed to render template", http.StatusInternalServerError)
        log.Printf("Template error: %v", err)
    }
}

// HandleLogSession handles logging a completed training session
func (s *Server) HandleLogSession(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    // Parse form
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    // Extract form values
    sessionDate := r.FormValue("session_date")
    sessionType := r.FormValue("session_type")
    durationStr := r.FormValue("duration_minutes")
    effortStr := r.FormValue("perceived_effort")
    notes := r.FormValue("notes")

    // Parse duration
    duration, err := strconv.Atoi(durationStr)
    if err != nil {
        http.Error(w, "Invalid duration", http.StatusBadRequest)
        return
    }

    // Parse effort
    effort, err := strconv.Atoi(effortStr)
    if err != nil || effort < 1 || effort > 10 {
        http.Error(w, "Invalid effort (must be 1-10)", http.StatusBadRequest)
        return
    }

    // Parse date
    date, err := time.Parse("2006-01-02", sessionDate)
    if err != nil {
        http.Error(w, "Invalid date format", http.StatusBadRequest)
        return
    }

    // Create session
    session := &db.TrainingSession{
        UserID:          user.ID,
        SessionDate:     date,
        SessionType:     sessionType,
        DurationMinutes: duration,
        PerceivedEffort: effort,
        Notes:           notes,
        Completed:       true,
        Planned:         false,
    }

    if err := s.db.CreateTrainingSession(session); err != nil {
        http.Error(w, "Failed to save session", http.StatusInternalServerError)
        log.Printf("Error saving session: %v", err)
        return
    }

    // Redirect back to dashboard
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleCurrentPlan returns the current week's plan as plain text
func (s *Server) HandleCurrentPlan(w http.ResponseWriter, r *http.Request) {
    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        return
    }

    weekStart := agent.GetWeekStart(time.Now())
    plan, err := s.db.GetWeeklyPlan(user.ID, weekStart)
    if err != nil {
        http.Error(w, "No plan found for this week", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(plan.PlanData))
}

// ModifyDayData holds data for the modify day page
type ModifyDayData struct {
    UserName    string
    CurrentWeek string
    Plan        *db.WeeklyPlan
    DayName     string
    Days        []string
    ThemeClass  string
    SportIcon   string
}

// HandleModifyDayForm shows the form to modify a day
func (s *Server) HandleModifyDayForm(w http.ResponseWriter, r *http.Request) {
    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    weekStart := agent.GetWeekStart(time.Now())
    plan, err := s.db.GetWeeklyPlan(user.ID, weekStart)
    if err != nil {
        http.Error(w, "No plan found for this week", http.StatusNotFound)
        return
    }

    weekEnd := weekStart.AddDate(0, 0, 6)
    days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

    dayName := r.URL.Query().Get("day")
    if dayName == "" {
        dayName = days[0]
    }

    data := ModifyDayData{
        UserName:    user.Name,
        CurrentWeek: fmt.Sprintf("%s - %s", weekStart.Format("Jan 2"), weekEnd.Format("Jan 2, 2006")),
        Plan:        plan,
        DayName:     dayName,
        Days:        days,
        ThemeClass:  s.GetThemeClass(user.ID),
        SportIcon:   s.GetSportIconForUser(user.ID),
    }

    if err := s.templates.ExecuteTemplate(w, "modify_day.html", data); err != nil {
        http.Error(w, "Failed to render template", http.StatusInternalServerError)
        log.Printf("Template error: %v", err)
    }
}

// ModifyDayProposalData holds data for the modification proposal page
type ModifyDayProposalData struct {
    UserName           string
    CurrentWeek        string
    DayName            string
    ModificationReason string
    NewDayPlan         string
    OriginalPlan       string
    ThemeClass         string
    SportIcon          string
}

// HandleModifyDayRequest generates a modified day plan
func (s *Server) HandleModifyDayRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    // Parse form
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    dayName := r.FormValue("day_name")
    modificationReason := r.FormValue("modification_reason")

    if dayName == "" || modificationReason == "" {
        http.Error(w, "Day name and modification reason are required", http.StatusBadRequest)
        return
    }

    weekStart := agent.GetWeekStart(time.Now())
    plan, err := s.db.GetWeeklyPlan(user.ID, weekStart)
    if err != nil {
        http.Error(w, "No plan found for this week", http.StatusNotFound)
        return
    }

    // Load user config
    userConfig, err := s.GetUserConfig(user.ID)
    if err != nil {
        http.Error(w, "Failed to load user config", http.StatusInternalServerError)
        log.Printf("Error loading user config: %v", err)
        return
    }

    // Check rate limit
    withinLimit, currentCount, err := s.db.CheckAIRateLimit(user.ID, DailyAICallLimit)
    if err != nil {
        http.Error(w, "Failed to check rate limit", http.StatusInternalServerError)
        log.Printf("Error checking rate limit: %v", err)
        return
    }
    if !withinLimit {
        http.Error(w, fmt.Sprintf("Daily AI call limit reached (%d/%d). Please try again tomorrow.", currentCount, DailyAICallLimit), http.StatusTooManyRequests)
        return
    }

    // Generate modified day
    ctx := r.Context()
    newDay, err := agent.GenerateDayModification(
        ctx,
        s.db,
        userConfig,
        weekStart,
        plan.PlanData,
        dayName,
        modificationReason,
    )
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to generate modification: %v", err), http.StatusInternalServerError)
        log.Printf("Error generating modification: %v", err)
        return
    }

    // Increment AI usage counter
    if err := s.db.IncrementAIUsage(user.ID, time.Now()); err != nil {
        log.Printf("Warning: Failed to increment AI usage: %v", err)
    }

    weekEnd := weekStart.AddDate(0, 0, 6)
    data := ModifyDayProposalData{
        UserName:           user.Name,
        CurrentWeek:        fmt.Sprintf("%s - %s", weekStart.Format("Jan 2"), weekEnd.Format("Jan 2, 2006")),
        DayName:            dayName,
        ModificationReason: modificationReason,
        NewDayPlan:         newDay,
        OriginalPlan:       plan.PlanData,
        ThemeClass:         s.GetThemeClass(user.ID),
        SportIcon:          s.GetSportIconForUser(user.ID),
    }

    if err := s.templates.ExecuteTemplate(w, "modify_proposal.html", data); err != nil {
        http.Error(w, "Failed to render template", http.StatusInternalServerError)
        log.Printf("Template error: %v", err)
    }
}

// HandleAcceptModification accepts and saves the modified plan
func (s *Server) HandleAcceptModification(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    // Parse form
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    dayName := r.FormValue("day_name")
    newDayPlan := r.FormValue("new_day_plan")

    weekStart := agent.GetWeekStart(time.Now())
    plan, err := s.db.GetWeeklyPlan(user.ID, weekStart)
    if err != nil {
        http.Error(w, "No plan found for this week", http.StatusNotFound)
        return
    }

    // Replace the day in the plan (simple text replacement)
    // This is basic - could be improved with better parsing
    updatedPlan := replaceDayInPlan(plan.PlanData, dayName, newDayPlan)

    // Update the plan in database
    plan.PlanData = updatedPlan
    if err := s.db.UpdateWeeklyPlan(plan); err != nil {
        http.Error(w, "Failed to update plan", http.StatusInternalServerError)
        log.Printf("Error updating plan: %v", err)
        return
    }

    // Redirect back to dashboard
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleCompleteSession marks a planned session as complete
func (s *Server) HandleCompleteSession(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    sessionIDStr := r.FormValue("session_id")
    durationStr := r.FormValue("duration_minutes")
    effortStr := r.FormValue("perceived_effort")
    notes := r.FormValue("notes")
    performanceNotes := r.FormValue("performance_notes")

    sessionID, err := strconv.Atoi(sessionIDStr)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    duration, err := strconv.Atoi(durationStr)
    if err != nil {
        http.Error(w, "Invalid duration", http.StatusBadRequest)
        return
    }

    effort, err := strconv.Atoi(effortStr)
    if err != nil || effort < 1 || effort > 10 {
        http.Error(w, "Invalid effort (must be 1-10)", http.StatusBadRequest)
        return
    }

    // Get the planned session
    weekStart := agent.GetWeekStart(time.Now())
    weekEnd := weekStart.AddDate(0, 0, 6)
    sessions, err := s.db.GetPlannedSessions(user.ID, weekStart, weekEnd)
    if err != nil {
        http.Error(w, "Failed to load planned sessions", http.StatusInternalServerError)
        return
    }

    var session *db.TrainingSession
    for _, sess := range sessions {
        if sess.ID == sessionID {
            session = sess
            break
        }
    }

    if session == nil {
        http.Error(w, "Session not found", http.StatusNotFound)
        return
    }

    // Update the session
    session.Completed = true
    session.DurationMinutes = duration
    session.PerceivedEffort = effort
    session.Notes = notes
    session.PerformanceNotes = performanceNotes
    session.Skipped = false

    if err := s.db.UpdateTrainingSession(session); err != nil {
        http.Error(w, "Failed to update session", http.StatusInternalServerError)
        log.Printf("Error updating session: %v", err)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleSkipSession marks a planned session as skipped
func (s *Server) HandleSkipSession(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    sessionIDStr := r.FormValue("session_id")
    skipReason := r.FormValue("skip_reason")

    sessionID, err := strconv.Atoi(sessionIDStr)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    if skipReason == "" {
        http.Error(w, "Skip reason is required", http.StatusBadRequest)
        return
    }

    // Get the planned session
    weekStart := agent.GetWeekStart(time.Now())
    weekEnd := weekStart.AddDate(0, 0, 6)
    sessions, err := s.db.GetPlannedSessions(user.ID, weekStart, weekEnd)
    if err != nil {
        http.Error(w, "Failed to load planned sessions", http.StatusInternalServerError)
        return
    }

    var session *db.TrainingSession
    for _, sess := range sessions {
        if sess.ID == sessionID {
            session = sess
            break
        }
    }

    if session == nil {
        http.Error(w, "Session not found", http.StatusNotFound)
        return
    }

    // Update the session
    session.Skipped = true
    session.SkipReason = skipReason
    session.Completed = false

    if err := s.db.UpdateTrainingSession(session); err != nil {
        http.Error(w, "Failed to update session", http.StatusInternalServerError)
        log.Printf("Error updating session: %v", err)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// replaceDayInPlan replaces a specific day's section in the plan
// This is a simple implementation - could be improved with better parsing
func replaceDayInPlan(originalPlan string, dayName string, newDayPlan string) string {
    // Find the day section
    dayMarker := fmt.Sprintf("\n%s:", dayName)
    startIdx := -1

    for i := 0; i < len(originalPlan)-len(dayMarker); i++ {
        if originalPlan[i:i+len(dayMarker)] == dayMarker {
            startIdx = i + 1 // Include the newline
            break
        }
    }

    if startIdx == -1 {
        // Day not found, append at end
        return originalPlan + "\n\n" + newDayPlan
    }

    // Find the next day or end of plan
    nextDays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday", "NOTES:"}
    endIdx := len(originalPlan)

    for _, nextDay := range nextDays {
        if nextDay == dayName {
            continue
        }
        marker := fmt.Sprintf("\n%s:", nextDay)
        for i := startIdx; i < len(originalPlan)-len(marker); i++ {
            if originalPlan[i:i+len(marker)] == marker {
                if i < endIdx {
                    endIdx = i
                }
                break
            }
        }
    }

    // Replace the section
    return originalPlan[:startIdx] + newDayPlan + "\n" + originalPlan[endIdx:]
}

// SettingsData holds data for the settings template
type SettingsData struct {
    User          *db.User
    Config        *config.UserConfig
    UserActivities []*db.UserSport
    Success       bool
    Error         string
    ThemeClass    string
    SportIcon     string
    AIUsageCount  int
    AIUsageLimit  int
}

// NewSettingsData creates a SettingsData struct with common fields populated
func (s *Server) NewSettingsData(user *db.User, userConfig *config.UserConfig) SettingsData {
    aiUsageCount, err := s.db.GetDailyAIUsage(user.ID, time.Now())
    if err != nil {
        log.Printf("Warning: Failed to get AI usage: %v", err)
        aiUsageCount = 0
    }

    // Load user activities for the activity selector
    userActivities, err := s.db.GetUserSports(user.ID)
    if err != nil {
        log.Printf("Warning: Failed to get user activities: %v", err)
        userActivities = []*db.UserSport{}
    }

    return SettingsData{
        User:           user,
        Config:         userConfig,
        UserActivities: userActivities,
        ThemeClass:     s.GetThemeClass(user.ID),
        SportIcon:      s.GetSportIconForUser(user.ID),
        AIUsageCount:   aiUsageCount,
        AIUsageLimit:   DailyAICallLimit,
    }
}

// HandleSettings displays the settings page
func (s *Server) HandleSettings(w http.ResponseWriter, r *http.Request) {
    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    // Load user config, or create default if doesn't exist
    userConfig, err := s.GetUserConfig(user.ID)
    if err != nil {
        log.Printf("ERROR loading user config for user %d: %v - CREATING DEFAULT CONFIG", user.ID, err)
        // Create a minimal default config
        userConfig = &config.UserConfig{
            User: config.UserProfile{
                Name:            user.Name,
                Age:             user.Age,
                Weight:          user.Weight,
                Height:          user.Height,
                ExperienceLevel: user.ExperienceLevel,
            },
            Sports: []config.UserSport{
                {
                    Name:    "boxing",
                    Primary: true,
                },
            },
            Equipment: config.EquipmentAccess{
                Home: []string{},
                Gyms: []config.Gym{},
            },
            Availability: make(map[string]config.DayAvailability),
            Goals: config.Goals{
                ShortTerm:  []string{},
                MediumTerm: []string{},
                LongTerm:   []string{},
            },
            Fitness:     &config.FitnessBaseline{},
            Supplements: []config.Supplement{},
            Preferences: config.UserPreferences{
                SessionsPerWeek:   3,
                PreferredDuration: 60,
            },
            Coach: config.CoachSettings{
                Model:             "claude-haiku-4-5",
                Temperature:       0.7,
                CoachingStyle:     "motivational",
                ExplanationDetail: "balanced",
            },
            Tracking: config.TrackingSettings{
                HistoryMonths:    6,
                TrackSupplements: true,
                TrackSleep:       false,
                TrackWeight:      false,
            },
        }
        log.Printf("Created default config for user %d", user.ID)
    } else {
        log.Printf("Successfully loaded user config for user %d - Availability days: %d, Equipment: %d, Goals: %d/%d/%d",
            user.ID, len(userConfig.Availability), len(userConfig.Equipment.Home),
            len(userConfig.Goals.ShortTerm), len(userConfig.Goals.MediumTerm), len(userConfig.Goals.LongTerm))
    }

    data := s.NewSettingsData(user, userConfig)

    if err := s.templates.ExecuteTemplate(w, "settings.html", data); err != nil {
        http.Error(w, "Failed to render settings", http.StatusInternalServerError)
        log.Printf("Template error: %v", err)
    }
}

// HandleSettingsSave saves updated settings
func (s *Server) HandleSettingsSave(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    // Get current user
    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    // Load user config, or create default if doesn't exist
    userConfig, err := s.GetUserConfig(user.ID)
    if err != nil {
        // Create a minimal default config
        userConfig = &config.UserConfig{
            User: config.UserProfile{
                Name:            user.Name,
                Age:             user.Age,
                Weight:          user.Weight,
                Height:          user.Height,
                ExperienceLevel: user.ExperienceLevel,
            },
            Sports: []config.UserSport{
                {
                    Name:    "boxing",
                    Primary: true,
                },
            },
            Equipment: config.EquipmentAccess{
                Home: []string{},
                Gyms: []config.Gym{},
            },
            Availability: make(map[string]config.DayAvailability),
            Goals: config.Goals{
                ShortTerm:  []string{},
                MediumTerm: []string{},
                LongTerm:   []string{},
            },
            Fitness:     &config.FitnessBaseline{},
            Supplements: []config.Supplement{},
            Preferences: config.UserPreferences{
                SessionsPerWeek:   3,
                PreferredDuration: 60,
            },
            Coach: config.CoachSettings{
                Model:             "claude-haiku-4-5",
                Temperature:       0.7,
                CoachingStyle:     "motivational",
                ExplanationDetail: "balanced",
            },
            Tracking: config.TrackingSettings{
                HistoryMonths:    6,
                TrackSupplements: true,
                TrackSleep:       false,
                TrackWeight:      false,
            },
        }
        log.Printf("Created default config for user %d during save", user.ID)
    }

    // Update user profile
    if name := r.FormValue("user_name"); name != "" {
        user.Name = strings.TrimSpace(name)
        userConfig.User.Name = user.Name
    }
    if ageStr := r.FormValue("user_age"); ageStr != "" {
        if age, err := strconv.Atoi(ageStr); err == nil {
            user.Age = age
            userConfig.User.Age = age
        }
    }
    if weightStr := r.FormValue("user_weight"); weightStr != "" {
        if weight, err := strconv.ParseFloat(weightStr, 64); err == nil {
            user.Weight = weight
            userConfig.User.Weight = weight
        }
    }
    if heightStr := r.FormValue("user_height"); heightStr != "" {
        if height, err := strconv.ParseFloat(heightStr, 64); err == nil {
            user.Height = height
            userConfig.User.Height = height
        }
    }
    if experience := r.FormValue("user_experience"); experience != "" {
        user.ExperienceLevel = experience
        userConfig.User.ExperienceLevel = experience
    }

    // Update user in database
    if err := s.db.UpdateUser(user); err != nil {
        log.Printf("Failed to update user: %v", err)
        data := s.NewSettingsData(user, userConfig)
        data.Error = fmt.Sprintf("Failed to update user: %v", err)
        s.templates.ExecuteTemplate(w, "settings.html", data)
        return
    }

    // Update primary sport and experience years together
    if sport := r.FormValue("primary_sport"); sport != "" {
        experienceYears := 0
        if experienceYearsStr := r.FormValue("sport_experience_years"); experienceYearsStr != "" {
            if years, err := strconv.Atoi(experienceYearsStr); err == nil {
                experienceYears = years
            }
        }

        // Update existing sports or create new one
        found := false
        for i := range userConfig.Sports {
            if userConfig.Sports[i].Primary {
                userConfig.Sports[i].Name = sport
                userConfig.Sports[i].ExperienceYears = experienceYears
                found = true
                break
            }
        }
        if !found {
            userConfig.Sports = []config.UserSport{{
                Name:            sport,
                Primary:         true,
                ExperienceYears: experienceYears,
            }}
        }
    }

    // Update goals
    if shortTermText := r.FormValue("short_term_goals"); shortTermText != "" {
        goals := []string{}
        for _, line := range strings.Split(shortTermText, "\n") {
            line = strings.TrimSpace(line)
            if line != "" {
                goals = append(goals, line)
            }
        }
        userConfig.Goals.ShortTerm = goals
    } else {
        userConfig.Goals.ShortTerm = []string{}
    }

    if mediumTermText := r.FormValue("medium_term_goals"); mediumTermText != "" {
        goals := []string{}
        for _, line := range strings.Split(mediumTermText, "\n") {
            line = strings.TrimSpace(line)
            if line != "" {
                goals = append(goals, line)
            }
        }
        userConfig.Goals.MediumTerm = goals
    } else {
        userConfig.Goals.MediumTerm = []string{}
    }

    if longTermText := r.FormValue("long_term_goals"); longTermText != "" {
        goals := []string{}
        for _, line := range strings.Split(longTermText, "\n") {
            line = strings.TrimSpace(line)
            if line != "" {
                goals = append(goals, line)
            }
        }
        userConfig.Goals.LongTerm = goals
    } else {
        userConfig.Goals.LongTerm = []string{}
    }

    // Update availability
    days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
    for _, day := range days {
        morning := r.FormValue(fmt.Sprintf("avail_%s_morning", day)) != ""
        lunch := r.FormValue(fmt.Sprintf("avail_%s_lunch", day)) != ""
        evening := r.FormValue(fmt.Sprintf("avail_%s_evening", day)) != ""

        if userConfig.Availability == nil {
            userConfig.Availability = make(map[string]config.DayAvailability)
        }

        userConfig.Availability[day] = config.DayAvailability{
            Morning: morning,
            Lunch:   lunch,
            Evening: evening,
        }
    }

    // Update preferences
    if goal := r.FormValue("primary_goal"); goal != "" {
        userConfig.Preferences.PrimaryGoal = goal
    }
    if sessionsStr := r.FormValue("sessions_per_week"); sessionsStr != "" {
        if sessions, err := strconv.ParseFloat(sessionsStr, 64); err == nil {
            userConfig.Preferences.SessionsPerWeek = sessions
        }
    }
    allowShortSessions := r.FormValue("allow_short_sessions") == "on"
    userConfig.Preferences.AllowShortSessions = allowShortSessions

    if maxSessionsStr := r.FormValue("max_sessions_per_day"); maxSessionsStr != "" {
        if maxSessions, err := strconv.Atoi(maxSessionsStr); err == nil {
            userConfig.Preferences.MaxSessionsPerDay = maxSessions
        }
    }
    if durationStr := r.FormValue("session_duration"); durationStr != "" {
        if duration, err := strconv.Atoi(durationStr); err == nil {
            userConfig.Preferences.PreferredDuration = duration
        }
    }
    notes := r.FormValue("training_notes")
    userConfig.Preferences.Notes = strings.TrimSpace(notes)

    // Update home equipment
    homeEquipment := []string{}
    for _, item := range r.Form["home_equipment[]"] {
        item = strings.TrimSpace(item)
        if item != "" {
            homeEquipment = append(homeEquipment, item)
        }
    }
    userConfig.Equipment.Home = homeEquipment

    // Update gyms
    gymNames := r.Form["gym_name[]"]
    gymTypes := r.Form["gym_type[]"]
    gymMemberships := r.Form["gym_membership[]"]
    gymSportIDs := r.Form["gym_sport_id[]"]

    newGyms := []config.Gym{}
    for i := range gymNames {
        name := strings.TrimSpace(gymNames[i])
        if name == "" {
            continue
        }

        gym := config.Gym{
            Name:       name,
            Type:       gymTypes[i],
            Membership: gymMemberships[i],
            Sessions:   []config.ClubSession{},
        }

        // Get sessions for this gym
        sessionNames := r.Form[fmt.Sprintf("gym_%d_session_name[]", i)]
        sessionDescriptions := r.Form[fmt.Sprintf("gym_%d_session_description[]", i)]
        sessionOccurrences := r.Form[fmt.Sprintf("gym_%d_session_occurrences[]", i)]
        sessionDurations := r.Form[fmt.Sprintf("gym_%d_session_duration[]", i)]
        sessionCosts := r.Form[fmt.Sprintf("gym_%d_session_cost[]", i)]

        for j := range sessionNames {
            sessionName := strings.TrimSpace(sessionNames[j])
            if sessionName == "" {
                continue
            }

            sessionDesc := ""
            if j < len(sessionDescriptions) {
                sessionDesc = strings.TrimSpace(sessionDescriptions[j])
            }

            sessionOccur := ""
            if j < len(sessionOccurrences) {
                sessionOccur = strings.TrimSpace(sessionOccurrences[j])
            }

            sessionDuration := ""
            if j < len(sessionDurations) {
                sessionDuration = strings.TrimSpace(sessionDurations[j])
            }

            sessionCost := ""
            if j < len(sessionCosts) {
                sessionCost = strings.TrimSpace(sessionCosts[j])
            }

            session := config.ClubSession{
                Name:        sessionName,
                Description: sessionDesc,
                Occurrences: sessionOccur,
                Duration:    sessionDuration,
                Cost:        sessionCost,
            }
            gym.Sessions = append(gym.Sessions, session)
        }

        newGyms = append(newGyms, gym)
    }
    userConfig.Equipment.Gyms = newGyms

    // Save to database instead of YAML
    // 1. Save primary sport
    if len(userConfig.Sports) > 0 && userConfig.Sports[0].Name != "" {
        if err := s.db.UpdatePrimarySport(user.ID, userConfig.Sports[0].Name, userConfig.Sports[0].ExperienceYears); err != nil {
            log.Printf("Failed to save primary sport: %v", err)
            data := s.NewSettingsData(user, userConfig)
            data.Error = fmt.Sprintf("Failed to save primary sport: %v", err)
            s.templates.ExecuteTemplate(w, "settings.html", data)
            return
        }
    }

    // 2. Save goals - delete all and re-insert
    if err := s.db.DeleteUserGoals(user.ID); err != nil {
        log.Printf("Failed to delete existing goals: %v", err)
        data := s.NewSettingsData(user, userConfig)
        data.Error = fmt.Sprintf("Failed to save goals: %v", err)
        s.templates.ExecuteTemplate(w, "settings.html", data)
        return
    }
    for _, goal := range userConfig.Goals.ShortTerm {
        if err := s.db.CreateGoalSimple(user.ID, "short_term", goal); err != nil {
            log.Printf("Failed to save short-term goal: %v", err)
        }
    }
    for _, goal := range userConfig.Goals.MediumTerm {
        if err := s.db.CreateGoalSimple(user.ID, "medium_term", goal); err != nil {
            log.Printf("Failed to save medium-term goal: %v", err)
        }
    }
    for _, goal := range userConfig.Goals.LongTerm {
        if err := s.db.CreateGoalSimple(user.ID, "long_term", goal); err != nil {
            log.Printf("Failed to save long-term goal: %v", err)
        }
    }

    // 3. Save availability
    for day, avail := range userConfig.Availability {
        dbAvail := &db.Availability{
            UserID:    user.ID,
            DayOfWeek: day,
            Morning:   avail.Morning,
            Lunch:     avail.Lunch,
            Evening:   avail.Evening,
        }
        if err := s.db.UpsertAvailability(dbAvail); err != nil {
            log.Printf("Failed to save availability for %s: %v", day, err)
        }
    }

    // 4. Save preferences
    prefs := &db.UserPreferences{
        UserID:                    user.ID,
        PrimaryGoal:               userConfig.Preferences.PrimaryGoal,
        SessionsPerWeek:           userConfig.Preferences.SessionsPerWeek,
        PreferredDuration:         userConfig.Preferences.PreferredDuration,
        SessionDurationPreference: userConfig.Preferences.SessionDurationPreference,
        IntensityPreference:       userConfig.Preferences.IntensityPreference,
        RecoveryPriority:          userConfig.Preferences.RecoveryPriority,
        PlanFrequency:             userConfig.Preferences.PlanFrequency,
        AllowShortSessions:        userConfig.Preferences.AllowShortSessions,
        MaxSessionsPerDay:         userConfig.Preferences.MaxSessionsPerDay,
        Notes:                     userConfig.Preferences.Notes,
    }
    if err := s.db.UpsertUserPreferences(prefs); err != nil {
        log.Printf("Failed to save preferences: %v", err)
        data := s.NewSettingsData(user, userConfig)
        data.Error = fmt.Sprintf("Failed to save preferences: %v", err)
        s.templates.ExecuteTemplate(w, "settings.html", data)
        return
    }

    // 5. Save coach settings
    if coachModel := r.FormValue("coach_model"); coachModel != "" {
        userConfig.Coach.Model = coachModel
    }
    if coachTemp := r.FormValue("coach_temperature"); coachTemp != "" {
        if temp, err := strconv.ParseFloat(coachTemp, 64); err == nil {
            userConfig.Coach.Temperature = temp
        }
    }
    if coachingStyle := r.FormValue("coaching_style"); coachingStyle != "" {
        userConfig.Coach.CoachingStyle = coachingStyle
    }
    if explanationDetail := r.FormValue("explanation_detail"); explanationDetail != "" {
        userConfig.Coach.ExplanationDetail = explanationDetail
    }

    coachSettings := &db.CoachSettings{
        UserID:            user.ID,
        Model:             userConfig.Coach.Model,
        Temperature:       userConfig.Coach.Temperature,
        CoachingStyle:     userConfig.Coach.CoachingStyle,
        ExplanationDetail: userConfig.Coach.ExplanationDetail,
    }
    if err := s.db.UpsertCoachSettings(coachSettings); err != nil {
        log.Printf("Failed to save coach settings: %v", err)
        data := s.NewSettingsData(user, userConfig)
        data.Error = fmt.Sprintf("Failed to save coach settings: %v", err)
        s.templates.ExecuteTemplate(w, "settings.html", data)
        return
    }

    // 5. Save home equipment - clear and re-insert
    if err := s.db.ClearUserEquipment(user.ID); err != nil {
        log.Printf("Failed to clear equipment: %v", err)
    }
    for _, item := range userConfig.Equipment.Home {
        if err := s.db.AddEquipment(user.ID, "home", item); err != nil {
            log.Printf("Failed to save equipment '%s': %v", item, err)
        }
    }

    // 6. Save gyms and club sessions
    if err := s.db.ClearUserGyms(user.ID); err != nil {
        log.Printf("Failed to clear gyms: %v", err)
    }
    for i, gym := range userConfig.Equipment.Gyms {
        // Create gym
        dbGym := &db.Gym{
            UserID:     user.ID,
            Name:       gym.Name,
            Type:       gym.Type,
            Membership: gym.Membership,
        }

        // Parse sport_id if provided
        if i < len(gymSportIDs) && gymSportIDs[i] != "" && gymSportIDs[i] != "0" {
            if sportID, err := strconv.Atoi(gymSportIDs[i]); err == nil && sportID > 0 {
                dbGym.SportID = &sportID
            }
        }

        // Parse sessions limit fields if membership is "limited"
        if gym.Membership == "limited" {
            gymSessionsLimitSlice := r.Form["gym_sessions_limit[]"]
            gymLimitPeriodSlice := r.Form["gym_limit_period[]"]

            if i < len(gymSessionsLimitSlice) && gymSessionsLimitSlice[i] != "" {
                if limit, err := strconv.Atoi(gymSessionsLimitSlice[i]); err == nil {
                    dbGym.SessionsLimit = &limit
                }
            }

            if i < len(gymLimitPeriodSlice) && gymLimitPeriodSlice[i] != "" {
                period := gymLimitPeriodSlice[i]
                dbGym.LimitPeriod = &period
            }
        }

        gymID, err := s.db.CreateGym(dbGym)
        if err != nil {
            log.Printf("Failed to create gym '%s': %v", gym.Name, err)
            continue
        }

        // Create club sessions for this gym
        for _, session := range gym.Sessions {
            // Parse day of week and time from Occurrences field
            // Occurrences format: "Monday 7pm and 8pm, Thursday 12:30pm"
            // or "Thursday 12:30, Thursday 17:30"
            // Duration is separate now: "60 min"

            var durationMinutes int
            // Parse duration from Duration field
            if session.Duration != "" {
                durationStr := strings.TrimSpace(session.Duration)
                durationStr = strings.TrimSuffix(durationStr, " min")
                durationStr = strings.TrimSpace(durationStr)
                fmt.Sscanf(durationStr, "%d", &durationMinutes)
            }

            // Split by comma to handle multiple day/time pairs
            occurrencePairs := strings.Split(session.Occurrences, ",")

            for _, pair := range occurrencePairs {
                pair = strings.TrimSpace(pair)
                if pair == "" {
                    continue
                }

                // Extract day and time from each pair (e.g., "Thursday 12:30")
                dayTimeParts := strings.Fields(pair)

                var dayOfWeek, time string
                if len(dayTimeParts) >= 2 {
                    dayOfWeek = dayTimeParts[0]
                    time = dayTimeParts[1]
                } else {
                    continue // Skip invalid entries
                }

                gymIDInt := int(gymID)
                dbSession := &db.ClubSession{
                    UserID:          user.ID,
                    GymID:           &gymIDInt,
                    SportID:         dbGym.SportID, // Link session to same sport as gym
                    SessionName:     session.Name,
                    Description:     session.Description,
                    Occurrences:     session.Occurrences,
                    Cost:            session.Cost,
                    DayOfWeek:       dayOfWeek,
                    Time:            time,
                    DurationMinutes: durationMinutes,
                    SessionType:     "club",
                    Active:          true,
                }

                if _, err := s.db.CreateClubSession(dbSession); err != nil {
                    log.Printf("Failed to create club session '%s' for %s %s: %v", session.Name, dayOfWeek, time, err)
                }
            }
        }
    }

    // Redirect back to settings with success message
    data := s.NewSettingsData(user, userConfig)
    data.Success = true
    if err := s.templates.ExecuteTemplate(w, "settings.html", data); err != nil {
        http.Error(w, "Failed to render settings", http.StatusInternalServerError)
        log.Printf("Template error: %v", err)
    }
}

// validatePlanGenerationSettings checks if sufficient configuration exists for plan generation
func validatePlanGenerationSettings(userConfig *config.UserConfig) []string {
    var errors []string

    // Check primary sport
    if len(userConfig.Sports) == 0 || userConfig.Sports[0].Name == "" {
        errors = append(errors, "Primary sport must be selected")
    }

    // Check at least one goal
    if len(userConfig.Goals.ShortTerm) == 0 && len(userConfig.Goals.MediumTerm) == 0 && len(userConfig.Goals.LongTerm) == 0 {
        errors = append(errors, "At least one goal must be defined (short-term, medium-term, or long-term)")
    }

    // Check availability
    hasAvailability := false
    if userConfig.Availability != nil {
        for _, day := range userConfig.Availability {
            if day.Morning || day.Lunch || day.Evening {
                hasAvailability = true
                break
            }
        }
    }
    if !hasAvailability {
        errors = append(errors, "Weekly availability must be set for at least one time slot")
    }

    // Check API key
    if os.Getenv("ANTHROPIC_API_KEY") == "" {
        errors = append(errors, "ANTHROPIC_API_KEY environment variable is not set")
    }

    return errors
}

// HandleGeneratePlan generates a new weekly training plan
func (s *Server) HandleGeneratePlan(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    // Load user config
    userConfig, err := s.GetUserConfig(user.ID)
    if err != nil {
        http.Error(w, "Failed to load user config", http.StatusInternalServerError)
        log.Printf("Error loading user config: %v", err)
        return
    }

    // Validate settings
    validationErrors := validatePlanGenerationSettings(userConfig)
    if len(validationErrors) > 0 {
        // Show settings page with validation errors
        errorMsg := "Cannot generate plan. Please complete your profile first:\n"
        for _, err := range validationErrors {
            errorMsg += "• " + err + "\n"
        }
        log.Printf("Plan generation validation failed for user %s (ID: %d): %v", user.Name, user.ID, validationErrors)
        data := s.NewSettingsData(user, userConfig)
        data.Error = errorMsg
        s.templates.ExecuteTemplate(w, "settings.html", data)
        return
    }

    // Check rate limit
    withinLimit, currentCount, err := s.db.CheckAIRateLimit(user.ID, DailyAICallLimit)
    if err != nil {
        http.Error(w, "Failed to check rate limit", http.StatusInternalServerError)
        log.Printf("Error checking rate limit: %v", err)
        return
    }
    if !withinLimit {
        data := s.NewSettingsData(user, userConfig)
        data.Error = fmt.Sprintf("Daily AI call limit reached (%d/%d). Please try again tomorrow.", currentCount, DailyAICallLimit)
        s.templates.ExecuteTemplate(w, "settings.html", data)
        return
    }

    // Parse form to get special constraints
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    specialConstraints := r.FormValue("plan_constraints")

    // Generate plan
    weekStart := agent.GetWeekStart(time.Now())
    planText, err := agent.GenerateWeeklyPlan(r.Context(), s.db, userConfig, weekStart, specialConstraints)
    if err != nil {
        log.Printf("Failed to generate plan: %v", err)
        data := s.NewSettingsData(user, userConfig)
        data.Error = fmt.Sprintf("Failed to generate plan: %v", err)
        s.templates.ExecuteTemplate(w, "settings.html", data)
        return
    }

    // Save plan
    if err := agent.SaveWeeklyPlan(s.db, user.ID, weekStart, planText); err != nil {
        log.Printf("Failed to save plan: %v", err)
        data := s.NewSettingsData(user, userConfig)
        data.Error = fmt.Sprintf("Failed to save plan: %v", err)
        s.templates.ExecuteTemplate(w, "settings.html", data)
        return
    }

    log.Printf("Successfully generated and saved weekly plan for user: %s (ID: %d)", user.Name, user.ID)

    // Increment AI usage counter
    if err := s.db.IncrementAIUsage(user.ID, time.Now()); err != nil {
        log.Printf("Warning: Failed to increment AI usage: %v", err)
    }

    // Redirect to dashboard to show new plan
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleLogRestDay handles logging a rest day note
func (s *Server) HandleLogRestDay(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    restDate := r.FormValue("rest_date")
    wellnessStr := r.FormValue("wellness_rating")
    sorenessStr := r.FormValue("soreness_level")
    motivationStr := r.FormValue("motivation_level")
    recoveryActivities := r.FormValue("recovery_activities")
    notes := r.FormValue("notes")

    date, err := time.Parse("2006-01-02", restDate)
    if err != nil {
        http.Error(w, "Invalid date format", http.StatusBadRequest)
        return
    }

    wellness, err := strconv.Atoi(wellnessStr)
    if err != nil || wellness < 1 || wellness > 10 {
        http.Error(w, "Invalid wellness rating (must be 1-10)", http.StatusBadRequest)
        return
    }

    soreness, err := strconv.Atoi(sorenessStr)
    if err != nil || soreness < 1 || soreness > 10 {
        http.Error(w, "Invalid soreness level (must be 1-10)", http.StatusBadRequest)
        return
    }

    motivation, err := strconv.Atoi(motivationStr)
    if err != nil || motivation < 1 || motivation > 10 {
        http.Error(w, "Invalid motivation level (must be 1-10)", http.StatusBadRequest)
        return
    }

    restDayNote := &db.RestDayNote{
        UserID:             user.ID,
        RestDate:           date,
        WellnessRating:     wellness,
        SorenessLevel:      soreness,
        MotivationLevel:    motivation,
        RecoveryActivities: recoveryActivities,
        Notes:              notes,
    }

    if err := s.db.CreateRestDayNote(restDayNote); err != nil {
        http.Error(w, "Failed to save rest day note", http.StatusInternalServerError)
        log.Printf("Error saving rest day note: %v", err)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleCopyPreviousWeek copies the previous week's plan to the current week
func (s *Server) HandleCopyPreviousWeek(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    adjustWithAI := r.FormValue("adjust_with_ai") == "true"
    adjustmentNotes := r.FormValue("adjustment_notes")

    // Get current week start
    currentWeekStart := agent.GetWeekStart(time.Now())

    // Get previous week start (7 days before current week)
    previousWeekStart := currentWeekStart.AddDate(0, 0, -7)

    // Get previous week's plan
    previousPlan, err := s.db.GetWeeklyPlan(user.ID, previousWeekStart)
    if err != nil {
        http.Error(w, "No plan found for previous week", http.StatusNotFound)
        log.Printf("Error loading previous week plan: %v", err)
        return
    }

    var newPlanData string
    var newRationale string

    if adjustWithAI && adjustmentNotes != "" {
        // Load user config for AI adjustment
        userConfig, err := s.GetUserConfig(user.ID)
        if err != nil {
            http.Error(w, "Failed to load user config", http.StatusInternalServerError)
            log.Printf("Error loading user config: %v", err)
            return
        }

        // Check rate limit for AI usage
        withinLimit, currentCount, err := s.db.CheckAIRateLimit(user.ID, DailyAICallLimit)
        if err != nil {
            http.Error(w, "Failed to check rate limit", http.StatusInternalServerError)
            log.Printf("Error checking rate limit: %v", err)
            return
        }
        if !withinLimit {
            data := s.NewSettingsData(user, userConfig)
            data.Error = fmt.Sprintf("Daily AI call limit reached (%d/%d). Please try again tomorrow or copy without AI adjustment.", currentCount, DailyAICallLimit)
            s.templates.ExecuteTemplate(w, "settings.html", data)
            return
        }

        // Generate adjusted plan using AI
        ctx := r.Context()
        adjustedPlan, err := agent.AdjustWeeklyPlan(
            ctx,
            s.db,
            userConfig,
            currentWeekStart,
            previousPlan.PlanData,
            adjustmentNotes,
        )
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to adjust plan: %v", err), http.StatusInternalServerError)
            log.Printf("Error adjusting plan: %v", err)
            return
        }

        // Increment AI usage counter
        if err := s.db.IncrementAIUsage(user.ID, time.Now()); err != nil {
            log.Printf("Warning: Failed to increment AI usage: %v", err)
        }

        newPlanData = adjustedPlan
        newRationale = fmt.Sprintf("Adjusted from previous week with notes: %s", adjustmentNotes)
    } else {
        // Simple copy without AI adjustment
        newPlanData = previousPlan.PlanData
        newRationale = "Copied from previous week"
    }

    // Create new plan for current week
    currentWeekEnd := currentWeekStart.AddDate(0, 0, 6)
    newPlan := &db.WeeklyPlan{
        UserID:        user.ID,
        WeekStartDate: currentWeekStart,
        WeekEndDate:   currentWeekEnd,
        PlanData:      newPlanData,
        Rationale:     newRationale,
    }

    if err := s.db.CreateWeeklyPlan(newPlan); err != nil {
        http.Error(w, "Failed to save copied plan", http.StatusInternalServerError)
        log.Printf("Error saving copied plan: %v", err)
        return
    }

    log.Printf("Successfully copied plan from previous week for user: %s (ID: %d)", user.Name, user.ID)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// StatsData holds data for the stats template
type StatsData struct {
    User             *db.User
    WeeklyStats      *db.WeeklyStats
    MonthlyStats     *db.MonthlyStats
    ActivityStats    []*db.ActivityStats
    StreakStats      *db.StreakStats
    VolumeData       []*db.VolumeDataPoint
    EffortData       []*db.EffortDataPoint
    ThemeClass       string
    SportIcon        string
}

// HandleStats displays the stats page
func (s *Server) HandleStats(w http.ResponseWriter, r *http.Request) {
    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    // Get week start for current week
    weekStart := agent.GetWeekStart(time.Now())

    // Get weekly stats
    weeklyStats, err := s.db.GetWeeklyStats(user.ID, weekStart)
    if err != nil {
        log.Printf("Error loading weekly stats: %v", err)
        weeklyStats = nil
    }

    // Get monthly stats (last 30 days)
    monthlyStats, err := s.db.GetMonthlyStats(user.ID, 30)
    if err != nil {
        log.Printf("Error loading monthly stats: %v", err)
        monthlyStats = nil
    }

    // Get activity breakdown (last 30 days)
    activityStats, err := s.db.GetActivityStats(user.ID, 30)
    if err != nil {
        log.Printf("Error loading activity stats: %v", err)
        activityStats = []*db.ActivityStats{}
    }

    // Get streak stats
    streakStats, err := s.db.GetStreakStats(user.ID)
    if err != nil {
        log.Printf("Error loading streak stats: %v", err)
        streakStats = nil
    }

    // Get volume data for chart (last 30 days)
    volumeData, err := s.db.GetVolumeDataPoints(user.ID, 30)
    if err != nil {
        log.Printf("Error loading volume data: %v", err)
        volumeData = []*db.VolumeDataPoint{}
    }

    // Get effort data for chart (last 30 days)
    effortData, err := s.db.GetEffortDataPoints(user.ID, 30)
    if err != nil {
        log.Printf("Error loading effort data: %v", err)
        effortData = []*db.EffortDataPoint{}
    }

    data := StatsData{
        User:          user,
        WeeklyStats:   weeklyStats,
        MonthlyStats:  monthlyStats,
        ActivityStats: activityStats,
        StreakStats:   streakStats,
        VolumeData:    volumeData,
        EffortData:    effortData,
        ThemeClass:    s.GetThemeClass(user.ID),
        SportIcon:     s.GetSportIconForUser(user.ID),
    }

    if err := s.templates.ExecuteTemplate(w, "stats.html", data); err != nil {
        http.Error(w, "Failed to render stats", http.StatusInternalServerError)
        log.Printf("Template error: %v", err)
    }
}
