package agent

import (
    "context"
    "fmt"
    "time"

    "github.com/tuxnam/iamfeel/internal/config"
    "github.com/tuxnam/iamfeel/internal/db"
)

// GenerateWeeklyPlan generates a training plan for the given week
func GenerateWeeklyPlan(
    ctx context.Context,
    database *db.DB,
    userConfig *config.UserConfig,
    weekStart time.Time,
) (string, error) {
    // Load sport config
    primarySport := userConfig.GetPrimarySport()
    if primarySport == nil {
        return "", fmt.Errorf("no primary sport configured")
    }

    sportConfig, err := config.LoadSportConfig(primarySport.ConfigFile)
    if err != nil {
        return "", fmt.Errorf("failed to load sport config: %w", err)
    }

    // Create AI client
    client, err := NewClient(userConfig.Coach.Model, userConfig.Coach.Temperature)
    if err != nil {
        return "", fmt.Errorf("failed to create AI client: %w", err)
    }

    // Build system prompt
    systemPrompt := BuildSystemPrompt(
        sportConfig,
        userConfig.Coach.CoachingStyle,
        userConfig.Coach.ExplanationDetail,
    )

    // Load user context
    userContext, err := LoadUserContext(database, userConfig, weekStart)
    if err != nil {
        return "", fmt.Errorf("failed to load user context: %w", err)
    }

    // Build user prompt
    userPrompt := BuildUserPrompt(userContext)

    // Generate plan
    plan, err := client.GenerateCompletion(ctx, systemPrompt, userPrompt)
    if err != nil {
        return "", fmt.Errorf("failed to generate plan: %w", err)
    }

    return plan, nil
}

// SaveWeeklyPlan saves the generated plan to the database and creates planned sessions
func SaveWeeklyPlan(database *db.DB, userID int, weekStart time.Time, planText string) error {
    weekEnd := weekStart.AddDate(0, 0, 6)

    plan := &db.WeeklyPlan{
        UserID:        userID,
        WeekStartDate: weekStart,
        WeekEndDate:   weekEnd,
        PlanData:      planText,
        Rationale:     extractRationale(planText),
    }

    if err := database.CreateWeeklyPlan(plan); err != nil {
        return fmt.Errorf("failed to save weekly plan: %w", err)
    }

    // Create planned sessions from the plan
    if err := CreatePlannedSessions(database, userID, weekStart, planText); err != nil {
        // Don't fail the whole save if planned sessions fail
        // Just log and continue
        fmt.Printf("Warning: failed to create planned sessions: %v\n", err)
    }

    return nil
}

// extractRationale attempts to extract the rationale section from the plan
func extractRationale(planText string) string {
    // Simple extraction - look for RATIONALE: section
    // This is basic and could be improved with better parsing
    const rationaleMarker = "RATIONALE:"
    const nextSection = "\n\nWEEKLY OVERVIEW:"

    startIdx := -1
    for i := 0; i < len(planText)-len(rationaleMarker); i++ {
        if planText[i:i+len(rationaleMarker)] == rationaleMarker {
            startIdx = i + len(rationaleMarker)
            break
        }
    }

    if startIdx == -1 {
        return ""
    }

    endIdx := len(planText)
    for i := startIdx; i < len(planText)-len(nextSection); i++ {
        if planText[i:i+len(nextSection)] == nextSection {
            endIdx = i
            break
        }
    }

    rationale := planText[startIdx:endIdx]
    // Trim whitespace
    for len(rationale) > 0 && (rationale[0] == ' ' || rationale[0] == '\n') {
        rationale = rationale[1:]
    }
    for len(rationale) > 0 && (rationale[len(rationale)-1] == ' ' || rationale[len(rationale)-1] == '\n') {
        rationale = rationale[:len(rationale)-1]
    }

    return rationale
}

// GenerateDayModification generates a modified plan for a single day
func GenerateDayModification(
    ctx context.Context,
    database *db.DB,
    userConfig *config.UserConfig,
    weekStart time.Time,
    currentPlan string,
    dayName string,
    modificationReason string,
) (string, error) {
    // Load sport config
    primarySport := userConfig.GetPrimarySport()
    if primarySport == nil {
        return "", fmt.Errorf("no primary sport configured")
    }

    sportConfig, err := config.LoadSportConfig(primarySport.ConfigFile)
    if err != nil {
        return "", fmt.Errorf("failed to load sport config: %w", err)
    }

    // Create AI client
    client, err := NewClient(userConfig.Coach.Model, userConfig.Coach.Temperature)
    if err != nil {
        return "", fmt.Errorf("failed to create AI client: %w", err)
    }

    // Build system prompt
    systemPrompt := BuildSystemPrompt(
        sportConfig,
        userConfig.Coach.CoachingStyle,
        userConfig.Coach.ExplanationDetail,
    )

    // Load user context
    userContext, err := LoadUserContext(database, userConfig, weekStart)
    if err != nil {
        return "", fmt.Errorf("failed to load user context: %w", err)
    }

    // Build day modification prompt
    modificationPrompt := BuildDayModificationPrompt(userContext, currentPlan, dayName, modificationReason)

    // Generate modified day
    newDay, err := client.GenerateCompletion(ctx, systemPrompt, modificationPrompt)
    if err != nil {
        return "", fmt.Errorf("failed to generate day modification: %w", err)
    }

    return newDay, nil
}

// AdjustWeeklyPlan adjusts a weekly plan based on user notes
func AdjustWeeklyPlan(
    ctx context.Context,
    database *db.DB,
    userConfig *config.UserConfig,
    weekStart time.Time,
    previousPlan string,
    adjustmentNotes string,
) (string, error) {
    // Load sport config
    primarySport := userConfig.GetPrimarySport()
    if primarySport == nil {
        return "", fmt.Errorf("no primary sport configured")
    }

    sportConfig, err := config.LoadSportConfig(primarySport.ConfigFile)
    if err != nil {
        return "", fmt.Errorf("failed to load sport config: %w", err)
    }

    // Create AI client
    client, err := NewClient(userConfig.Coach.Model, userConfig.Coach.Temperature)
    if err != nil {
        return "", fmt.Errorf("failed to create AI client: %w", err)
    }

    // Build system prompt
    systemPrompt := BuildSystemPrompt(
        sportConfig,
        userConfig.Coach.CoachingStyle,
        userConfig.Coach.ExplanationDetail,
    )

    // Load user context
    userContext, err := LoadUserContext(database, userConfig, weekStart)
    if err != nil {
        return "", fmt.Errorf("failed to load user context: %w", err)
    }

    // Build adjustment prompt
    adjustmentPrompt := BuildPlanAdjustmentPrompt(userContext, previousPlan, adjustmentNotes)

    // Generate adjusted plan
    adjustedPlan, err := client.GenerateCompletion(ctx, systemPrompt, adjustmentPrompt)
    if err != nil {
        return "", fmt.Errorf("failed to generate adjusted plan: %w", err)
    }

    return adjustedPlan, nil
}

// GetWeekStart returns the start of the week (Monday) for the given date
func GetWeekStart(date time.Time) time.Time {
    weekday := int(date.Weekday())
    if weekday == 0 { // Sunday
        weekday = 7
    }
    daysToMonday := weekday - 1
    return date.AddDate(0, 0, -daysToMonday).Truncate(24 * time.Hour)
}
