package api

import (
    "fmt"
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
    // Create template functions
    funcMap := template.FuncMap{
        "Iterate": func(n int) []int {
            result := make([]int, n)
            for i := 0; i < n; i++ {
                result[i] = i + 1
            }
            return result
        },
        "printf": func(format string, args ...interface{}) string {
            return fmt.Sprintf(format, args...)
        },
    }

    // Load templates with custom functions
    templates := template.Must(template.New("").Funcs(funcMap).ParseGlob(filepath.Join("web", "templates", "*.html")))

    return &Server{
        db:        database,
        templates: templates,
    }
}

// GetUserConfig loads the config for a specific user from the database
func (s *Server) GetUserConfig(userID int) (*config.UserConfig, error) {
    // Get user basic info
    user, err := s.db.GetUser(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    // Get equipment
    equipment, err := s.db.GetUserEquipment(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get equipment: %w", err)
    }

    // Convert equipment to slice format
    equipmentAccess := config.EquipmentAccess{
        Home: []string{},
    }
    for _, eq := range equipment {
        if eq.Location == "home" {
            equipmentAccess.Home = append(equipmentAccess.Home, eq.EquipmentName)
        }
    }

    // Get gyms
    gyms, err := s.db.GetUserGyms(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get gyms: %w", err)
    }

    // Convert gyms to config format
    configGyms := []config.Gym{}
    for _, gym := range gyms {
        // Get club sessions for this gym
        sessions, err := s.db.GetGymSessions(gym.ID)
        if err != nil {
            return nil, fmt.Errorf("failed to get gym sessions: %w", err)
        }

        clubSessions := []config.ClubSession{}
        for _, session := range sessions {
            // Format occurrences from database fields
            occurrences := fmt.Sprintf("%s %s, %d min", session.DayOfWeek, session.Time, session.DurationMinutes)
            clubSessions = append(clubSessions, config.ClubSession{
                Name:        session.SessionName,
                Description: session.Description,
                Occurrences: occurrences,
                Cost:        session.CostType,
            })
        }

        configGyms = append(configGyms, config.Gym{
            Name:     gym.Name,
            Type:     gym.Type,
            Sessions: clubSessions,
        })
    }
    equipmentAccess.Gyms = configGyms

    // Get availability
    availList, err := s.db.GetUserAvailability(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get availability: %w", err)
    }

    availability := make(map[string]config.DayAvailability)
    for _, avail := range availList {
        availability[avail.DayOfWeek] = config.DayAvailability{
            Morning: avail.Morning,
            Lunch:   avail.Lunch,
            Evening: avail.Evening,
        }
    }

    // Get goals
    goals, err := s.db.GetUserGoals(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get goals: %w", err)
    }

    configGoals := config.Goals{
        ShortTerm:  []string{},
        MediumTerm: []string{},
        LongTerm:   []string{},
    }
    for _, goal := range goals {
        switch goal.GoalType {
        case "short_term":
            configGoals.ShortTerm = append(configGoals.ShortTerm, goal.Description)
        case "medium_term":
            configGoals.MediumTerm = append(configGoals.MediumTerm, goal.Description)
        case "long_term":
            configGoals.LongTerm = append(configGoals.LongTerm, goal.Description)
        }
    }

    // Get supplements
    supplements, err := s.db.GetUserSupplements(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get supplements: %w", err)
    }

    configSupplements := []config.Supplement{}
    for _, supp := range supplements {
        configSupplements = append(configSupplements, config.Supplement{
            Name:   supp.Name,
            Dosage: supp.Dosage,
            Timing: supp.Timing,
        })
    }

    // Get user preferences
    prefs, err := s.db.GetUserPreferences(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user preferences: %w", err)
    }

    preferences := config.UserPreferences{
        PrimaryGoal:               prefs.PrimaryGoal,
        SessionsPerWeek:           prefs.SessionsPerWeek,
        PreferredDuration:         prefs.PreferredDuration,
        SessionDurationPreference: prefs.SessionDurationPreference,
        IntensityPreference:       prefs.IntensityPreference,
        RecoveryPriority:          prefs.RecoveryPriority,
        PlanFrequency:             prefs.PlanFrequency,
        Notes:                     prefs.Notes,
    }

    // Get coach settings
    coachSettings, err := s.db.GetCoachSettings(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get coach settings: %w", err)
    }

    coach := config.CoachSettings{
        Model:             coachSettings.Model,
        Temperature:       coachSettings.Temperature,
        CoachingStyle:     coachSettings.CoachingStyle,
        ExplanationDetail: coachSettings.ExplanationDetail,
    }

    // Get tracking settings
    trackingSettings, err := s.db.GetTrackingSettings(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get tracking settings: %w", err)
    }

    tracking := config.TrackingSettings{
        HistoryMonths:    trackingSettings.HistoryMonths,
        TrackSupplements: trackingSettings.TrackSupplements,
        TrackSleep:       trackingSettings.TrackSleep,
        TrackWeight:      trackingSettings.TrackWeight,
    }

    // Build complete UserConfig
    userConfig := &config.UserConfig{
        User: config.UserProfile{
            Name:            user.Name,
            Age:             user.Age,
            ExperienceLevel: user.ExperienceLevel,
        },
        Equipment:    equipmentAccess,
        Availability: availability,
        Goals:        configGoals,
        Supplements:  configSupplements,
        Preferences:  preferences,
        Coach:        coach,
        Tracking:     tracking,
    }

    return userConfig, nil
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
