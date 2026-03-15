package api

import (
    "fmt"
    "html/template"
    "path/filepath"
    "reflect"

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
        "deref": func(ptr interface{}) interface{} {
            if ptr == nil {
                return ""
            }
            val := reflect.ValueOf(ptr)
            if val.Kind() == reflect.Ptr {
                if val.IsNil() {
                    return ""
                }
                return val.Elem().Interface()
            }
            return ptr
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

    // Get user sports
    sports, err := s.db.GetUserSports(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get sports: %w", err)
    }

    // Convert sports to config format
    configSports := []config.UserSport{}
    for _, sport := range sports {
        configSport := config.UserSport{
            Name:                  sport.SportName,
            ConfigFile:            sport.ConfigPath,
            Primary:               sport.IsPrimary,
            ExperienceYears:       sport.ExperienceYears,
            CurrentPhase:          sport.CurrentPhase,
            GoalType:              sport.GoalType,
            Priority:              sport.Priority,
            TargetSessionsPerWeek: sport.TargetSessionsPerWeek,
            Notes:                 sport.Notes,
        }
        if sport.PhaseStartDate != nil {
            configSport.PhaseStartDate = sport.PhaseStartDate.Format("2006-01-02")
        }
        if sport.PhaseEndDate != nil {
            configSport.PhaseEndDate = sport.PhaseEndDate.Format("2006-01-02")
        }
        configSports = append(configSports, configSport)
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

        // Group sessions by name+description+duration+cost to combine multiple day/time entries
        sessionMap := make(map[string]*config.ClubSession)
        for _, session := range sessions {
            // Create key from name+description+duration+cost
            duration := fmt.Sprintf("%d min", session.DurationMinutes)
            key := fmt.Sprintf("%s|%s|%s|%s", session.SessionName, session.Description, duration, session.Cost)

            dayTime := fmt.Sprintf("%s %s", session.DayOfWeek, session.Time)

            if existing, found := sessionMap[key]; found {
                // Append to existing occurrences
                existing.Occurrences = existing.Occurrences + ", " + dayTime
            } else {
                // Create new session
                sessionMap[key] = &config.ClubSession{
                    Name:        session.SessionName,
                    Description: session.Description,
                    Occurrences: dayTime,
                    Duration:    duration,
                    Cost:        session.Cost,
                }
            }
        }

        // Convert map to slice
        clubSessions := []config.ClubSession{}
        for _, session := range sessionMap {
            clubSessions = append(clubSessions, *session)
        }

        configGyms = append(configGyms, config.Gym{
            Name:          gym.Name,
            Type:          gym.Type,
            Membership:    gym.Membership,
            SportID:       gym.SportID,
            SessionsLimit: gym.SessionsLimit,
            LimitPeriod:   gym.LimitPeriod,
            Sessions:      clubSessions,
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
        AllowShortSessions:        prefs.AllowShortSessions,
        MaxSessionsPerDay:         prefs.MaxSessionsPerDay,
        Notes:                     prefs.Notes,
    }

    // Get coach settings with defaults
    coach := config.CoachSettings{
        Model:             "claude-haiku-4-5",
        Temperature:       0.7,
        CoachingStyle:     "motivational",
        ExplanationDetail: "balanced",
    }
    coachSettings, err := s.db.GetCoachSettings(userID)
    if err == nil && coachSettings != nil {
        coach.Model = coachSettings.Model
        coach.Temperature = coachSettings.Temperature
        coach.CoachingStyle = coachSettings.CoachingStyle
        coach.ExplanationDetail = coachSettings.ExplanationDetail
    }

    // Get tracking settings with defaults
    tracking := config.TrackingSettings{
        HistoryMonths:    6,
        TrackSupplements: true,
        TrackSleep:       false,
        TrackWeight:      false,
    }
    trackingSettings, err := s.db.GetTrackingSettings(userID)
    if err == nil && trackingSettings != nil {
        tracking.HistoryMonths = trackingSettings.HistoryMonths
        tracking.TrackSupplements = trackingSettings.TrackSupplements
        tracking.TrackSleep = trackingSettings.TrackSleep
        tracking.TrackWeight = trackingSettings.TrackWeight
    }

    // Build complete UserConfig
    userConfig := &config.UserConfig{
        User: config.UserProfile{
            Name:            user.Name,
            Age:             user.Age,
            ExperienceLevel: user.ExperienceLevel,
        },
        Sports:       configSports,
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
