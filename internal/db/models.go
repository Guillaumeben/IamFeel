package db

import "time"

// User represents the core user profile
type User struct {
    ID              int       `json:"id"`
    Name            string    `json:"name"`
    Age             int       `json:"age"`
    Weight          float64   `json:"weight"`
    Height          float64   `json:"height"`
    ExperienceLevel string    `json:"experience_level"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// UserSport represents a sport the user practices
type UserSport struct {
    ID               int       `json:"id"`
    UserID           int       `json:"user_id"`
    SportName        string    `json:"sport_name"`
    ConfigPath       string    `json:"config_path"`
    IsPrimary        bool      `json:"is_primary"`
    ExperienceYears  int       `json:"experience_years"`
    CurrentPhase     string    `json:"current_phase"`
    PhaseStartDate   *time.Time `json:"phase_start_date"`
    PhaseEndDate     *time.Time `json:"phase_end_date"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

// TrainingSession represents a logged workout
type TrainingSession struct {
    ID               int       `json:"id"`
    UserID           int       `json:"user_id"`
    SportID          *int      `json:"sport_id"`
    SessionDate      time.Time `json:"session_date"`
    SessionType      string    `json:"session_type"`
    DurationMinutes  int       `json:"duration_minutes"`
    PerceivedEffort  int       `json:"perceived_effort"`
    Notes            string    `json:"notes"`
    PerformanceNotes string    `json:"performance_notes"` // Weights, reps, PBs, etc.
    Skipped          bool      `json:"skipped"`
    SkipReason       string    `json:"skip_reason"`
    Completed        bool      `json:"completed"`
    Planned          bool      `json:"planned"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

// WeeklyPlan represents an AI-generated training plan
type WeeklyPlan struct {
    ID            int       `json:"id"`
    UserID        int       `json:"user_id"`
    WeekStartDate time.Time `json:"week_start_date"`
    WeekEndDate   time.Time `json:"week_end_date"`
    PlanData      string    `json:"plan_data"` // JSON string
    Rationale     string    `json:"rationale"`
    GeneratedAt   time.Time `json:"generated_at"`
    CreatedAt     time.Time `json:"created_at"`
}

// NutritionLog represents daily nutrition tracking
type NutritionLog struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    LogDate     time.Time `json:"log_date"`
    Calories    int       `json:"calories"`
    ProteinGrams int      `json:"protein_grams"`
    CarbsGrams  int       `json:"carbs_grams"`
    FatGrams    int       `json:"fat_grams"`
    MealsData   string    `json:"meals_data"` // JSON string
    Notes       string    `json:"notes"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// SupplementLog represents supplement tracking
type SupplementLog struct {
    ID             int       `json:"id"`
    UserID         int       `json:"user_id"`
    LogDate        time.Time `json:"log_date"`
    SupplementName string    `json:"supplement_name"`
    Dosage         string    `json:"dosage"`
    Timing         string    `json:"timing"`
    Taken          bool      `json:"taken"`
    CreatedAt      time.Time `json:"created_at"`
}

// Goal represents a user goal
type Goal struct {
    ID          int        `json:"id"`
    UserID      int        `json:"user_id"`
    GoalType    string     `json:"goal_type"`
    Description string     `json:"description"`
    TargetDate  *time.Time `json:"target_date"`
    Completed   bool       `json:"completed"`
    CompletedAt *time.Time `json:"completed_at"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

// Equipment represents available equipment
type Equipment struct {
    ID            int       `json:"id"`
    UserID        int       `json:"user_id"`
    Location      string    `json:"location"`
    EquipmentName string    `json:"equipment_name"`
    CreatedAt     time.Time `json:"created_at"`
}

// ClubSession represents a scheduled club/gym session
type ClubSession struct {
    ID              int       `json:"id"`
    UserID          int       `json:"user_id"`
    GymID           *int      `json:"gym_id"`
    SportID         *int      `json:"sport_id"`
    SessionName     string    `json:"session_name"`
    Description     string    `json:"description"`
    Occurrences     string    `json:"occurrences"`
    Cost            string    `json:"cost"`
    DayOfWeek       string    `json:"day_of_week"`
    Time            string    `json:"time"`
    DurationMinutes int       `json:"duration_minutes"`
    SessionType     string    `json:"session_type"`
    CostType        string    `json:"cost_type"`
    CostAmount      float64   `json:"cost_amount"`
    Notes           string    `json:"notes"`
    Active          bool      `json:"active"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// Availability represents weekly time availability
type Availability struct {
    ID                int       `json:"id"`
    UserID            int       `json:"user_id"`
    DayOfWeek         string    `json:"day_of_week"`
    Morning           bool      `json:"morning"`
    Lunch             bool      `json:"lunch"`
    Evening           bool      `json:"evening"`
    PreferredLocation string    `json:"preferred_location"`
    Notes             string    `json:"notes"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}

// ChatMessage represents a chat message for context
type ChatMessage struct {
    ID             int       `json:"id"`
    UserID         int       `json:"user_id"`
    MessageRole    string    `json:"message_role"`
    MessageContent string    `json:"message_content"`
    CreatedAt      time.Time `json:"created_at"`
}

// Gym represents a gym or club membership
type Gym struct {
    ID             int       `json:"id"`
    UserID         int       `json:"user_id"`
    Name           string    `json:"name"`
    Type           string    `json:"type"`
    Membership     string    `json:"membership"`
    AvailableDays  string    `json:"available_days"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

// SupplementDefinition represents a supplement that user takes regularly
type SupplementDefinition struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Name      string    `json:"name"`
    Dosage    string    `json:"dosage"`
    Timing    string    `json:"timing"`
    Active    bool      `json:"active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// UserPreferences represents user training preferences
type UserPreferences struct {
    ID                        int       `json:"id"`
    UserID                    int       `json:"user_id"`
    PrimaryGoal               string    `json:"primary_goal"`
    SessionsPerWeek           int       `json:"sessions_per_week"`
    PreferredDuration         int       `json:"preferred_duration"`
    PreferredSessionTimes     string    `json:"preferred_session_times"`
    SessionDurationPreference string    `json:"session_duration_preference"`
    IntensityPreference       string    `json:"intensity_preference"`
    RecoveryPriority          string    `json:"recovery_priority"`
    PlanFrequency             string    `json:"plan_frequency"`
    Notes                     string    `json:"notes"`
    CreatedAt                 time.Time `json:"created_at"`
    UpdatedAt                 time.Time `json:"updated_at"`
}

// FitnessBaseline represents baseline fitness metrics
type FitnessBaseline struct {
    ID                      int       `json:"id"`
    UserID                  int       `json:"user_id"`
    RestingHeartRate        int       `json:"resting_heart_rate"`
    VO2MaxEstimate          int       `json:"vo2_max_estimate"`
    Squat1RM                int       `json:"squat_1rm"`
    Deadlift1RM             int       `json:"deadlift_1rm"`
    Bench1RM                int       `json:"bench_1rm"`
    OverheadPress1RM        int       `json:"overhead_press_1rm"`
    MaxRoundsHeavyBag       int       `json:"max_rounds_heavy_bag"`
    MaxRoundsSparring       int       `json:"max_rounds_sparring"`
    ComfortableSparringPace string    `json:"comfortable_sparring_pace"`
    CardioNotes             string    `json:"cardio_notes"`
    StrengthNotes           string    `json:"strength_notes"`
    BoxingNotes             string    `json:"boxing_notes"`
    CreatedAt               time.Time `json:"created_at"`
    UpdatedAt               time.Time `json:"updated_at"`
}

// CoachSettings represents AI coach configuration
type CoachSettings struct {
    ID                int       `json:"id"`
    UserID            int       `json:"user_id"`
    Model             string    `json:"model"`
    Temperature       float64   `json:"temperature"`
    CoachingStyle     string    `json:"coaching_style"`
    ExplanationDetail string    `json:"explanation_detail"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}

// TrackingSettings represents tracking preferences
type TrackingSettings struct {
    ID               int       `json:"id"`
    UserID           int       `json:"user_id"`
    HistoryMonths    int       `json:"history_months"`
    TrackSupplements bool      `json:"track_supplements"`
    TrackSleep       bool      `json:"track_sleep"`
    TrackWeight      bool      `json:"track_weight"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

// SessionTemplate represents a reusable session configuration
type SessionTemplate struct {
    ID              int       `json:"id"`
    UserID          int       `json:"user_id"`
    TemplateName    string    `json:"template_name"`
    SportName       string    `json:"sport_name"`
    SessionType     string    `json:"session_type"`
    DurationMinutes int       `json:"duration_minutes"`
    PerceivedEffort int       `json:"perceived_effort"`
    Description     string    `json:"description"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// RestDayNote represents a rest day entry with wellness tracking
type RestDayNote struct {
    ID                 int       `json:"id"`
    UserID             int       `json:"user_id"`
    RestDate           time.Time `json:"rest_date"`
    WellnessRating     int       `json:"wellness_rating"`
    SorenessLevel      int       `json:"soreness_level"`
    MotivationLevel    int       `json:"motivation_level"`
    RecoveryActivities string    `json:"recovery_activities"`
    Notes              string    `json:"notes"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
}

// AIUsage represents daily AI call tracking for rate limiting
type AIUsage struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    UsageDate time.Time `json:"usage_date"`
    CallCount int       `json:"call_count"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
