package main

import (
    "fmt"
    "os"
    "time"

    "github.com/tuxnam/iamfeel/internal/config"
    "github.com/tuxnam/iamfeel/internal/db"
)

const (
    defaultDBPath         = "data/coach.db"
    defaultUserConfigPath = "data/user_config.yaml"
)

func runOnboarding() error {
    fmt.Println("╔══════════════════════════════════════════════════════════╗")
    fmt.Println("║              Welcome to IamFeel Onboarding!              ║")
    fmt.Println("║                                                          ║")
    fmt.Println("║  Let's set up your personalized training assistant.     ║")
    fmt.Println("║  This will take about 5-10 minutes.                      ║")
    fmt.Println("╚══════════════════════════════════════════════════════════╝")
    fmt.Println()

    // Check if user config already exists
    if _, err := os.Stat(defaultUserConfigPath); err == nil {
        fmt.Println("⚠️  User configuration already exists!")
        if !askYesNo("Do you want to overwrite it?", false) {
            fmt.Println("Onboarding cancelled.")
            return nil
        }
    }

    userConfig := &config.UserConfig{}

    // Step 1: Basic Information
    if err := collectBasicInfo(userConfig); err != nil {
        return err
    }

    // Step 2: Sport Selection
    if err := collectSportInfo(userConfig); err != nil {
        return err
    }

    // Step 3: Equipment & Gym Access (includes club sessions)
    if err := collectEquipmentInfo(userConfig); err != nil {
        return err
    }

    // Step 4: Weekly Availability
    if err := collectAvailability(userConfig); err != nil {
        return err
    }

    // Step 5: Goals
    if err := collectGoals(userConfig); err != nil {
        return err
    }

    // Step 6: Nutrition
    if err := collectNutrition(userConfig); err != nil {
        return err
    }

    // Step 7: Preferences
    if err := collectPreferences(userConfig); err != nil {
        return err
    }

    // Step 8: Coach Settings
    if err := collectCoachSettings(userConfig); err != nil {
        return err
    }

    // Step 9: Tracking Settings
    collectTrackingSettings(userConfig)

    // Save configuration
    fmt.Println()
    fmt.Println("💾 Saving your configuration...")
    if err := config.SaveUserConfig(defaultUserConfigPath, userConfig); err != nil {
        return fmt.Errorf("failed to save configuration: %w", err)
    }
    fmt.Printf("✓ Configuration saved to %s\n", defaultUserConfigPath)

    // Initialize database
    fmt.Println()
    fmt.Println("🗄️  Setting up your database...")
    if err := seedDatabase(userConfig); err != nil {
        return fmt.Errorf("failed to initialize database: %w", err)
    }
    fmt.Printf("✓ Database initialized at %s\n", defaultDBPath)

    // Done!
    fmt.Println()
    fmt.Println("╔══════════════════════════════════════════════════════════╗")
    fmt.Println("║                 🎉 Onboarding Complete! 🎉               ║")
    fmt.Println("║                                                          ║")
    fmt.Println("║  Your IamFeel coach is ready to help you train!         ║")
    fmt.Println("║                                                          ║")
    fmt.Println("║  Next steps:                                             ║")
    fmt.Println("║    • Generate your first plan: iamfeel plan              ║")
    fmt.Println("║    • Start the web server: make run-server               ║")
    fmt.Println("╚══════════════════════════════════════════════════════════╝")
    fmt.Println()

    return nil
}

func collectBasicInfo(cfg *config.UserConfig) error {
    fmt.Println()
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("  STEP 1: Basic Information")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

    cfg.User.Name = askString("What's your name?", "")
    cfg.User.Age = askInt("What's your age?", 30)

    experienceLevel := askChoice(
        "What's your overall fitness experience level?",
        []string{"Beginner", "Intermediate", "Advanced"},
        1,
    )
    levels := []string{"beginner", "intermediate", "advanced"}
    cfg.User.ExperienceLevel = levels[experienceLevel]

    return nil
}

func collectSportInfo(cfg *config.UserConfig) error {
    fmt.Println()
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("  STEP 2: Sport Selection")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

    // For now, we'll default to boxing
    // In the future, we can scan the configs/sports directory
    sportChoice := askChoice(
        "Which sport do you primarily train?",
        []string{"Boxing"},
        0,
    )

    var sportName, configPath string
    switch sportChoice {
    case 0:
        sportName = "boxing"
        configPath = "configs/sports/boxing.yaml"
    }

    // Load sport config to get phases
    sportConfig, err := config.LoadSportConfig(configPath)
    if err != nil {
        return fmt.Errorf("failed to load sport config: %w", err)
    }

    experienceYears := askInt(fmt.Sprintf("How many years have you been training %s?", sportName), 1)

    // Ask about current phase
    phaseNames := make([]string, len(sportConfig.Phases))
    for i, phase := range sportConfig.Phases {
        phaseNames[i] = fmt.Sprintf("%s - %s", phase.DisplayName, phase.Focus)
    }
    phaseIdx := askChoice("What training phase are you currently in?", phaseNames, 0)
    currentPhase := sportConfig.Phases[phaseIdx].Name

    // Ask about phase dates
    fmt.Println()
    fmt.Println("When did this phase start? (YYYY-MM-DD)")
    phaseStart := askString("Phase start date", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))

    fmt.Println("When should this phase end? (YYYY-MM-DD)")
    phaseEnd := askString("Phase end date", time.Now().AddDate(0, 1, 0).Format("2006-01-02"))

    cfg.Sports = []config.UserSport{
        {
            Name:            sportName,
            ConfigFile:      configPath,
            Primary:         true,
            ExperienceYears: experienceYears,
            CurrentPhase:    currentPhase,
            PhaseStartDate:  phaseStart,
            PhaseEndDate:    phaseEnd,
        },
    }

    return nil
}

func collectEquipmentInfo(cfg *config.UserConfig) error {
    return collectEquipmentInfoWithRetry(cfg, 0)
}

func collectEquipmentInfoWithRetry(cfg *config.UserConfig, retryCount int) error {
    const maxRetries = 3

    if retryCount >= maxRetries {
        return fmt.Errorf("maximum retries exceeded for equipment info collection")
    }

    fmt.Println()
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("  STEP 3: Training Locations & Equipment")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println()
    fmt.Println("Where do you typically train? Select all that apply:")

    locationOptions := []string{
        "Home/Garage gym (own equipment)",
        "Gym(s) or Club(s) with membership",
        "Outdoors (running, cycling, calisthenics, etc.)",
    }

    selectedLocations := askMultiChoice("Select your training locations", locationOptions)

    if len(selectedLocations) == 0 {
        fmt.Println()
        fmt.Println("⚠️  You must select at least one training location!")
        return collectEquipmentInfoWithRetry(cfg, retryCount+1) // Retry with incremented counter
    }

    // Process home equipment
    hasHome := false
    for _, idx := range selectedLocations {
        if idx == 0 { // Home/Garage gym
            hasHome = true
            break
        }
    }

    if hasHome {
        fmt.Println()
        fmt.Println("── Home Equipment ──")
        cfg.Equipment.Home = askList("What equipment do you have at home?")
    }

    // Process gyms/clubs
    hasGyms := false
    for _, idx := range selectedLocations {
        if idx == 1 { // Gyms/Clubs
            hasGyms = true
            break
        }
    }

    if hasGyms {
        fmt.Println()
        fmt.Println("── Gym & Club Memberships ──")
        fmt.Println("Let's add your gym/club memberships and their available classes.")
        fmt.Println()

        for {
            gym := config.Gym{}

            gym.Name = askString("Gym/club name", "")
            if gym.Name == "" {
                break
            }

            gym.Type = askString("Type (e.g., boxing_club, commercial_gym, CrossFit, hockey_club)", "boxing_club")
            gym.Membership = askString("Membership type (e.g., unlimited, 3x/week, drop-in)", "unlimited")

            // Ask about scheduled sessions at this gym
            fmt.Println()
            hasSessions := askYesNo(fmt.Sprintf("Does %s have scheduled classes/sessions?", gym.Name), true)

            if hasSessions {
                fmt.Println()
                fmt.Printf("Let's add the sessions available at %s.\n", gym.Name)
                fmt.Println("For each type of session, you can specify multiple occurrences per week.")
                fmt.Println()

                gym.Sessions = collectSessionsForGym(gym.Name)
            }

            cfg.Equipment.Gyms = append(cfg.Equipment.Gyms, gym)

            fmt.Println()
            if !askYesNo("Add another gym/club?", false) {
                break
            }
            fmt.Println()
        }
    }

    // Process outdoors
    hasOutdoors := false
    for _, idx := range selectedLocations {
        if idx == 2 { // Outdoors
            hasOutdoors = true
            break
        }
    }

    if hasOutdoors {
        fmt.Println()
        fmt.Println("── Outdoor Training ──")
        fmt.Println("Note: Outdoor activities (running, cycling, calisthenics) are noted.")
        fmt.Println("You can track these as 'Outdoor Run', 'Outdoor Workout', etc. in your sessions.")
    }

    // Summary
    fmt.Println()
    fmt.Println("✓ Training locations configured:")
    if hasHome {
        fmt.Printf("  • Home gym with %d items\n", len(cfg.Equipment.Home))
    }
    if hasGyms {
        fmt.Printf("  • %d gym/club membership(s)\n", len(cfg.Equipment.Gyms))
    }
    if hasOutdoors {
        fmt.Println("  • Outdoor training")
    }

    return nil
}

func collectSessionsForGym(gymName string) []config.ClubSession {
    var sessions []config.ClubSession

    for {
        session := config.ClubSession{}

        session.Name = askString("Session name/type (e.g., 'Bag Work', 'Boxing Class')", "")
        if session.Name == "" {
            break
        }

        // Ask for description (what it is)
        fmt.Println()
        fmt.Println("Describe WHAT this session involves.")
        fmt.Println("Example: 'Technical drilling with partner work and combinations'")
        session.Description = askString("Description", "")

        // Ask for occurrences (when it happens)
        fmt.Println()
        fmt.Println("Specify WHEN this session occurs (days, times, duration).")
        fmt.Println("Example: 'Tuesdays & Thursdays 7pm, 60 min'")
        session.Occurrences = askString("Occurrences", "")

        // Ask for cost
        fmt.Println()
        session.Cost = askString("Cost (e.g., 'included', '$10', '$15/session')", "included")

        sessions = append(sessions, session)

        fmt.Println()
        if !askYesNo(fmt.Sprintf("Add another session type at %s?", gymName), false) {
            break
        }
        fmt.Println()
    }

    return sessions
}

func collectAvailability(cfg *config.UserConfig) error {
    fmt.Println()
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("  STEP 4: Weekly Availability")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println()
    fmt.Println("For each day, tell us when you're available to train:")
    fmt.Println("  Morning: 6-9am")
    fmt.Println("  Lunch:   12-2pm")
    fmt.Println("  Evening: 6-9pm")
    fmt.Println()

    cfg.Availability = make(map[string]config.DayAvailability)

    days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
    for _, day := range days {
        fmt.Printf("\n%s:\n", day)
        morning := askYesNo("  Available in the morning?", false)
        lunch := askYesNo("  Available at lunch?", false)
        evening := askYesNo("  Available in the evening?", true)

        var preferredLocation string
        if morning || lunch || evening {
            locationIdx := askChoice(
                "  Where do you prefer to workout this day?",
                []string{"Flexible (home or gym)", "Home", "Gym"},
                0,
            )
            locationOptions := []string{"flexible", "home", "gym"}
            preferredLocation = locationOptions[locationIdx]
        }

        notes := askString("  Any notes for this day? (optional)", "")

        cfg.Availability[day] = config.DayAvailability{
            Morning:           morning,
            Lunch:             lunch,
            Evening:           evening,
            PreferredLocation: preferredLocation,
            Notes:             notes,
        }
    }

    return nil
}

func collectGoals(cfg *config.UserConfig) error {
    fmt.Println()
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("  STEP 5: Goals")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

    fmt.Println()
    fmt.Println("Short-term goals (1-3 months):")
    cfg.Goals.ShortTerm = askList("Enter your short-term goals")

    fmt.Println()
    fmt.Println("Medium-term goals (3-6 months):")
    cfg.Goals.MediumTerm = askList("Enter your medium-term goals")

    fmt.Println()
    fmt.Println("Long-term goals (6-12 months):")
    cfg.Goals.LongTerm = askList("Enter your long-term goals")

    return nil
}

func collectNutrition(cfg *config.UserConfig) error {
    fmt.Println()
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("  STEP 6: Supplements")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

    takesSupplements := askYesNo("Do you take any supplements?", false)
    if takesSupplements {
        fmt.Println()
        fmt.Println("Let's add your supplements:")
        for {
            name := askString("Supplement name (or blank to finish)", "")
            if name == "" {
                break
            }
            dosage := askString("Dosage (e.g., '5g')", "")
            timing := askString("When do you take it? (e.g., 'daily', 'post-workout')", "")

            cfg.Supplements = append(cfg.Supplements, config.Supplement{
                Name:   name,
                Dosage: dosage,
                Timing: timing,
            })
        }
    }

    return nil
}

func collectPreferences(cfg *config.UserConfig) error {
    fmt.Println()
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("  STEP 7: Training Preferences")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

    intensityIdx := askChoice(
        "How intense do you like your training?",
        []string{"Conservative (focus on recovery)", "Moderate", "Moderate-to-high", "Aggressive (push hard)"},
        1,
    )
    intensities := []string{"conservative", "moderate", "moderate-to-high", "aggressive"}
    cfg.Preferences.IntensityPreference = intensities[intensityIdx]

    recoveryIdx := askChoice(
        "How do you prioritize recovery?",
        []string{"Aggressive (minimal rest)", "Balanced", "Conservative (plenty of rest)"},
        1,
    )
    recoveries := []string{"aggressive", "balanced", "conservative"}
    cfg.Preferences.RecoveryPriority = recoveries[recoveryIdx]

    planFreqIdx := askChoice(
        "How often do you want new plans generated?",
        []string{"Weekly", "Bi-weekly", "Monthly"},
        0,
    )
    frequencies := []string{"weekly", "bi-weekly", "monthly"}
    cfg.Preferences.PlanFrequency = frequencies[planFreqIdx]

    cfg.Preferences.SessionDurationPreference = askString("Preferred session duration (e.g., '60-75 minutes')", "60-75 minutes")

    return nil
}

func collectCoachSettings(cfg *config.UserConfig) error {
    fmt.Println()
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("  STEP 8: AI Coach Settings")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

    modelIdx := askChoice(
        "Which Claude model should the coach use?",
        []string{
            "claude-3-5-sonnet-20241022 (Recommended - balanced)",
            "claude-3-opus-20240229 (Most capable, higher cost)",
            "claude-3-haiku-20240307 (Fastest, lower cost)",
        },
        0,
    )
    models := []string{
        "claude-3-5-sonnet-20241022",
        "claude-3-opus-20240229",
        "claude-3-haiku-20240307",
    }
    cfg.Coach.Model = models[modelIdx]
    cfg.Coach.Temperature = 0.7

    styleIdx := askChoice(
        "What coaching style do you prefer?",
        []string{"Aggressive (push hard)", "Balanced", "Supportive (encouraging)"},
        1,
    )
    styles := []string{"aggressive", "balanced", "supportive"}
    cfg.Coach.CoachingStyle = styles[styleIdx]

    detailIdx := askChoice(
        "How detailed should explanations be?",
        []string{"Brief", "Moderate", "Detailed"},
        1,
    )
    details := []string{"brief", "moderate", "detailed"}
    cfg.Coach.ExplanationDetail = details[detailIdx]

    return nil
}

func collectTrackingSettings(cfg *config.UserConfig) {
    fmt.Println()
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("  STEP 9: Tracking Settings")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

    cfg.Tracking.HistoryMonths = 6
    cfg.Tracking.TrackSupplements = askYesNo("Track supplements?", len(cfg.Supplements) > 0)
    cfg.Tracking.TrackSleep = false // Future feature
    cfg.Tracking.TrackWeight = false // Future feature
}

func seedDatabase(cfg *config.UserConfig) error {
    // Create database
    database, err := db.New(defaultDBPath)
    if err != nil {
        return err
    }
    defer database.Close()

    // Create user
    user, err := database.CreateUser(cfg.User.Name, cfg.User.Age, cfg.User.ExperienceLevel)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    // Add goals
    for _, goal := range cfg.Goals.ShortTerm {
        if err := database.CreateGoal(&db.Goal{
            UserID:      user.ID,
            GoalType:    "short_term",
            Description: goal,
        }); err != nil {
            return fmt.Errorf("failed to create goal: %w", err)
        }
    }

    for _, goal := range cfg.Goals.MediumTerm {
        if err := database.CreateGoal(&db.Goal{
            UserID:      user.ID,
            GoalType:    "medium_term",
            Description: goal,
        }); err != nil {
            return fmt.Errorf("failed to create goal: %w", err)
        }
    }

    for _, goal := range cfg.Goals.LongTerm {
        if err := database.CreateGoal(&db.Goal{
            UserID:      user.ID,
            GoalType:    "long_term",
            Description: goal,
        }); err != nil {
            return fmt.Errorf("failed to create goal: %w", err)
        }
    }

    fmt.Printf("✓ Created user profile for %s (ID: %d)\n", user.Name, user.ID)
    fmt.Printf("✓ Added %d goals to database\n", len(cfg.Goals.ShortTerm)+len(cfg.Goals.MediumTerm)+len(cfg.Goals.LongTerm))

    return nil
}
