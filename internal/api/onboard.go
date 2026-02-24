package api

import (
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/tuxnam/iamfeel/internal/config"
    "github.com/tuxnam/iamfeel/internal/db"
)

// OnboardingData holds data for the onboarding template
type OnboardingData struct {
    User       *db.User
    Error      string
    ThemeClass string
}

// HandleOnboarding displays the onboarding form
func (s *Server) HandleOnboarding(w http.ResponseWriter, r *http.Request) {
    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    // Check if user already has a config
    if config.UserConfigExists(user.ID) {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    data := OnboardingData{
        User:       user,
        ThemeClass: s.GetThemeClass(user.ID),
    }

    if err := s.templates.ExecuteTemplate(w, "onboard.html", data); err != nil {
        http.Error(w, "Failed to render onboarding", http.StatusInternalServerError)
        log.Printf("Template error: %v", err)
    }
}

// HandleOnboardingSubmit processes the onboarding form
func (s *Server) HandleOnboardingSubmit(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    user, err := GetCurrentUser(r)
    if err != nil {
        http.Error(w, "Failed to load user", http.StatusInternalServerError)
        log.Printf("Error loading user: %v", err)
        return
    }

    // Check if config already exists
    if config.UserConfigExists(user.ID) {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Parse form values
    sportID := r.FormValue("sport")
    primaryGoal := r.FormValue("primary_goal")
    trainingLocation := r.FormValue("training_location")
    sessionsPerWeek := 3
    sessionDuration := 60
    trainingNotes := strings.TrimSpace(r.FormValue("training_notes"))

    if val := r.FormValue("sessions_per_week"); val != "" {
        fmt.Sscanf(val, "%d", &sessionsPerWeek)
    }
    if val := r.FormValue("session_duration"); val != "" {
        fmt.Sscanf(val, "%d", &sessionDuration)
    }

    // Parse equipment
    equipment := r.Form["equipment"]
    homeEquipment := []string{}
    for _, item := range equipment {
        if item != "" {
            homeEquipment = append(homeEquipment, item)
        }
    }

    // Parse availability
    availability := r.Form["availability"]
    availabilityMap := make(map[string]config.DayAvailability)
    days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
    for _, day := range days {
        isAvailable := false
        for _, selected := range availability {
            if selected == day {
                isAvailable = true
                break
            }
        }
        // Set default availability to evening for selected days
        availabilityMap[day] = config.DayAvailability{
            Morning:           false,
            Lunch:             false,
            Evening:           isAvailable,
            PreferredLocation: trainingLocation,
        }
    }

    // Create user config
    userConfig := &config.UserConfig{
        User: config.UserProfile{
            Name:            user.Name,
            Age:             user.Age,
            ExperienceLevel: user.ExperienceLevel,
        },
        Sports: []config.UserSport{
            {
                Name:    sportID,
                Primary: true,
            },
        },
        Equipment: config.EquipmentAccess{
            Home: homeEquipment,
            Gyms: []config.Gym{},
        },
        Availability: availabilityMap,
        Goals: config.Goals{
            ShortTerm:  []string{primaryGoal},
            MediumTerm: []string{},
            LongTerm:   []string{},
        },
        Fitness:     &config.FitnessBaseline{},
        Supplements: []config.Supplement{},
        Preferences: config.UserPreferences{
            PrimaryGoal:       primaryGoal,
            SessionsPerWeek:   sessionsPerWeek,
            PreferredDuration: sessionDuration,
            Notes:             trainingNotes,
        },
        Coach: config.CoachSettings{
            Model:             "claude-3-5-sonnet-20241022",
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

    // Save config
    if err := config.SaveUserConfigByID(user.ID, userConfig); err != nil {
        log.Printf("Failed to save config: %v", err)
        data := OnboardingData{
            User:       user,
            Error:      "Failed to save your configuration. Please try again.",
            ThemeClass: s.GetThemeClass(user.ID),
        }
        s.templates.ExecuteTemplate(w, "onboard.html", data)
        return
    }

    log.Printf("Onboarding completed for user: %s (ID: %d)", user.Name, user.ID)

    // Redirect to dashboard
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
