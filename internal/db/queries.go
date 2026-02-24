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
        SELECT id, name, age, experience_level, created_at, updated_at
        FROM users
        WHERE id = ?
    `

    var user User
    err := db.conn.QueryRow(query, userID).Scan(
        &user.ID, &user.Name, &user.Age, &user.ExperienceLevel,
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
        SELECT id, name, age, experience_level, created_at, updated_at
        FROM users
        ORDER BY id ASC
        LIMIT 1
    `

    var user User
    err := db.conn.QueryRow(query).Scan(
        &user.ID, &user.Name, &user.Age, &user.ExperienceLevel,
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
        SET name = ?, age = ?, experience_level = ?, updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `

    _, err := db.conn.Exec(query, user.Name, user.Age, user.ExperienceLevel, user.ID)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }

    return nil
}

// GetAllUsers retrieves all users from the database
func (db *DB) GetAllUsers() ([]*User, error) {
    query := `
        SELECT id, name, age, experience_level, created_at, updated_at
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
            &user.ID, &user.Name, &user.Age, &user.ExperienceLevel,
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
