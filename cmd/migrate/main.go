package main

import (
    "fmt"
    "log"
    "path/filepath"
    "strings"

    "github.com/tuxnam/iamfeel/internal/config"
    "github.com/tuxnam/iamfeel/internal/db"
)

const (
    defaultDBPath = "data/coach.db"
    dataDir       = "data"
)

func main() {
    log.Println("Starting YAML to Database migration...")

    // Open database
    database, err := db.New(defaultDBPath)
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }
    defer database.Close()

    // Find all user config YAML files
    files, err := filepath.Glob(filepath.Join(dataDir, "user_*_config.yaml"))
    if err != nil {
        log.Fatalf("Failed to find config files: %v", err)
    }

    if len(files) == 0 {
        log.Println("No user config files found to migrate")
        return
    }

    log.Printf("Found %d user config file(s) to migrate\n", len(files))

    // Migrate each user config file
    for _, file := range files {
        if err := migrateUserConfig(database, file); err != nil {
            log.Printf("ERROR migrating %s: %v\n", file, err)
        } else {
            log.Printf("✓ Successfully migrated %s\n", file)
        }
    }

    log.Println("Migration complete!")
}

func migrateUserConfig(database *db.DB, configPath string) error {
    // Extract user ID from filename (user_1_config.yaml -> 1)
    basename := filepath.Base(configPath)
    userIDStr := strings.TrimPrefix(basename, "user_")
    userIDStr = strings.TrimSuffix(userIDStr, "_config.yaml")

    var userID int
    if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
        return fmt.Errorf("failed to parse user ID from filename: %w", err)
    }

    // Check if user exists in database
    userName, err := database.GetUserName(userID)
    if err != nil {
        return fmt.Errorf("user %d not found in database, skipping", userID)
    }

    log.Printf("\nMigrating user %d (%s) from %s...", userID, userName, configPath)

    // Load YAML config
    cfg, err := config.LoadUserConfig(configPath)
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    // 1. Migrate Equipment
    log.Printf("  - Migrating equipment...")
    if err := database.ClearUserEquipment(userID); err != nil {
        return fmt.Errorf("failed to clear existing equipment: %w", err)
    }

    if cfg.Equipment.Home != nil {
        for _, item := range cfg.Equipment.Home {
            if err := database.AddEquipment(userID, "home", item); err != nil {
                return fmt.Errorf("failed to add home equipment '%s': %w", item, err)
            }
        }
        log.Printf("    ✓ Migrated %d home equipment items", len(cfg.Equipment.Home))
    }

    // Note: Gyms field in YAML contains gym memberships, not equipment
    // Equipment is only stored in Home field

    // 2. Migrate Gyms/Clubs (note: YAML doesn't store gyms, so skip this)
    log.Printf("  - Skipping gyms (not in YAML)")

    // 3. Migrate Availability
    log.Printf("  - Migrating availability...")
    availCount := 0
    for day, slots := range cfg.Availability {
        avail := &db.Availability{
            UserID:    userID,
            DayOfWeek: day,
            Morning:   slots.Morning,
            Lunch:     slots.Lunch,
            Evening:   slots.Evening,
        }
        if err := database.UpsertAvailability(avail); err != nil {
            return fmt.Errorf("failed to upsert availability for %s: %w", day, err)
        }
        availCount++
    }
    log.Printf("    ✓ Migrated %d availability entries", availCount)

    // 4. Migrate Goals (note: Goals table needs to be populated, but YAML structure is different)
    log.Printf("  - Skipping goals (need to migrate to new structure separately)")

    // 5. Migrate Supplements
    log.Printf("  - Migrating supplements...")
    if cfg.Supplements != nil {
        for _, supp := range cfg.Supplements {
            supplement := &db.SupplementDefinition{
                UserID:  userID,
                Name:    supp.Name,
                Dosage:  supp.Dosage,
                Timing:  supp.Timing,
            }
            if err := database.UpsertSupplement(supplement); err != nil {
                return fmt.Errorf("failed to upsert supplement '%s': %w", supp.Name, err)
            }
        }
        log.Printf("    ✓ Migrated %d supplements", len(cfg.Supplements))
    }

    // 6. Migrate User Preferences
    log.Printf("  - Migrating user preferences...")
    prefs := &db.UserPreferences{
        UserID:                     userID,
        PrimaryGoal:                cfg.Preferences.PrimaryGoal,
        SessionsPerWeek:            cfg.Preferences.SessionsPerWeek,
        PreferredDuration:          cfg.Preferences.PreferredDuration,
        SessionDurationPreference:  cfg.Preferences.SessionDurationPreference,
        IntensityPreference:        cfg.Preferences.IntensityPreference,
        RecoveryPriority:           cfg.Preferences.RecoveryPriority,
        PlanFrequency:              cfg.Preferences.PlanFrequency,
        Notes:                      cfg.Preferences.Notes,
    }
    if err := database.UpsertUserPreferences(prefs); err != nil {
        return fmt.Errorf("failed to upsert user preferences: %w", err)
    }
    log.Printf("    ✓ Migrated user preferences")

    // 7. Migrate Coach Settings
    log.Printf("  - Migrating coach settings...")
    coachSettings := &db.CoachSettings{
        UserID:           userID,
        Model:            cfg.Coach.Model,
        Temperature:      cfg.Coach.Temperature,
        CoachingStyle:    cfg.Coach.CoachingStyle,
        ExplanationDetail: cfg.Coach.ExplanationDetail,
    }
    if err := database.UpsertCoachSettings(coachSettings); err != nil {
        return fmt.Errorf("failed to upsert coach settings: %w", err)
    }
    log.Printf("    ✓ Migrated coach settings")

    // 8. Migrate Tracking Settings
    log.Printf("  - Migrating tracking settings...")
    trackingSettings := &db.TrackingSettings{
        UserID:           userID,
        HistoryMonths:    cfg.Tracking.HistoryMonths,
        TrackSupplements: cfg.Tracking.TrackSupplements,
        TrackSleep:       cfg.Tracking.TrackSleep,
        TrackWeight:      cfg.Tracking.TrackWeight,
    }
    if err := database.UpsertTrackingSettings(trackingSettings); err != nil {
        return fmt.Errorf("failed to upsert tracking settings: %w", err)
    }
    log.Printf("    ✓ Migrated tracking settings")

    return nil
}
