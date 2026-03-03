package db

import (
    "database/sql"
    "fmt"
    "time"
)

// ===== User Queries =====

// CreateUser creates a new user
func (db *DB) CreateUser(name string, age int, experienceLevel string) (*User, error) {
    query := `
        INSERT INTO users (name, age, experience_level)
        VALUES (?, ?, ?)
        RETURNING id, name, age, experience_level, created_at, updated_at
    `

    var user User
    err := db.conn.QueryRow(query, name, age, experienceLevel).Scan(
        &user.ID, &user.Name, &user.Age, &user.ExperienceLevel,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return &user, nil
}

// GetUser retrieves a user by ID
func (db *DB) GetUser(userID int) (*User, error) {
    query := `
        SELECT id, name, age, weight, height, experience_level, created_at, updated_at
        FROM users
        WHERE id = ?
    `

    var user User
    err := db.conn.QueryRow(query, userID).Scan(
        &user.ID, &user.Name, &user.Age, &user.Weight, &user.Height, &user.ExperienceLevel,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    return &user, nil
}

// GetFirstUser gets the first user (for single-user setup)
func (db *DB) GetFirstUser() (*User, error) {
    query := `
        SELECT id, name, age, weight, height, experience_level, created_at, updated_at
        FROM users
        ORDER BY id ASC
        LIMIT 1
    `

    var user User
    err := db.conn.QueryRow(query).Scan(
        &user.ID, &user.Name, &user.Age, &user.Weight, &user.Height, &user.ExperienceLevel,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("no users found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    return &user, nil
}

// UpdateUser updates an existing user
func (db *DB) UpdateUser(user *User) error {
    query := `
        UPDATE users
        SET name = ?, age = ?, weight = ?, height = ?, experience_level = ?, updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `

    _, err := db.conn.Exec(query, user.Name, user.Age, user.Weight, user.Height, user.ExperienceLevel, user.ID)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }

    return nil
}

// GetAllUsers retrieves all users from the database
func (db *DB) GetAllUsers() ([]*User, error) {
    query := `
        SELECT id, name, age, weight, height, experience_level, created_at, updated_at
        FROM users
        ORDER BY created_at ASC
    `

    rows, err := db.conn.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query users: %w", err)
    }
    defer rows.Close()

    var users []*User
    for rows.Next() {
        var user User
        err := rows.Scan(
            &user.ID, &user.Name, &user.Age, &user.Weight, &user.Height, &user.ExperienceLevel,
            &user.CreatedAt, &user.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, &user)
    }

    return users, nil
}

// ===== Training Session Queries =====

// CreateTrainingSession creates a new training session
func (db *DB) CreateTrainingSession(session *TrainingSession) error {
    query := `
        INSERT INTO training_sessions
        (user_id, sport_id, session_date, session_type, duration_minutes,
         perceived_effort, notes, performance_notes, skipped, skip_reason, completed, planned)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        RETURNING id, created_at, updated_at
    `

    return db.conn.QueryRow(
        query,
        session.UserID, session.SportID, session.SessionDate, session.SessionType,
        session.DurationMinutes, session.PerceivedEffort, session.Notes,
        session.PerformanceNotes, session.Skipped, session.SkipReason,
        session.Completed, session.Planned,
    ).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)
}

// GetTrainingSessions retrieves training sessions for a user within a date range
func (db *DB) GetTrainingSessions(userID int, startDate, endDate time.Time) ([]*TrainingSession, error) {
    query := `
        SELECT id, user_id, sport_id, session_date, session_type, duration_minutes,
               perceived_effort, notes, performance_notes, skipped, skip_reason, completed, planned, created_at, updated_at
        FROM training_sessions
        WHERE user_id = ? AND session_date BETWEEN ? AND ?
        ORDER BY session_date DESC
    `

    rows, err := db.conn.Query(query, userID, startDate, endDate)
    if err != nil {
        return nil, fmt.Errorf("failed to query training sessions: %w", err)
    }
    defer rows.Close()

    var sessions []*TrainingSession
    for rows.Next() {
        var session TrainingSession
        err := rows.Scan(
            &session.ID, &session.UserID, &session.SportID, &session.SessionDate,
            &session.SessionType, &session.DurationMinutes, &session.PerceivedEffort,
            &session.Notes, &session.PerformanceNotes, &session.Skipped, &session.SkipReason,
            &session.Completed, &session.Planned,
            &session.CreatedAt, &session.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan training session: %w", err)
        }
        sessions = append(sessions, &session)
    }

    return sessions, nil
}

// GetRecentTrainingSessions gets the most recent N training sessions
func (db *DB) GetRecentTrainingSessions(userID int, limit int) ([]*TrainingSession, error) {
    query := `
        SELECT id, user_id, sport_id, session_date, session_type, duration_minutes,
               perceived_effort, notes, performance_notes, skipped, skip_reason, completed, planned, created_at, updated_at
        FROM training_sessions
        WHERE user_id = ?
        ORDER BY session_date DESC, created_at DESC
        LIMIT ?
    `

    rows, err := db.conn.Query(query, userID, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to query recent training sessions: %w", err)
    }
    defer rows.Close()

    var sessions []*TrainingSession
    for rows.Next() {
        var session TrainingSession
        err := rows.Scan(
            &session.ID, &session.UserID, &session.SportID, &session.SessionDate,
            &session.SessionType, &session.DurationMinutes, &session.PerceivedEffort,
            &session.Notes, &session.PerformanceNotes, &session.Skipped, &session.SkipReason,
            &session.Completed, &session.Planned,
            &session.CreatedAt, &session.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan training session: %w", err)
        }
        sessions = append(sessions, &session)
    }

    return sessions, nil
}

// GetPastTrainingSessions returns training sessions before today
func (db *DB) GetPastTrainingSessions(userID int, limit int) ([]*TrainingSession, error) {
    query := `
        SELECT id, user_id, sport_id, session_date, session_type, duration_minutes,
               perceived_effort, notes, performance_notes, skipped, skip_reason, completed, planned, created_at, updated_at
        FROM training_sessions
        WHERE user_id = ? AND session_date < DATE('now')
        ORDER BY session_date DESC, created_at DESC
        LIMIT ?
    `

    rows, err := db.conn.Query(query, userID, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to query past training sessions: %w", err)
    }
    defer rows.Close()

    var sessions []*TrainingSession
    for rows.Next() {
        var session TrainingSession
        err := rows.Scan(
            &session.ID, &session.UserID, &session.SportID, &session.SessionDate,
            &session.SessionType, &session.DurationMinutes, &session.PerceivedEffort,
            &session.Notes, &session.PerformanceNotes, &session.Skipped, &session.SkipReason,
            &session.Completed, &session.Planned, &session.CreatedAt, &session.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan training session: %w", err)
        }
        sessions = append(sessions, &session)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating training sessions: %w", err)
    }

    return sessions, nil
}

// SessionFilters holds optional filter parameters for querying training sessions
type SessionFilters struct {
    StartDate    string // YYYY-MM-DD format
    EndDate      string // YYYY-MM-DD format
    SessionType  string
    MinEffort    int
    MaxEffort    int
}

// GetFilteredTrainingSessions returns training sessions filtered by the given criteria
func (db *DB) GetFilteredTrainingSessions(userID int, filters SessionFilters, limit int) ([]*TrainingSession, error) {
    query := `
        SELECT id, user_id, sport_id, session_date, session_type, duration_minutes,
               perceived_effort, notes, performance_notes, skipped, skip_reason, completed, planned, created_at, updated_at
        FROM training_sessions
        WHERE user_id = ?
    `

    args := []interface{}{userID}

    // Add date range filter
    if filters.StartDate != "" {
        query += " AND session_date >= ?"
        args = append(args, filters.StartDate)
    }
    if filters.EndDate != "" {
        query += " AND session_date <= ?"
        args = append(args, filters.EndDate)
    }

    // Add session type filter (case-insensitive)
    if filters.SessionType != "" {
        query += " AND LOWER(session_type) LIKE LOWER(?)"
        args = append(args, "%"+filters.SessionType+"%")
    }

    // Add effort level filter
    if filters.MinEffort > 0 {
        query += " AND perceived_effort >= ?"
        args = append(args, filters.MinEffort)
    }
    if filters.MaxEffort > 0 && filters.MaxEffort <= 10 {
        query += " AND perceived_effort <= ?"
        args = append(args, filters.MaxEffort)
    }

    query += `
        ORDER BY session_date DESC, created_at DESC
        LIMIT ?
    `
    args = append(args, limit)

    rows, err := db.conn.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to query filtered training sessions: %w", err)
    }
    defer rows.Close()

    var sessions []*TrainingSession
    for rows.Next() {
        var session TrainingSession
        err := rows.Scan(
            &session.ID, &session.UserID, &session.SportID, &session.SessionDate,
            &session.SessionType, &session.DurationMinutes, &session.PerceivedEffort,
            &session.Notes, &session.PerformanceNotes, &session.Skipped, &session.SkipReason,
            &session.Completed, &session.Planned, &session.CreatedAt, &session.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan training session: %w", err)
        }
        sessions = append(sessions, &session)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating training sessions: %w", err)
    }

    return sessions, nil
}

// UpdateTrainingSession updates an existing training session
func (db *DB) UpdateTrainingSession(session *TrainingSession) error {
    query := `
        UPDATE training_sessions
        SET duration_minutes = ?, perceived_effort = ?, notes = ?,
            performance_notes = ?, skipped = ?, skip_reason = ?,
            completed = ?, updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `

    _, err := db.conn.Exec(
        query,
        session.DurationMinutes, session.PerceivedEffort, session.Notes,
        session.PerformanceNotes, session.Skipped, session.SkipReason,
        session.Completed, session.ID,
    )
    if err != nil {
        return fmt.Errorf("failed to update training session: %w", err)
    }

    return nil
}

// GetPlannedSessions retrieves planned but not completed sessions for a date range
func (db *DB) GetPlannedSessions(userID int, startDate, endDate time.Time) ([]*TrainingSession, error) {
    query := `
        SELECT id, user_id, sport_id, session_date, session_type, duration_minutes,
               perceived_effort, notes, performance_notes, skipped, skip_reason, completed, planned, created_at, updated_at
        FROM training_sessions
        WHERE user_id = ? AND planned = 1 AND session_date BETWEEN ? AND ?
        ORDER BY session_date ASC
    `

    rows, err := db.conn.Query(query, userID, startDate, endDate)
    if err != nil {
        return nil, fmt.Errorf("failed to query planned sessions: %w", err)
    }
    defer rows.Close()

    var sessions []*TrainingSession
    for rows.Next() {
        var session TrainingSession
        err := rows.Scan(
            &session.ID, &session.UserID, &session.SportID, &session.SessionDate,
            &session.SessionType, &session.DurationMinutes, &session.PerceivedEffort,
            &session.Notes, &session.PerformanceNotes, &session.Skipped, &session.SkipReason,
            &session.Completed, &session.Planned,
            &session.CreatedAt, &session.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan planned session: %w", err)
        }
        sessions = append(sessions, &session)
    }

    return sessions, nil
}

// ===== Weekly Plan Queries =====

// CreateWeeklyPlan creates a new weekly plan
func (db *DB) CreateWeeklyPlan(plan *WeeklyPlan) error {
    query := `
        INSERT INTO weekly_plans
        (user_id, week_start_date, week_end_date, plan_data, rationale)
        VALUES (?, ?, ?, ?, ?)
        ON CONFLICT(user_id, week_start_date) DO UPDATE SET
            week_end_date = excluded.week_end_date,
            plan_data = excluded.plan_data,
            rationale = excluded.rationale,
            generated_at = CURRENT_TIMESTAMP
        RETURNING id, generated_at, created_at
    `

    return db.conn.QueryRow(
        query,
        plan.UserID, plan.WeekStartDate, plan.WeekEndDate,
        plan.PlanData, plan.Rationale,
    ).Scan(&plan.ID, &plan.GeneratedAt, &plan.CreatedAt)
}

// GetWeeklyPlan retrieves a weekly plan for a specific week
func (db *DB) GetWeeklyPlan(userID int, weekStart time.Time) (*WeeklyPlan, error) {
    query := `
        SELECT id, user_id, week_start_date, week_end_date, plan_data, rationale,
               generated_at, created_at
        FROM weekly_plans
        WHERE user_id = ? AND week_start_date = ?
    `

    var plan WeeklyPlan
    err := db.conn.QueryRow(query, userID, weekStart).Scan(
        &plan.ID, &plan.UserID, &plan.WeekStartDate, &plan.WeekEndDate,
        &plan.PlanData, &plan.Rationale, &plan.GeneratedAt, &plan.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("weekly plan not found")
        }
        return nil, fmt.Errorf("failed to get weekly plan: %w", err)
    }

    return &plan, nil
}

// UpdateWeeklyPlan updates an existing weekly plan
func (db *DB) UpdateWeeklyPlan(plan *WeeklyPlan) error {
    query := `
        UPDATE weekly_plans
        SET plan_data = ?, rationale = ?, generated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `

    _, err := db.conn.Exec(query, plan.PlanData, plan.Rationale, plan.ID)
    if err != nil {
        return fmt.Errorf("failed to update weekly plan: %w", err)
    }

    return nil
}

// GetRecentWeeklyPlans gets the most recent N weekly plans
func (db *DB) GetRecentWeeklyPlans(userID int, limit int) ([]*WeeklyPlan, error) {
    query := `
        SELECT id, user_id, week_start_date, week_end_date, plan_data, rationale,
               generated_at, created_at
        FROM weekly_plans
        WHERE user_id = ?
        ORDER BY week_start_date DESC
        LIMIT ?
    `

    rows, err := db.conn.Query(query, userID, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to query weekly plans: %w", err)
    }
    defer rows.Close()

    var plans []*WeeklyPlan
    for rows.Next() {
        var plan WeeklyPlan
        err := rows.Scan(
            &plan.ID, &plan.UserID, &plan.WeekStartDate, &plan.WeekEndDate,
            &plan.PlanData, &plan.Rationale, &plan.GeneratedAt, &plan.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan weekly plan: %w", err)
        }
        plans = append(plans, &plan)
    }

    return plans, nil
}

// ===== User Sport Queries =====

// CreateUserSport creates a new sport entry for a user
func (db *DB) CreateUserSport(userID int, sportName string, isPrimary bool) error {
    configPath := fmt.Sprintf("configs/%s.yaml", sportName)
    query := `
        INSERT INTO user_sports (user_id, sport_name, config_path, is_primary, current_phase)
        VALUES (?, ?, ?, ?, ?)
    `

    _, err := db.conn.Exec(query, userID, sportName, configPath, isPrimary, "foundation")
    if err != nil {
        return fmt.Errorf("failed to create user sport: %w", err)
    }

    return nil
}

// GetUserSports retrieves all sports for a user
func (db *DB) GetUserSports(userID int) ([]*UserSport, error) {
    query := `
        SELECT id, user_id, sport_name, config_path, is_primary, experience_years,
               current_phase, phase_start_date, phase_end_date, created_at, updated_at
        FROM user_sports
        WHERE user_id = ?
        ORDER BY is_primary DESC, created_at ASC
    `

    rows, err := db.conn.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query user sports: %w", err)
    }
    defer rows.Close()

    var sports []*UserSport
    for rows.Next() {
        var sport UserSport
        err := rows.Scan(
            &sport.ID, &sport.UserID, &sport.SportName, &sport.ConfigPath,
            &sport.IsPrimary, &sport.ExperienceYears, &sport.CurrentPhase,
            &sport.PhaseStartDate, &sport.PhaseEndDate,
            &sport.CreatedAt, &sport.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan user sport: %w", err)
        }
        sports = append(sports, &sport)
    }

    return sports, nil
}

// GetPrimarySport retrieves the primary sport for a user
func (db *DB) GetPrimarySport(userID int) (*UserSport, error) {
    query := `
        SELECT id, user_id, sport_name, config_path, is_primary, experience_years,
               current_phase, phase_start_date, phase_end_date, created_at, updated_at
        FROM user_sports
        WHERE user_id = ? AND is_primary = 1
        LIMIT 1
    `

    var sport UserSport
    err := db.conn.QueryRow(query, userID).Scan(
        &sport.ID, &sport.UserID, &sport.SportName, &sport.ConfigPath,
        &sport.IsPrimary, &sport.ExperienceYears, &sport.CurrentPhase,
        &sport.PhaseStartDate, &sport.PhaseEndDate,
        &sport.CreatedAt, &sport.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("no primary sport found")
        }
        return nil, fmt.Errorf("failed to get primary sport: %w", err)
    }

    return &sport, nil
}

// ===== Goal Queries =====

// CreateGoal creates a new goal
func (db *DB) CreateGoal(goal *Goal) error {
    query := `
        INSERT INTO goals (user_id, goal_type, description, target_date)
        VALUES (?, ?, ?, ?)
        RETURNING id, created_at, updated_at
    `

    return db.conn.QueryRow(
        query,
        goal.UserID, goal.GoalType, goal.Description, goal.TargetDate,
    ).Scan(&goal.ID, &goal.CreatedAt, &goal.UpdatedAt)
}

// GetGoalsByType retrieves goals by type for a user
func (db *DB) GetGoalsByType(userID int, goalType string) ([]*Goal, error) {
    query := `
        SELECT id, user_id, goal_type, description, target_date, completed,
               completed_at, created_at, updated_at
        FROM goals
        WHERE user_id = ? AND goal_type = ?
        ORDER BY completed ASC, created_at DESC
    `

    rows, err := db.conn.Query(query, userID, goalType)
    if err != nil {
        return nil, fmt.Errorf("failed to query goals: %w", err)
    }
    defer rows.Close()

    var goals []*Goal
    for rows.Next() {
        var goal Goal
        err := rows.Scan(
            &goal.ID, &goal.UserID, &goal.GoalType, &goal.Description,
            &goal.TargetDate, &goal.Completed, &goal.CompletedAt,
            &goal.CreatedAt, &goal.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan goal: %w", err)
        }
        goals = append(goals, &goal)
    }

    return goals, nil
}

// GetAllGoals retrieves all goals for a user
func (db *DB) GetAllGoals(userID int) ([]*Goal, error) {
    query := `
        SELECT id, user_id, goal_type, description, target_date, completed,
               completed_at, created_at, updated_at
        FROM goals
        WHERE user_id = ?
        ORDER BY goal_type, completed ASC, created_at DESC
    `

    rows, err := db.conn.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query goals: %w", err)
    }
    defer rows.Close()

    var goals []*Goal
    for rows.Next() {
        var goal Goal
        err := rows.Scan(
            &goal.ID, &goal.UserID, &goal.GoalType, &goal.Description,
            &goal.TargetDate, &goal.Completed, &goal.CompletedAt,
            &goal.CreatedAt, &goal.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan goal: %w", err)
        }
        goals = append(goals, &goal)
    }

    return goals, nil
}

// ===== Session Template Queries =====

// CreateSessionTemplate creates a new session template
func (db *DB) CreateSessionTemplate(template *SessionTemplate) error {
    query := `
        INSERT INTO session_templates (user_id, template_name, sport_name, session_type,
                                       duration_minutes, perceived_effort, description)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `

    result, err := db.conn.Exec(query,
        template.UserID, template.TemplateName, template.SportName,
        template.SessionType, template.DurationMinutes, template.PerceivedEffort,
        template.Description,
    )
    if err != nil {
        return fmt.Errorf("failed to create session template: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get template ID: %w", err)
    }

    template.ID = int(id)
    template.CreatedAt = time.Now()
    template.UpdatedAt = time.Now()

    return nil
}

// GetSessionTemplates retrieves all templates for a user
func (db *DB) GetSessionTemplates(userID int) ([]*SessionTemplate, error) {
    query := `
        SELECT id, user_id, template_name, sport_name, session_type,
               duration_minutes, perceived_effort, description,
               created_at, updated_at
        FROM session_templates
        WHERE user_id = ?
        ORDER BY sport_name, template_name
    `

    rows, err := db.conn.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query session templates: %w", err)
    }
    defer rows.Close()

    var templates []*SessionTemplate
    for rows.Next() {
        var tmpl SessionTemplate
        err := rows.Scan(
            &tmpl.ID, &tmpl.UserID, &tmpl.TemplateName, &tmpl.SportName,
            &tmpl.SessionType, &tmpl.DurationMinutes, &tmpl.PerceivedEffort,
            &tmpl.Description, &tmpl.CreatedAt, &tmpl.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan template: %w", err)
        }
        templates = append(templates, &tmpl)
    }

    return templates, nil
}

// GetSessionTemplatesBySport retrieves templates filtered by sport
func (db *DB) GetSessionTemplatesBySport(userID int, sportName string) ([]*SessionTemplate, error) {
    query := `
        SELECT id, user_id, template_name, sport_name, session_type,
               duration_minutes, perceived_effort, description,
               created_at, updated_at
        FROM session_templates
        WHERE user_id = ? AND sport_name = ?
        ORDER BY template_name
    `

    rows, err := db.conn.Query(query, userID, sportName)
    if err != nil {
        return nil, fmt.Errorf("failed to query session templates by sport: %w", err)
    }
    defer rows.Close()

    var templates []*SessionTemplate
    for rows.Next() {
        var tmpl SessionTemplate
        err := rows.Scan(
            &tmpl.ID, &tmpl.UserID, &tmpl.TemplateName, &tmpl.SportName,
            &tmpl.SessionType, &tmpl.DurationMinutes, &tmpl.PerceivedEffort,
            &tmpl.Description, &tmpl.CreatedAt, &tmpl.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan template: %w", err)
        }
        templates = append(templates, &tmpl)
    }

    return templates, nil
}

// GetSessionTemplate retrieves a single template by ID
func (db *DB) GetSessionTemplate(templateID int) (*SessionTemplate, error) {
    query := `
        SELECT id, user_id, template_name, sport_name, session_type,
               duration_minutes, perceived_effort, description,
               created_at, updated_at
        FROM session_templates
        WHERE id = ?
    `

    var tmpl SessionTemplate
    err := db.conn.QueryRow(query, templateID).Scan(
        &tmpl.ID, &tmpl.UserID, &tmpl.TemplateName, &tmpl.SportName,
        &tmpl.SessionType, &tmpl.DurationMinutes, &tmpl.PerceivedEffort,
        &tmpl.Description, &tmpl.CreatedAt, &tmpl.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("template not found")
        }
        return nil, fmt.Errorf("failed to get template: %w", err)
    }

    return &tmpl, nil
}

// UpdateSessionTemplate updates an existing template
func (db *DB) UpdateSessionTemplate(template *SessionTemplate) error {
    query := `
        UPDATE session_templates
        SET template_name = ?, sport_name = ?, session_type = ?,
            duration_minutes = ?, perceived_effort = ?, description = ?,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `

    _, err := db.conn.Exec(query,
        template.TemplateName, template.SportName, template.SessionType,
        template.DurationMinutes, template.PerceivedEffort, template.Description,
        template.ID,
    )
    if err != nil {
        return fmt.Errorf("failed to update template: %w", err)
    }

    return nil
}

// DeleteSessionTemplate deletes a template
func (db *DB) DeleteSessionTemplate(templateID int) error {
    query := `DELETE FROM session_templates WHERE id = ?`

    _, err := db.conn.Exec(query, templateID)
    if err != nil {
        return fmt.Errorf("failed to delete template: %w", err)
    }

    return nil
}

// GetTrainingSession retrieves a single training session by ID
func (db *DB) GetTrainingSession(sessionID int) (*TrainingSession, error) {
    query := `
        SELECT id, user_id, session_date, session_type, duration_minutes,
               perceived_effort, performance_notes, notes, completed, skipped,
               skip_reason, planned, created_at, updated_at
        FROM training_sessions
        WHERE id = ?
    `

    session := &TrainingSession{}
    err := db.conn.QueryRow(query, sessionID).Scan(
        &session.ID,
        &session.UserID,
        &session.SessionDate,
        &session.SessionType,
        &session.DurationMinutes,
        &session.PerceivedEffort,
        &session.PerformanceNotes,
        &session.Notes,
        &session.Completed,
        &session.Skipped,
        &session.SkipReason,
        &session.Planned,
        &session.CreatedAt,
        &session.UpdatedAt,
    )

    if err != nil {
        return nil, fmt.Errorf("failed to get training session: %w", err)
    }

    return session, nil
}

// ===== Rest Day Notes Queries =====

// CreateRestDayNote creates or updates a rest day note
func (db *DB) CreateRestDayNote(note *RestDayNote) error {
    query := `
        INSERT INTO rest_day_notes
        (user_id, rest_date, wellness_rating, soreness_level, motivation_level,
         recovery_activities, notes)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(user_id, rest_date) DO UPDATE SET
            wellness_rating = excluded.wellness_rating,
            soreness_level = excluded.soreness_level,
            motivation_level = excluded.motivation_level,
            recovery_activities = excluded.recovery_activities,
            notes = excluded.notes,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id, created_at, updated_at
    `

    return db.conn.QueryRow(
        query,
        note.UserID, note.RestDate, note.WellnessRating, note.SorenessLevel,
        note.MotivationLevel, note.RecoveryActivities, note.Notes,
    ).Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
}

// GetRestDayNote retrieves a rest day note for a specific date
func (db *DB) GetRestDayNote(userID int, date time.Time) (*RestDayNote, error) {
    query := `
        SELECT id, user_id, rest_date, wellness_rating, soreness_level,
               motivation_level, recovery_activities, notes, created_at, updated_at
        FROM rest_day_notes
        WHERE user_id = ? AND rest_date = ?
    `

    var note RestDayNote
    err := db.conn.QueryRow(query, userID, date).Scan(
        &note.ID, &note.UserID, &note.RestDate, &note.WellnessRating,
        &note.SorenessLevel, &note.MotivationLevel, &note.RecoveryActivities,
        &note.Notes, &note.CreatedAt, &note.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("rest day note not found")
        }
        return nil, fmt.Errorf("failed to get rest day note: %w", err)
    }

    return &note, nil
}

// GetRestDayNotes retrieves rest day notes within a date range
func (db *DB) GetRestDayNotes(userID int, startDate, endDate time.Time) ([]*RestDayNote, error) {
    query := `
        SELECT id, user_id, rest_date, wellness_rating, soreness_level,
               motivation_level, recovery_activities, notes, created_at, updated_at
        FROM rest_day_notes
        WHERE user_id = ? AND rest_date BETWEEN ? AND ?
        ORDER BY rest_date DESC
    `

    rows, err := db.conn.Query(query, userID, startDate, endDate)
    if err != nil {
        return nil, fmt.Errorf("failed to query rest day notes: %w", err)
    }
    defer rows.Close()

    var notes []*RestDayNote
    for rows.Next() {
        var note RestDayNote
        err := rows.Scan(
            &note.ID, &note.UserID, &note.RestDate, &note.WellnessRating,
            &note.SorenessLevel, &note.MotivationLevel, &note.RecoveryActivities,
            &note.Notes, &note.CreatedAt, &note.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan rest day note: %w", err)
        }
        notes = append(notes, &note)
    }

    return notes, nil
}

// GetRecentRestDayNotes retrieves the most recent N rest day notes
func (db *DB) GetRecentRestDayNotes(userID int, limit int) ([]*RestDayNote, error) {
    query := `
        SELECT id, user_id, rest_date, wellness_rating, soreness_level,
               motivation_level, recovery_activities, notes, created_at, updated_at
        FROM rest_day_notes
        WHERE user_id = ?
        ORDER BY rest_date DESC
        LIMIT ?
    `

    rows, err := db.conn.Query(query, userID, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to query recent rest day notes: %w", err)
    }
    defer rows.Close()

    var notes []*RestDayNote
    for rows.Next() {
        var note RestDayNote
        err := rows.Scan(
            &note.ID, &note.UserID, &note.RestDate, &note.WellnessRating,
            &note.SorenessLevel, &note.MotivationLevel, &note.RecoveryActivities,
            &note.Notes, &note.CreatedAt, &note.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan rest day note: %w", err)
        }
        notes = append(notes, &note)
    }

    return notes, nil
}

// DeleteRestDayNote deletes a rest day note
func (db *DB) DeleteRestDayNote(noteID int) error {
    query := `DELETE FROM rest_day_notes WHERE id = ?`

    _, err := db.conn.Exec(query, noteID)
    if err != nil {
        return fmt.Errorf("failed to delete rest day note: %w", err)
    }

    return nil
}

// GetDailyAIUsage retrieves the AI usage count for a user on a specific date
func (db *DB) GetDailyAIUsage(userID int, date time.Time) (int, error) {
    dateStr := date.Format("2006-01-02")
    query := `SELECT call_count FROM ai_usage WHERE user_id = ? AND usage_date = ?`

    var count int
    err := db.conn.QueryRow(query, userID, dateStr).Scan(&count)
    if err == sql.ErrNoRows {
        return 0, nil
    }
    if err != nil {
        return 0, fmt.Errorf("failed to get AI usage: %w", err)
    }

    return count, nil
}

// IncrementAIUsage increments the AI usage count for a user on a specific date
func (db *DB) IncrementAIUsage(userID int, date time.Time) error {
    dateStr := date.Format("2006-01-02")
    query := `
        INSERT INTO ai_usage (user_id, usage_date, call_count, created_at, updated_at)
        VALUES (?, ?, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT(user_id, usage_date) DO UPDATE SET
            call_count = call_count + 1,
            updated_at = CURRENT_TIMESTAMP
    `

    _, err := db.conn.Exec(query, userID, dateStr)
    if err != nil {
        return fmt.Errorf("failed to increment AI usage: %w", err)
    }

    return nil
}

// CheckAIRateLimit checks if the user has exceeded their daily AI call limit
// Returns (withinLimit, currentCount, error)
func (db *DB) CheckAIRateLimit(userID int, dailyLimit int) (bool, int, error) {
    today := time.Now()
    count, err := db.GetDailyAIUsage(userID, today)
    if err != nil {
        return false, 0, err
    }

    return count < dailyLimit, count, nil
}
