package db

const Schema = `
-- Users table: core user profile
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    age INTEGER,
    weight REAL,
    height REAL,
    experience_level TEXT CHECK(experience_level IN ('beginner', 'intermediate', 'advanced')) DEFAULT 'beginner',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sports configuration: tracks which sports/activities user practices
CREATE TABLE IF NOT EXISTS user_sports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    sport_name TEXT NOT NULL,
    config_path TEXT NOT NULL,
    is_primary BOOLEAN DEFAULT 0,
    experience_years INTEGER DEFAULT 0,
    current_phase TEXT NOT NULL,
    phase_start_date DATE,
    phase_end_date DATE,
    goal_type TEXT CHECK(goal_type IN ('competition_prep', 'maintenance', 'learning', 'recreation')) DEFAULT 'maintenance',
    priority TEXT CHECK(priority IN ('high', 'medium', 'low')) DEFAULT 'medium',
    target_sessions_per_week REAL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Training sessions: logged workouts
CREATE TABLE IF NOT EXISTS training_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    sport_id INTEGER,
    session_date DATE NOT NULL,
    session_type TEXT NOT NULL,
    duration_minutes INTEGER,
    perceived_effort INTEGER CHECK(perceived_effort = 0 OR perceived_effort BETWEEN 1 AND 10),
    notes TEXT,
    performance_notes TEXT,
    skipped BOOLEAN DEFAULT 0,
    skip_reason TEXT,
    completed BOOLEAN DEFAULT 1,
    planned BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sport_id) REFERENCES user_sports(id) ON DELETE SET NULL
);

-- Weekly plans: AI-generated training plans
CREATE TABLE IF NOT EXISTS weekly_plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    week_start_date DATE NOT NULL,
    week_end_date DATE NOT NULL,
    plan_data JSON NOT NULL,
    rationale TEXT,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, week_start_date)
);

-- Nutrition logs: daily nutrition tracking
CREATE TABLE IF NOT EXISTS nutrition_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    log_date DATE NOT NULL,
    calories INTEGER,
    protein_grams INTEGER,
    carbs_grams INTEGER,
    fat_grams INTEGER,
    meals_data JSON,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, log_date)
);

-- Supplement definitions: base supplement information
CREATE TABLE IF NOT EXISTS supplement_definitions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    dosage TEXT,
    timing TEXT,
    active BOOLEAN DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, name)
);

-- Supplement tracking: daily logs of whether supplements were taken
CREATE TABLE IF NOT EXISTS supplement_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    supplement_id INTEGER,
    log_date DATE NOT NULL,
    supplement_name TEXT NOT NULL,
    dosage TEXT,
    timing TEXT,
    taken BOOLEAN DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (supplement_id) REFERENCES supplement_definitions(id) ON DELETE SET NULL
);

-- User goals
CREATE TABLE IF NOT EXISTS goals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    goal_type TEXT CHECK(goal_type IN ('short_term', 'medium_term', 'long_term')) NOT NULL,
    description TEXT NOT NULL,
    target_date DATE,
    completed BOOLEAN DEFAULT 0,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Equipment availability
CREATE TABLE IF NOT EXISTS equipment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    location TEXT CHECK(location IN ('home', 'gym', 'club')) NOT NULL,
    equipment_name TEXT NOT NULL,
    sport_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sport_id) REFERENCES user_sports(id) ON DELETE SET NULL
);

-- Gyms: gym/club memberships
CREATE TABLE IF NOT EXISTS gyms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    membership TEXT,
    available_days TEXT,
    sport_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sport_id) REFERENCES user_sports(id) ON DELETE SET NULL
);

-- Club sessions: scheduled club/gym sessions
CREATE TABLE IF NOT EXISTS club_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    gym_id INTEGER,
    sport_id INTEGER,
    session_name TEXT NOT NULL,
    description TEXT,
    occurrences TEXT,
    cost TEXT,
    day_of_week TEXT CHECK(day_of_week IN ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday')),
    time TEXT,
    duration_minutes INTEGER,
    session_type TEXT,
    cost_type TEXT CHECK(cost_type IN ('included', 'per_session')) DEFAULT 'included',
    cost_amount REAL,
    notes TEXT,
    active BOOLEAN DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,
    FOREIGN KEY (sport_id) REFERENCES user_sports(id) ON DELETE SET NULL
);

-- Weekly availability
CREATE TABLE IF NOT EXISTS availability (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    day_of_week TEXT CHECK(day_of_week IN ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday')) NOT NULL,
    morning BOOLEAN DEFAULT 0,
    lunch BOOLEAN DEFAULT 0,
    evening BOOLEAN DEFAULT 0,
    preferred_location TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, day_of_week)
);

-- User preferences: training preferences and goals
CREATE TABLE IF NOT EXISTS user_preferences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    primary_goal TEXT,
    sessions_per_week REAL,
    preferred_duration INTEGER,
    preferred_session_times TEXT,
    session_duration_preference TEXT,
    intensity_preference TEXT,
    recovery_priority TEXT,
    plan_frequency TEXT,
    allow_short_sessions BOOLEAN DEFAULT 0,
    max_sessions_per_day INTEGER DEFAULT 1,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Fitness baseline: cardio, strength, and sport-specific metrics
CREATE TABLE IF NOT EXISTS fitness_baseline (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    resting_heart_rate INTEGER,
    vo2_max_estimate INTEGER,
    squat_1rm INTEGER,
    deadlift_1rm INTEGER,
    bench_1rm INTEGER,
    overhead_press_1rm INTEGER,
    max_rounds_heavy_bag INTEGER,
    max_rounds_sparring INTEGER,
    comfortable_sparring_pace TEXT,
    cardio_notes TEXT,
    strength_notes TEXT,
    boxing_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Coach settings: AI coach configuration
CREATE TABLE IF NOT EXISTS coach_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    model TEXT DEFAULT 'claude-haiku-4-5',
    temperature REAL DEFAULT 0.7,
    coaching_style TEXT DEFAULT 'balanced',
    explanation_detail TEXT DEFAULT 'moderate',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Tracking settings: what data to track
CREATE TABLE IF NOT EXISTS tracking_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    history_months INTEGER DEFAULT 6,
    track_supplements BOOLEAN DEFAULT 1,
    track_sleep BOOLEAN DEFAULT 0,
    track_weight BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Chat history (for context in conversations)
CREATE TABLE IF NOT EXISTS chat_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    message_role TEXT CHECK(message_role IN ('user', 'assistant')) NOT NULL,
    message_content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Session templates: reusable session configurations
CREATE TABLE IF NOT EXISTS session_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    template_name TEXT NOT NULL,
    sport_name TEXT NOT NULL,
    session_type TEXT NOT NULL,
    duration_minutes INTEGER,
    perceived_effort INTEGER CHECK(perceived_effort = 0 OR perceived_effort BETWEEN 1 AND 10),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Rest day notes: logging for rest/recovery days
CREATE TABLE IF NOT EXISTS rest_day_notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    rest_date DATE NOT NULL,
    wellness_rating INTEGER CHECK(wellness_rating BETWEEN 1 AND 10),
    soreness_level INTEGER CHECK(soreness_level BETWEEN 1 AND 10),
    motivation_level INTEGER CHECK(motivation_level BETWEEN 1 AND 10),
    recovery_activities TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, rest_date)
);

-- AI usage tracking: rate limiting for AI calls
CREATE TABLE IF NOT EXISTS ai_usage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    usage_date DATE NOT NULL,
    call_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, usage_date)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_training_sessions_user_date ON training_sessions(user_id, session_date DESC);
CREATE INDEX IF NOT EXISTS idx_training_sessions_sport ON training_sessions(sport_id);
CREATE INDEX IF NOT EXISTS idx_weekly_plans_user_date ON weekly_plans(user_id, week_start_date DESC);
CREATE INDEX IF NOT EXISTS idx_nutrition_logs_user_date ON nutrition_logs(user_id, log_date DESC);
CREATE INDEX IF NOT EXISTS idx_goals_user_type ON goals(user_id, goal_type);
CREATE INDEX IF NOT EXISTS idx_chat_history_user ON chat_history(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_session_templates_user ON session_templates(user_id, sport_name);
CREATE INDEX IF NOT EXISTS idx_rest_day_notes_user_date ON rest_day_notes(user_id, rest_date DESC);
CREATE INDEX IF NOT EXISTS idx_ai_usage_user_date ON ai_usage(user_id, usage_date DESC);
CREATE INDEX IF NOT EXISTS idx_equipment_sport ON equipment(sport_id);
CREATE INDEX IF NOT EXISTS idx_user_sports_priority ON user_sports(user_id, priority);
CREATE INDEX IF NOT EXISTS idx_club_sessions_sport_active ON club_sessions(sport_id, active);
`

// migrationSQL contains ALTER statements for schema updates
const migrationSQL = `
-- Add new columns to training_sessions if they don't exist
-- SQLite doesn't have IF NOT EXISTS for ALTER TABLE, so we check via pragma
`
