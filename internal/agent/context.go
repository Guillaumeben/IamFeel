package agent

import (
    "fmt"
    "time"

    "github.com/tuxnam/iamfeel/internal/config"
    "github.com/tuxnam/iamfeel/internal/db"
)

// UserContext contains all context needed for plan generation
type UserContext struct {
    // User profile
    Name            string
    Age             int
    ExperienceLevel string

    // Activities (supports multi-sport)
    Activities []ActivityInfo

    // Deprecated single-sport fields (for backwards compatibility)
    SportName            string
    SportExperienceYears int
    CurrentPhaseName     string
    PhaseStartDate       string
    PhaseEndDate         string

    // Goals
    ShortTermGoals  []string
    MediumTermGoals []string
    LongTermGoals   []string

    // Recent training
    RecentSessions []SessionSummary

    // Scheduled sessions
    ClubSessions []ClubSessionInfo

    // Availability
    Availability map[string]DayAvailabilityInfo

    // Equipment
    AvailableEquipment []string

    // Preferences
    PrimaryGoal               string
    SessionsPerWeek           float64
    IntensityPreference       string
    RecoveryPriority          string
    SessionDurationPreference string
    AllowShortSessions        bool
    MaxSessionsPerDay         int

    // Special constraints from user
    SpecialConstraints string

    // Plan timing
    WeekStart string
}

// ActivityInfo contains information about an activity/sport the user practices
type ActivityInfo struct {
    Name                  string
    ExperienceYears       int
    CurrentPhase          string
    PhaseStartDate        string
    PhaseEndDate          string
    GoalType              string  // competition_prep, maintenance, learning, recreation
    Priority              string  // high, medium, low
    TargetSessionsPerWeek float64
    Notes                 string
}

// SessionSummary is a simplified session for context
type SessionSummary struct {
    Date            string
    SessionType     string
    DurationMinutes int
    PerceivedEffort int
    Notes           string
}

// ClubSessionInfo contains club session details
type ClubSessionInfo struct {
    GymName     string
    Name        string
    Description string
    Occurrences string
    Cost        string
}

// DayAvailabilityInfo contains availability for a day
type DayAvailabilityInfo struct {
    Morning           bool
    Lunch             bool
    Evening           bool
    PreferredLocation string
    Notes             string
}

// LoadUserContext loads all context needed for plan generation
func LoadUserContext(database *db.DB, userConfig *config.UserConfig, weekStart time.Time) (*UserContext, error) {
    ctx := &UserContext{
        WeekStart: weekStart.Format("2006-01-02"),
    }

    // Load user profile
    user, err := database.GetFirstUser()
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    ctx.Name = user.Name
    ctx.Age = user.Age
    ctx.ExperienceLevel = user.ExperienceLevel

    // Load activities (all sports/activities user practices)
    for _, sport := range userConfig.Sports {
        activity := ActivityInfo{
            Name:                  sport.Name,
            ExperienceYears:       sport.ExperienceYears,
            CurrentPhase:          sport.CurrentPhase,
            PhaseStartDate:        sport.PhaseStartDate,
            PhaseEndDate:          sport.PhaseEndDate,
            GoalType:              sport.GoalType,
            Priority:              sport.Priority,
            TargetSessionsPerWeek: sport.TargetSessionsPerWeek,
            Notes:                 sport.Notes,
        }

        // Set defaults if not specified
        if activity.GoalType == "" {
            if sport.Primary {
                activity.GoalType = "competition_prep"
            } else {
                activity.GoalType = "maintenance"
            }
        }
        if activity.Priority == "" {
            if sport.Primary {
                activity.Priority = "high"
            } else {
                activity.Priority = "medium"
            }
        }

        ctx.Activities = append(ctx.Activities, activity)
    }

    // Backwards compatibility: populate single-sport fields with primary/highest priority sport
    primarySport := userConfig.GetPrimarySport()
    if primarySport == nil && len(ctx.Activities) > 0 {
        // No primary sport - use first activity (sorted by priority)
        firstActivity := ctx.Activities[0]
        ctx.SportName = firstActivity.Name
        ctx.SportExperienceYears = firstActivity.ExperienceYears
        ctx.CurrentPhaseName = firstActivity.CurrentPhase
        ctx.PhaseStartDate = firstActivity.PhaseStartDate
        ctx.PhaseEndDate = firstActivity.PhaseEndDate
    } else if primarySport != nil {
        ctx.SportName = primarySport.Name
        ctx.SportExperienceYears = primarySport.ExperienceYears
        ctx.CurrentPhaseName = primarySport.CurrentPhase
        ctx.PhaseStartDate = primarySport.PhaseStartDate
        ctx.PhaseEndDate = primarySport.PhaseEndDate
    } else {
        return nil, fmt.Errorf("no activities configured")
    }

    // Load goals
    goals, err := database.GetAllGoals(user.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to get goals: %w", err)
    }

    for _, goal := range goals {
        switch goal.GoalType {
        case "short_term":
            ctx.ShortTermGoals = append(ctx.ShortTermGoals, goal.Description)
        case "medium_term":
            ctx.MediumTermGoals = append(ctx.MediumTermGoals, goal.Description)
        case "long_term":
            ctx.LongTermGoals = append(ctx.LongTermGoals, goal.Description)
        }
    }

    // Load recent training history (last 14 days)
    endDate := weekStart.AddDate(0, 0, -1) // Day before week starts
    startDate := endDate.AddDate(0, 0, -14)

    sessions, err := database.GetTrainingSessions(user.ID, startDate, endDate)
    if err != nil {
        return nil, fmt.Errorf("failed to get training sessions: %w", err)
    }

    for _, session := range sessions {
        ctx.RecentSessions = append(ctx.RecentSessions, SessionSummary{
            Date:            session.SessionDate.Format("2006-01-02"),
            SessionType:     session.SessionType,
            DurationMinutes: session.DurationMinutes,
            PerceivedEffort: session.PerceivedEffort,
            Notes:           session.Notes,
        })
    }

    // Load club sessions from all gyms
    for _, gym := range userConfig.Equipment.Gyms {
        for _, cs := range gym.Sessions {
            sessionInfo := ClubSessionInfo{
                GymName:     gym.Name,
                Name:        cs.Name,
                Description: cs.Description,
                Occurrences: cs.Occurrences,
                Cost:        cs.Cost,
            }
            ctx.ClubSessions = append(ctx.ClubSessions, sessionInfo)
        }
    }

    // Load availability
    ctx.Availability = make(map[string]DayAvailabilityInfo)
    for day, avail := range userConfig.Availability {
        ctx.Availability[day] = DayAvailabilityInfo{
            Morning:           avail.Morning,
            Lunch:             avail.Lunch,
            Evening:           avail.Evening,
            PreferredLocation: avail.PreferredLocation,
            Notes:             avail.Notes,
        }
    }

    // Load equipment
    ctx.AvailableEquipment = append(ctx.AvailableEquipment, userConfig.Equipment.Home...)
    for _, gym := range userConfig.Equipment.Gyms {
        ctx.AvailableEquipment = append(ctx.AvailableEquipment,
            fmt.Sprintf("Gym access: %s (%s)", gym.Name, gym.Type))
    }

    // Load preferences
    ctx.PrimaryGoal = userConfig.Preferences.PrimaryGoal
    ctx.SessionsPerWeek = userConfig.Preferences.SessionsPerWeek
    ctx.IntensityPreference = userConfig.Preferences.IntensityPreference
    ctx.RecoveryPriority = userConfig.Preferences.RecoveryPriority
    ctx.SessionDurationPreference = userConfig.Preferences.SessionDurationPreference
    ctx.AllowShortSessions = userConfig.Preferences.AllowShortSessions
    ctx.MaxSessionsPerDay = userConfig.Preferences.MaxSessionsPerDay

    return ctx, nil
}
