package config

// SportConfig represents the complete sport configuration
type SportConfig struct {
    SportName   string                 `yaml:"sport_name"`
    SportType   string                 `yaml:"sport_type"`
    Phases      []Phase                `yaml:"phases"`
    SessionTypes []SessionType         `yaml:"session_types"`
    Equipment   []string               `yaml:"equipment"`
    TypicalWeeklyStructure WeeklyStructure `yaml:"typical_weekly_structure"`
    RecoveryGuidelines map[string]Recovery `yaml:"recovery_guidelines"`
    Considerations []string              `yaml:"considerations"`
    Nutrition   *NutritionGuidelines   `yaml:"nutrition,omitempty"`
    AgentContext string                `yaml:"agent_context"`
}

// Phase represents a training phase
type Phase struct {
    Name              string   `yaml:"name"`
    DisplayName       string   `yaml:"display_name"`
    TypicalDurationWeeks string `yaml:"typical_duration_weeks"`
    Focus             string   `yaml:"focus"`
    Priorities        []string `yaml:"priorities"`
}

// SessionType represents a type of training session
type SessionType struct {
    ID                   string `yaml:"id"`
    Name                 string `yaml:"name"`
    Description          string `yaml:"description"`
    TypicalDurationMinutes string `yaml:"typical_duration_minutes"`
    Intensity            string `yaml:"intensity"`
    PrimaryAdaptation    string `yaml:"primary_adaptation"`
    RecoveryImpact       string `yaml:"recovery_impact"`
    EquipmentNeeded      string `yaml:"equipment_needed,omitempty"`
    Notes                string `yaml:"notes,omitempty"`
}

// WeeklyStructure defines typical weekly training structure
type WeeklyStructure struct {
    Beginner     LevelStructure `yaml:"beginner"`
    Intermediate LevelStructure `yaml:"intermediate"`
    Advanced     LevelStructure `yaml:"advanced"`
}

// LevelStructure defines structure for a specific experience level
type LevelStructure struct {
    SessionsPerWeek int    `yaml:"sessions_per_week"`
    RestDays        int    `yaml:"rest_days"`
    Notes           string `yaml:"notes"`
}

// Recovery defines recovery requirements for different intensity levels
type Recovery struct {
    RestBeforeNextHardSession int    `yaml:"rest_before_next_hard_session"`
    RestBeforeSameType        int    `yaml:"rest_before_same_type"`
    Notes                     string `yaml:"notes,omitempty"`
}

// NutritionGuidelines provides nutrition guidance
type NutritionGuidelines struct {
    GeneralApproach string   `yaml:"general_approach"`
    KeySupplements  []string `yaml:"key_supplements"`
    HydrationFocus  string   `yaml:"hydration_focus"`
    PreTraining     string   `yaml:"pre_training,omitempty"`
    PostTraining    string   `yaml:"post_training,omitempty"`
}

// UserConfig represents the user's profile and preferences
type UserConfig struct {
    User         UserProfile      `yaml:"user"`
    Sports       []UserSport      `yaml:"sports"`
    Equipment    EquipmentAccess  `yaml:"equipment"`
    Availability map[string]DayAvailability `yaml:"availability"`
    Goals        Goals            `yaml:"goals"`
    Fitness      *FitnessBaseline `yaml:"fitness,omitempty"`
    Supplements  []Supplement     `yaml:"supplements,omitempty"`
    Preferences  UserPreferences  `yaml:"preferences"`
    Coach        CoachSettings    `yaml:"coach"`
    Tracking     TrackingSettings `yaml:"tracking"`
}

// UserProfile contains basic user information
type UserProfile struct {
    Name            string  `yaml:"name"`
    Age             int     `yaml:"age"`
    Weight          float64 `yaml:"weight,omitempty"`
    Height          float64 `yaml:"height,omitempty"`
    ExperienceLevel string  `yaml:"experience_level"`
}

// UserSport represents a sport the user practices
type UserSport struct {
    Name           string `yaml:"name"`
    ConfigFile     string `yaml:"config_file"`
    Primary        bool   `yaml:"primary"`
    ExperienceYears int   `yaml:"experience_years"`
    CurrentPhase   string `yaml:"current_phase"`
    PhaseStartDate string `yaml:"phase_start_date"`
    PhaseEndDate   string `yaml:"phase_end_date"`
}

// EquipmentAccess describes available equipment
type EquipmentAccess struct {
    Home []string `yaml:"home"`
    Gyms []Gym    `yaml:"gyms,omitempty"`
}

// Gym represents a gym or club membership
type Gym struct {
    Name           string         `yaml:"name"`
    Type           string         `yaml:"type"` // boxing_club, commercial_gym, CrossFit, etc.
    Membership     string         `yaml:"membership"`
    SessionsLimit  *int           `yaml:"sessions_limit,omitempty"`
    LimitPeriod    *string        `yaml:"limit_period,omitempty"`
    AvailableDays  []string       `yaml:"available_days,omitempty"`
    Sessions       []ClubSession  `yaml:"sessions"`
}

// ClubSession represents a type of session available at the gym/club
type ClubSession struct {
    Name        string `yaml:"name"`
    Description string `yaml:"description"`
    Occurrences string `yaml:"occurrences,omitempty"` // e.g., "Tuesdays & Thursdays 7pm, 60 min"
    Cost        string `yaml:"cost,omitempty"`        // e.g., "included", "$10", "$15/session"
}

// DayAvailability represents availability for a specific day
type DayAvailability struct {
    Morning           bool   `yaml:"morning"`
    Lunch             bool   `yaml:"lunch"`
    Evening           bool   `yaml:"evening"`
    PreferredLocation string `yaml:"preferred_location,omitempty"` // "home", "gym", "flexible"
    Notes             string `yaml:"notes,omitempty"`
}

// Goals contains user goals
type Goals struct {
    ShortTerm  []string `yaml:"short_term"`
    MediumTerm []string `yaml:"medium_term"`
    LongTerm   []string `yaml:"long_term"`
}

// FitnessBaseline contains fitness assessment data
type FitnessBaseline struct {
    Cardio         *CardioBaseline         `yaml:"cardio,omitempty"`
    Strength       *StrengthBaseline       `yaml:"strength,omitempty"`
    BoxingSpecific *BoxingSpecificBaseline `yaml:"boxing_specific,omitempty"`
}

// CardioBaseline contains cardio fitness metrics
type CardioBaseline struct {
    RestingHeartRate int    `yaml:"resting_heart_rate"`
    VO2MaxEstimate   *int   `yaml:"vo2_max_estimate,omitempty"`
    Notes            string `yaml:"notes,omitempty"`
}

// StrengthBaseline contains strength metrics
type StrengthBaseline struct {
    Squat1RM          int    `yaml:"squat_1rm"`
    Deadlift1RM       int    `yaml:"deadlift_1rm"`
    Bench1RM          int    `yaml:"bench_1rm"`
    OverheadPress1RM  int    `yaml:"overhead_press_1rm"`
    Notes             string `yaml:"notes,omitempty"`
}

// BoxingSpecificBaseline contains boxing-specific metrics
type BoxingSpecificBaseline struct {
    MaxRoundsHeavyBag        int    `yaml:"max_rounds_heavy_bag"`
    MaxRoundsSparring        int    `yaml:"max_rounds_sparring"`
    ComfortableSparringPace  string `yaml:"comfortable_sparring_pace"`
    Notes                    string `yaml:"notes,omitempty"`
}

// Supplement represents a supplement
type Supplement struct {
    Name   string `yaml:"name"`
    Dosage string `yaml:"dosage"`
    Timing string `yaml:"timing"`
}

// UserPreferences contains training preferences
type UserPreferences struct {
    PrimaryGoal                string   `yaml:"primary_goal,omitempty"`
    SessionsPerWeek            int      `yaml:"sessions_per_week,omitempty"`
    PreferredDuration          int      `yaml:"preferred_duration,omitempty"`
    PreferredSessionTimes      []string `yaml:"preferred_session_times"`
    SessionDurationPreference  string   `yaml:"session_duration_preference"`
    IntensityPreference        string   `yaml:"intensity_preference"`
    RecoveryPriority           string   `yaml:"recovery_priority"`
    PlanFrequency              string   `yaml:"plan_frequency"`
    Notes                      string   `yaml:"notes,omitempty"`
}

// CoachSettings contains AI coach settings
type CoachSettings struct {
    Model             string  `yaml:"model"`
    Temperature       float64 `yaml:"temperature"`
    CoachingStyle     string  `yaml:"coaching_style"`
    ExplanationDetail string  `yaml:"explanation_detail"`
}

// TrackingSettings contains tracking preferences
type TrackingSettings struct {
    HistoryMonths    int  `yaml:"history_months"`
    TrackSupplements bool `yaml:"track_supplements"`
    TrackSleep       bool `yaml:"track_sleep"`
    TrackWeight      bool `yaml:"track_weight"`
}
