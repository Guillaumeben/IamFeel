package db

import (
	"database/sql"
	"fmt"
)

// ==================== User Functions ====================

// GetUserName retrieves the name of a user by ID
func (db *DB) GetUserName(userID int) (string, error) {
	var name string
	err := db.conn.QueryRow("SELECT name FROM users WHERE id = ?", userID).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("failed to get user name: %w", err)
	}
	return name, nil
}

// ==================== Goal Functions ====================

// GetUserGoals retrieves all goals for a user
func (db *DB) GetUserGoals(userID int) ([]*Goal, error) {
	query := `SELECT id, user_id, goal_type, description, target_date, completed, completed_at, created_at, updated_at
              FROM goals WHERE user_id = ? ORDER BY goal_type, description`

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query goals: %w", err)
	}
	defer rows.Close()

	var goals []*Goal
	for rows.Next() {
		var goal Goal
		err := rows.Scan(&goal.ID, &goal.UserID, &goal.GoalType, &goal.Description,
			&goal.TargetDate, &goal.Completed, &goal.CompletedAt, &goal.CreatedAt, &goal.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}
		goals = append(goals, &goal)
	}

	return goals, nil
}

// DeleteUserGoals deletes all goals for a user (used before re-inserting updated goals)
func (db *DB) DeleteUserGoals(userID int) error {
	_, err := db.conn.Exec("DELETE FROM goals WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user goals: %w", err)
	}
	return nil
}

// CreateGoalSimple creates a new goal with minimal fields (for settings save)
func (db *DB) CreateGoalSimple(userID int, goalType, description string) error {
	query := `INSERT INTO goals (user_id, goal_type, description, completed, created_at, updated_at)
	          VALUES (?, ?, ?, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err := db.conn.Exec(query, userID, goalType, description)
	if err != nil {
		return fmt.Errorf("failed to create goal: %w", err)
	}
	return nil
}

// ==================== Sport Functions ====================

// UpdatePrimarySport sets a sport as the primary sport for a user
func (db *DB) UpdatePrimarySport(userID int, sportName string, experienceYears int) error {
	// First, unset all sports as non-primary
	_, err := db.conn.Exec("UPDATE user_sports SET is_primary = 0 WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to unset primary sports: %w", err)
	}

	// Check if sport exists for user
	var count int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM user_sports WHERE user_id = ? AND sport_name = ?", userID, sportName).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check sport existence: %w", err)
	}

	if count == 0 {
		// Sport doesn't exist, create it
		_, err = db.conn.Exec("INSERT INTO user_sports (user_id, sport_name, is_primary, experience_years, created_at, updated_at) VALUES (?, ?, 1, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", userID, sportName, experienceYears)
		if err != nil {
			return fmt.Errorf("failed to create primary sport: %w", err)
		}
	} else {
		// Sport exists, set it as primary and update experience years
		_, err = db.conn.Exec("UPDATE user_sports SET is_primary = 1, experience_years = ?, updated_at = CURRENT_TIMESTAMP WHERE user_id = ? AND sport_name = ?", experienceYears, userID, sportName)
		if err != nil {
			return fmt.Errorf("failed to update primary sport: %w", err)
		}
	}

	return nil
}

// ==================== Equipment Functions ====================

// GetUserEquipment retrieves all equipment for a user
func (db *DB) GetUserEquipment(userID int) ([]*Equipment, error) {
	query := `SELECT id, user_id, location, equipment_name, sport_id, created_at
              FROM equipment WHERE user_id = ? ORDER BY location, equipment_name`

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query equipment: %w", err)
	}
	defer rows.Close()

	var equipment []*Equipment
	for rows.Next() {
		var e Equipment
		if err := rows.Scan(&e.ID, &e.UserID, &e.Location, &e.EquipmentName, &e.SportID, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan equipment: %w", err)
		}
		equipment = append(equipment, &e)
	}

	return equipment, nil
}

// AddEquipment adds a piece of equipment for a user
func (db *DB) AddEquipment(userID int, location, equipmentName string) error {
	query := `INSERT INTO equipment (user_id, location, equipment_name, sport_id, created_at)
              VALUES (?, ?, ?, NULL, CURRENT_TIMESTAMP)`

	_, err := db.conn.Exec(query, userID, location, equipmentName)
	if err != nil {
		return fmt.Errorf("failed to add equipment: %w", err)
	}

	return nil
}

// AddEquipmentWithSport adds equipment tagged to a specific activity
func (db *DB) AddEquipmentWithSport(userID int, location, equipmentName string, sportID *int) error {
	query := `INSERT INTO equipment (user_id, location, equipment_name, sport_id, created_at)
              VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := db.conn.Exec(query, userID, location, equipmentName, sportID)
	if err != nil {
		return fmt.Errorf("failed to add equipment with sport: %w", err)
	}

	return nil
}

// DeleteEquipment removes a piece of equipment
func (db *DB) DeleteEquipment(equipmentID int) error {
	query := `DELETE FROM equipment WHERE id = ?`

	_, err := db.conn.Exec(query, equipmentID)
	if err != nil {
		return fmt.Errorf("failed to delete equipment: %w", err)
	}

	return nil
}

// ClearUserEquipment removes all equipment for a user (used during sync)
func (db *DB) ClearUserEquipment(userID int) error {
	query := `DELETE FROM equipment WHERE user_id = ?`

	_, err := db.conn.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to clear equipment: %w", err)
	}

	return nil
}

// ==================== Gym Functions ====================

// GetUserGyms retrieves all gyms for a user
func (db *DB) GetUserGyms(userID int) ([]*Gym, error) {
	query := `SELECT id, user_id, name, type, membership, available_days, sport_id, sessions_limit, limit_period, created_at, updated_at
              FROM gyms WHERE user_id = ? ORDER BY name`

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query gyms: %w", err)
	}
	defer rows.Close()

	var gyms []*Gym
	for rows.Next() {
		var g Gym
		if err := rows.Scan(&g.ID, &g.UserID, &g.Name, &g.Type, &g.Membership,
			&g.AvailableDays, &g.SportID, &g.SessionsLimit, &g.LimitPeriod, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan gym: %w", err)
		}
		gyms = append(gyms, &g)
	}

	return gyms, nil
}

// CreateGym adds a new gym membership
func (db *DB) CreateGym(gym *Gym) (int64, error) {
	query := `INSERT INTO gyms (user_id, name, type, membership, available_days, sport_id, sessions_limit, limit_period, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := db.conn.Exec(query, gym.UserID, gym.Name, gym.Type, gym.Membership, gym.AvailableDays, gym.SportID, gym.SessionsLimit, gym.LimitPeriod)
	if err != nil {
		return 0, fmt.Errorf("failed to create gym: %w", err)
	}

	return result.LastInsertId()
}

// UpdateGym updates gym information
func (db *DB) UpdateGym(gym *Gym) error {
	query := `UPDATE gyms SET name = ?, type = ?, membership = ?, available_days = ?, sport_id = ?, sessions_limit = ?, limit_period = ?,
              updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	_, err := db.conn.Exec(query, gym.Name, gym.Type, gym.Membership, gym.AvailableDays, gym.SportID, gym.SessionsLimit, gym.LimitPeriod, gym.ID)
	if err != nil {
		return fmt.Errorf("failed to update gym: %w", err)
	}

	return nil
}

// DeleteGym removes a gym (and cascades to club sessions)
func (db *DB) DeleteGym(gymID int) error {
	query := `DELETE FROM gyms WHERE id = ?`

	_, err := db.conn.Exec(query, gymID)
	if err != nil {
		return fmt.Errorf("failed to delete gym: %w", err)
	}

	return nil
}

// ClearUserGyms deletes all gyms for a user
func (db *DB) ClearUserGyms(userID int) error {
	_, err := db.conn.Exec("DELETE FROM gyms WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to clear user gyms: %w", err)
	}
	return nil
}

// ==================== Club Session Functions ====================

// GetGymSessions retrieves all club sessions for a gym
func (db *DB) GetGymSessions(gymID int) ([]*ClubSession, error) {
	query := `SELECT id, user_id, gym_id, sport_id, session_name, description, occurrences,
              cost, day_of_week, time, duration_minutes, session_type,
              notes, active, created_at, updated_at
              FROM club_sessions WHERE gym_id = ? AND active = 1 ORDER BY day_of_week, time`

	rows, err := db.conn.Query(query, gymID)
	if err != nil {
		return nil, fmt.Errorf("failed to query club sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*ClubSession
	for rows.Next() {
		var s ClubSession
		if err := rows.Scan(&s.ID, &s.UserID, &s.GymID, &s.SportID, &s.SessionName, &s.Description,
			&s.Occurrences, &s.Cost, &s.DayOfWeek, &s.Time, &s.DurationMinutes, &s.SessionType,
			&s.Notes, &s.Active, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan club session: %w", err)
		}
		sessions = append(sessions, &s)
	}

	return sessions, nil
}

// CreateClubSession adds a new club session
func (db *DB) CreateClubSession(session *ClubSession) (int64, error) {
	query := `INSERT INTO club_sessions (user_id, gym_id, sport_id, session_name, description,
              occurrences, cost, day_of_week, time, duration_minutes, session_type,
              notes, active, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := db.conn.Exec(query, session.UserID, session.GymID, session.SportID,
		session.SessionName, session.Description, session.Occurrences, session.Cost,
		session.DayOfWeek, session.Time, session.DurationMinutes, session.SessionType,
		session.Notes, session.Active)
	if err != nil {
		return 0, fmt.Errorf("failed to create club session: %w", err)
	}

	return result.LastInsertId()
}

// DeleteClubSession removes a club session
func (db *DB) DeleteClubSession(sessionID int) error {
	query := `DELETE FROM club_sessions WHERE id = ?`

	_, err := db.conn.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete club session: %w", err)
	}

	return nil
}

// ==================== Availability Functions ====================

// GetUserAvailability retrieves weekly availability for a user
func (db *DB) GetUserAvailability(userID int) ([]*Availability, error) {
	query := `SELECT id, user_id, day_of_week, morning, lunch, evening, preferred_location,
              notes, created_at, updated_at
              FROM availability WHERE user_id = ?
              ORDER BY CASE day_of_week
                  WHEN 'Monday' THEN 1
                  WHEN 'Tuesday' THEN 2
                  WHEN 'Wednesday' THEN 3
                  WHEN 'Thursday' THEN 4
                  WHEN 'Friday' THEN 5
                  WHEN 'Saturday' THEN 6
                  WHEN 'Sunday' THEN 7
              END`

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query availability: %w", err)
	}
	defer rows.Close()

	var availability []*Availability
	for rows.Next() {
		var a Availability
		var preferredLocation sql.NullString
		if err := rows.Scan(&a.ID, &a.UserID, &a.DayOfWeek, &a.Morning, &a.Lunch, &a.Evening,
			&preferredLocation, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan availability: %w", err)
		}
		if preferredLocation.Valid {
			a.PreferredLocation = preferredLocation.String
		}
		availability = append(availability, &a)
	}

	return availability, nil
}

// UpsertAvailability creates or updates availability for a specific day
func (db *DB) UpsertAvailability(avail *Availability) error {
	query := `INSERT INTO availability (user_id, day_of_week, morning, lunch, evening,
              notes, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
              ON CONFLICT(user_id, day_of_week) DO UPDATE SET
                  morning = excluded.morning,
                  lunch = excluded.lunch,
                  evening = excluded.evening,
                  notes = excluded.notes,
                  updated_at = CURRENT_TIMESTAMP`

	_, err := db.conn.Exec(query, avail.UserID, avail.DayOfWeek, avail.Morning, avail.Lunch,
		avail.Evening, avail.Notes)
	if err != nil {
		return fmt.Errorf("failed to upsert availability: %w", err)
	}

	return nil
}

// ==================== Supplement Functions ====================

// GetUserSupplements retrieves all active supplements for a user
func (db *DB) GetUserSupplements(userID int) ([]*SupplementDefinition, error) {
	query := `SELECT id, user_id, name, dosage, timing, active, created_at, updated_at
              FROM supplement_definitions WHERE user_id = ? AND active = 1 ORDER BY name`

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query supplements: %w", err)
	}
	defer rows.Close()

	var supplements []*SupplementDefinition
	for rows.Next() {
		var s SupplementDefinition
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.Dosage, &s.Timing, &s.Active,
			&s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan supplement: %w", err)
		}
		supplements = append(supplements, &s)
	}

	return supplements, nil
}

// UpsertSupplement creates or updates a supplement
func (db *DB) UpsertSupplement(supp *SupplementDefinition) error {
	query := `INSERT INTO supplement_definitions (user_id, name, dosage, timing, active, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
              ON CONFLICT(user_id, name) DO UPDATE SET
                  dosage = excluded.dosage,
                  timing = excluded.timing,
                  active = excluded.active,
                  updated_at = CURRENT_TIMESTAMP`

	_, err := db.conn.Exec(query, supp.UserID, supp.Name, supp.Dosage, supp.Timing, supp.Active)
	if err != nil {
		return fmt.Errorf("failed to upsert supplement: %w", err)
	}

	return nil
}

// DeleteSupplement soft-deletes a supplement (sets active = false)
func (db *DB) DeleteSupplement(supplementID int) error {
	query := `UPDATE supplement_definitions SET active = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	_, err := db.conn.Exec(query, supplementID)
	if err != nil {
		return fmt.Errorf("failed to delete supplement: %w", err)
	}

	return nil
}

// ==================== User Preferences Functions ====================

// GetUserPreferences retrieves preferences for a user
func (db *DB) GetUserPreferences(userID int) (*UserPreferences, error) {
	query := `SELECT id, user_id, primary_goal, sessions_per_week, preferred_duration,
              preferred_session_times, session_duration_preference, intensity_preference,
              recovery_priority, plan_frequency, allow_short_sessions, max_sessions_per_day,
              notes, created_at, updated_at
              FROM user_preferences WHERE user_id = ?`

	var prefs UserPreferences
	err := db.conn.QueryRow(query, userID).Scan(
		&prefs.ID, &prefs.UserID, &prefs.PrimaryGoal, &prefs.SessionsPerWeek,
		&prefs.PreferredDuration, &prefs.PreferredSessionTimes, &prefs.SessionDurationPreference,
		&prefs.IntensityPreference, &prefs.RecoveryPriority, &prefs.PlanFrequency,
		&prefs.AllowShortSessions, &prefs.MaxSessionsPerDay,
		&prefs.Notes, &prefs.CreatedAt, &prefs.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No preferences yet
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	return &prefs, nil
}

// UpsertUserPreferences creates or updates user preferences
func (db *DB) UpsertUserPreferences(prefs *UserPreferences) error {
	query := `INSERT INTO user_preferences (user_id, primary_goal, sessions_per_week, preferred_duration,
              preferred_session_times, session_duration_preference, intensity_preference, recovery_priority,
              plan_frequency, allow_short_sessions, max_sessions_per_day, notes, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
              ON CONFLICT(user_id) DO UPDATE SET
                  primary_goal = excluded.primary_goal,
                  sessions_per_week = excluded.sessions_per_week,
                  preferred_duration = excluded.preferred_duration,
                  preferred_session_times = excluded.preferred_session_times,
                  session_duration_preference = excluded.session_duration_preference,
                  intensity_preference = excluded.intensity_preference,
                  recovery_priority = excluded.recovery_priority,
                  plan_frequency = excluded.plan_frequency,
                  allow_short_sessions = excluded.allow_short_sessions,
                  max_sessions_per_day = excluded.max_sessions_per_day,
                  notes = excluded.notes,
                  updated_at = CURRENT_TIMESTAMP`

	_, err := db.conn.Exec(query, prefs.UserID, prefs.PrimaryGoal, prefs.SessionsPerWeek,
		prefs.PreferredDuration, prefs.PreferredSessionTimes, prefs.SessionDurationPreference,
		prefs.IntensityPreference, prefs.RecoveryPriority, prefs.PlanFrequency,
		prefs.AllowShortSessions, prefs.MaxSessionsPerDay, prefs.Notes)
	if err != nil {
		return fmt.Errorf("failed to upsert user preferences: %w", err)
	}

	return nil
}

// ==================== Coach Settings Functions ====================

// GetCoachSettings retrieves coach settings for a user
func (db *DB) GetCoachSettings(userID int) (*CoachSettings, error) {
	query := `SELECT id, user_id, model, temperature, coaching_style, explanation_detail,
              created_at, updated_at
              FROM coach_settings WHERE user_id = ?`

	var settings CoachSettings
	err := db.conn.QueryRow(query, userID).Scan(
		&settings.ID, &settings.UserID, &settings.Model, &settings.Temperature,
		&settings.CoachingStyle, &settings.ExplanationDetail, &settings.CreatedAt, &settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No settings yet, will use defaults
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get coach settings: %w", err)
	}

	return &settings, nil
}

// UpsertCoachSettings creates or updates coach settings
func (db *DB) UpsertCoachSettings(settings *CoachSettings) error {
	query := `INSERT INTO coach_settings (user_id, model, temperature, coaching_style,
              explanation_detail, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
              ON CONFLICT(user_id) DO UPDATE SET
                  model = excluded.model,
                  temperature = excluded.temperature,
                  coaching_style = excluded.coaching_style,
                  explanation_detail = excluded.explanation_detail,
                  updated_at = CURRENT_TIMESTAMP`

	_, err := db.conn.Exec(query, settings.UserID, settings.Model, settings.Temperature,
		settings.CoachingStyle, settings.ExplanationDetail)
	if err != nil {
		return fmt.Errorf("failed to upsert coach settings: %w", err)
	}

	return nil
}

// ==================== Tracking Settings Functions ====================

// GetTrackingSettings retrieves tracking settings for a user
func (db *DB) GetTrackingSettings(userID int) (*TrackingSettings, error) {
	query := `SELECT id, user_id, history_months, track_supplements, track_sleep, track_weight,
              created_at, updated_at
              FROM tracking_settings WHERE user_id = ?`

	var settings TrackingSettings
	err := db.conn.QueryRow(query, userID).Scan(
		&settings.ID, &settings.UserID, &settings.HistoryMonths, &settings.TrackSupplements,
		&settings.TrackSleep, &settings.TrackWeight, &settings.CreatedAt, &settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No settings yet, will use defaults
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tracking settings: %w", err)
	}

	return &settings, nil
}

// UpsertTrackingSettings creates or updates tracking settings
func (db *DB) UpsertTrackingSettings(settings *TrackingSettings) error {
	query := `INSERT INTO tracking_settings (user_id, history_months, track_supplements, track_sleep,
              track_weight, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
              ON CONFLICT(user_id) DO UPDATE SET
                  history_months = excluded.history_months,
                  track_supplements = excluded.track_supplements,
                  track_sleep = excluded.track_sleep,
                  track_weight = excluded.track_weight,
                  updated_at = CURRENT_TIMESTAMP`

	_, err := db.conn.Exec(query, settings.UserID, settings.HistoryMonths, settings.TrackSupplements,
		settings.TrackSleep, settings.TrackWeight)
	if err != nil {
		return fmt.Errorf("failed to upsert tracking settings: %w", err)
	}

	return nil
}

// ==================== Helper Functions ====================

// GetUserConfigFromDB retrieves all config data from database for a user
func (db *DB) GetUserConfigFromDB(userID int) (*UserConfigData, error) {
	config := &UserConfigData{UserID: userID}

	// Get equipment
	equipment, err := db.GetUserEquipment(userID)
	if err != nil {
		return nil, err
	}
	config.Equipment = equipment

	// Get gyms
	gyms, err := db.GetUserGyms(userID)
	if err != nil {
		return nil, err
	}
	config.Gyms = gyms

	// Get availability
	availability, err := db.GetUserAvailability(userID)
	if err != nil {
		return nil, err
	}
	config.Availability = availability

	// Get goals
	goals, err := db.GetAllGoals(userID)
	if err != nil {
		return nil, err
	}
	config.Goals = goals

	// Get supplements
	supplements, err := db.GetUserSupplements(userID)
	if err != nil {
		return nil, err
	}
	config.Supplements = supplements

	// Get preferences
	prefs, err := db.GetUserPreferences(userID)
	if err != nil {
		return nil, err
	}
	config.Preferences = prefs

	// Get coach settings
	coach, err := db.GetCoachSettings(userID)
	if err != nil {
		return nil, err
	}
	config.CoachSettings = coach

	// Get tracking settings
	tracking, err := db.GetTrackingSettings(userID)
	if err != nil {
		return nil, err
	}
	config.TrackingSettings = tracking

	return config, nil
}

// UserConfigData holds all configuration data for a user
type UserConfigData struct {
	UserID           int
	Equipment        []*Equipment
	Gyms             []*Gym
	Availability     []*Availability
	Goals            []*Goal
	Supplements      []*SupplementDefinition
	Preferences      *UserPreferences
	CoachSettings    *CoachSettings
	TrackingSettings *TrackingSettings
}
