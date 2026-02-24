package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/tuxnam/iamfeel/internal/agent"
    "github.com/tuxnam/iamfeel/internal/config"
    "github.com/tuxnam/iamfeel/internal/db"
)

func runPlanGeneration() error {
    fmt.Println("╔══════════════════════════════════════════════════════════╗")
    fmt.Println("║            IamFeel - Weekly Plan Generation             ║")
    fmt.Println("╚══════════════════════════════════════════════════════════╝")
    fmt.Println()

    // Check if user config exists
    if _, err := os.Stat(defaultUserConfigPath); os.IsNotExist(err) {
        return fmt.Errorf("user configuration not found. Please run 'iamfeel onboard' first")
    }

    // Check if ANTHROPIC_API_KEY is set
    demoMode := os.Getenv("ANTHROPIC_API_KEY") == ""
    if demoMode {
        fmt.Println("⚠️  DEMO MODE: ANTHROPIC_API_KEY not set")
        fmt.Println("   Generating sample plan for testing purposes.")
        fmt.Println("   For personalized AI plans, get an API key from:")
        fmt.Println("   https://console.anthropic.com/")
        fmt.Println()
    }

    // Load user config
    fmt.Println("📖 Loading your profile...")
    userConfig, err := config.LoadUserConfig(defaultUserConfigPath)
    if err != nil {
        return fmt.Errorf("failed to load user config: %w", err)
    }
    fmt.Printf("✓ Loaded profile for %s\n", userConfig.User.Name)

    // Open database
    fmt.Println("🗄️  Connecting to database...")
    database, err := db.New(defaultDBPath)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    defer database.Close()
    fmt.Println("✓ Database connected")

    // Get user
    user, err := database.GetFirstUser()
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }

    // Determine week to plan for
    now := time.Now()
    weekStart := agent.GetWeekStart(now)

    fmt.Println()
    fmt.Printf("Generating plan for week of %s to %s\n",
        weekStart.Format("Mon Jan 2"),
        weekStart.AddDate(0, 0, 6).Format("Mon Jan 2, 2006"))
    fmt.Println()

    // Check if plan already exists
    existingPlan, _ := database.GetWeeklyPlan(user.ID, weekStart)
    if existingPlan != nil {
        fmt.Println("⚠️  A plan already exists for this week!")
        if !askYesNo("Do you want to regenerate it?", false) {
            fmt.Println()
            fmt.Println("Showing existing plan:")
            fmt.Println()
            fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
            fmt.Println(existingPlan.PlanData)
            fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
            return nil
        }
        fmt.Println()
    }

    // Generate plan
    var planText string
    if demoMode {
        fmt.Println("📝 Generating demo training plan...")
        fmt.Println()
        planText = agent.GenerateDemoPlan(user.Name, weekStart)
    } else {
        fmt.Println("🤖 Generating your personalized training plan...")
        fmt.Println("   (This may take 10-30 seconds)")
        fmt.Println()

        ctx := context.Background()
        var err error
        planText, err = agent.GenerateWeeklyPlan(ctx, database, userConfig, weekStart)
        if err != nil {
            return fmt.Errorf("failed to generate plan: %w", err)
        }
    }

    // Save plan to database
    fmt.Println("💾 Saving plan to database...")
    if err := agent.SaveWeeklyPlan(database, user.ID, weekStart, planText); err != nil {
        return fmt.Errorf("failed to save plan: %w", err)
    }
    fmt.Println("✓ Plan saved")
    fmt.Println()

    // Display plan
    fmt.Println("╔══════════════════════════════════════════════════════════╗")
    fmt.Println("║                   YOUR WEEKLY PLAN                       ║")
    fmt.Println("╚══════════════════════════════════════════════════════════╝")
    fmt.Println()
    fmt.Println(planText)
    fmt.Println()
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println()
    fmt.Println("💡 Next steps:")
    fmt.Println("   • View your plan anytime on the dashboard: make run-server")
    fmt.Println("   • Log your workouts as you complete them")
    fmt.Println("   • Generate next week's plan: iamfeel plan")
    fmt.Println()

    return nil
}
