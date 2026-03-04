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
    Name                string
    Age                 int
    ExperienceLevel     string
    SportName            string
    SportExperienceYears int

    // Current phase
    CurrentPhaseName string
    PhaseStartDate   string
    PhaseEndDate     string

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
    IntensityPreference       string
    RecoveryPriority          string
    SessionDurationPreference string

    // Special constraints from user
    SpecialConstraints string

    // Plan timing
    WeekStart string
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

    // Load sport info
    primarySport := userConfig.GetPrimarySport()
    if primarySport == nil {
        return nil, fmt.Errorf("no primary sport configured")
    }

    ctx.SportName = primarySport.Name
    ctx.SportExperienceYears = primarySport.ExperienceYears
    ctx.CurrentPhaseName = primarySport.CurrentPhase
    ctx.PhaseStartDate = primarySport.PhaseStartDate
    ctx.PhaseEndDate = primarySport.PhaseEndDate

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
    ctx.IntensityPreference = userConfig.Preferences.IntensityPreference
    ctx.RecoveryPriority = userConfig.Preferences.RecoveryPriority
    ctx.SessionDurationPreference = userConfig.Preferences.SessionDurationPreference

    return ctx, nil
}
