package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/tuxnam/iamfeel/internal/api"
    "github.com/tuxnam/iamfeel/internal/db"
)

const (
    defaultDBPath         = "data/coach.db"
    defaultUserConfigPath = "data/user_config.yaml"
    defaultPort           = "8080"
)

func main() {
    // Open database
    database, err := db.New(defaultDBPath)
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }
    defer database.Close()

    // Create server
    server := api.NewServer(database)

    // Setup router
    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Compress(5))
    r.Use(server.UserMiddleware) // User session management

    // Static files
    fileServer := http.FileServer(http.Dir("web/static"))
    r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

    // Routes
    r.Get("/", server.HandleDashboard)
    r.Get("/history", server.HandleHistory)
    r.Get("/templates", server.HandleTemplates)
    r.Post("/templates/create", server.HandleCreateTemplate)
    r.Post("/templates/delete", server.HandleDeleteTemplate)
    r.Post("/templates/use", server.HandleUseTemplate)
    r.Post("/templates/save", server.HandleSaveAsTemplate)
    r.Get("/settings", server.HandleSettings)
    r.Post("/settings", server.HandleSettingsSave)
    r.Post("/settings/generate-plan", server.HandleGeneratePlan)
    r.Post("/sessions/log", server.HandleLogSession)
    r.Post("/sessions/complete", server.HandleCompleteSession)
    r.Post("/sessions/skip", server.HandleSkipSession)
    r.Get("/plan/current", server.HandleCurrentPlan)
    r.Get("/plan/modify", server.HandleModifyDayForm)
    r.Post("/plan/modify", server.HandleModifyDayRequest)
    r.Post("/plan/modify/accept", server.HandleAcceptModification)

    // User management routes
    r.Get("/users", server.HandleUsers)
    r.Post("/users/create", server.HandleCreateUser)
    r.Post("/users/switch", server.HandleSwitchUser)

    // Calendar export routes
    r.Get("/calendar/export", server.HandleCalendarExport)
    r.Get("/calendar/export/session", server.HandleSessionCalendarExport)

    // Get port from environment or use default
    port := os.Getenv("PORT")
    if port == "" {
        port = defaultPort
    }

    // Start server
    addr := fmt.Sprintf(":%s", port)
    fmt.Println("╔══════════════════════════════════════════════════════════╗")
    fmt.Println("║              IamFeel - Web Dashboard                     ║")
    fmt.Println("╚══════════════════════════════════════════════════════════╝")
    fmt.Println()
    fmt.Printf("🚀 Server starting on http://localhost:%s\n", port)
    fmt.Println()
    fmt.Println("Available routes:")
    fmt.Println("  • Dashboard:  http://localhost:" + port)
    fmt.Println("  • History:    http://localhost:" + port + "/history")
    fmt.Println()
    fmt.Println("Press Ctrl+C to stop")
    fmt.Println()

    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
