package db

import (
	"fmt"
	"time"
)

// WeeklyStats represents training statistics for a week
type WeeklyStats struct {
	TotalSessions     int
	CompletedSessions int
	SkippedSessions   int
	TotalMinutes      int
	AvgEffort         float64
	CompletionRate    float64
	WeekStart         time.Time
	WeekEnd           time.Time
}

// MonthlyStats represents training statistics for a month
type MonthlyStats struct {
	TotalSessions     int
	CompletedSessions int
	SkippedSessions   int
	TotalMinutes      int
	TotalHours        float64
	AvgEffort         float64
	CompletionRate    float64
	StartDate         time.Time
	EndDate           time.Time
}

// ActivityStats represents statistics for a specific activity
type ActivityStats struct {
	ActivityID   int
	ActivityName string
	Sessions     int
	TotalMinutes int
	TotalHours   float64
	AvgEffort    float64
	LastSession  *time.Time
}

// StreakStats represents training streak information
type StreakStats struct {
	CurrentStreak int
	LongestStreak int
	LastTraining  *time.Time
}

// VolumeDataPoint represents a data point for volume trends
type VolumeDataPoint struct {
	Date     time.Time
	Sessions int
	Minutes  int
}

// EffortDataPoint represents a data point for effort trends
type EffortDataPoint struct {
	Date      time.Time
	AvgEffort float64
}

// GetWeeklyStats calculates statistics for a specific week
func (db *DB) GetWeeklyStats(userID int, weekStart time.Time) (*WeeklyStats, error) {
	weekEnd := weekStart.AddDate(0, 0, 6)

	query := `
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN completed = 1 THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN skipped = 1 THEN 1 ELSE 0 END) as skipped,
			SUM(CASE WHEN completed = 1 THEN duration_minutes ELSE 0 END) as total_minutes,
			AVG(CASE WHEN completed = 1 AND perceived_effort > 0 THEN perceived_effort ELSE NULL END) as avg_effort
		FROM training_sessions
		WHERE user_id = ?
		AND session_date >= ?
		AND session_date <= ?
		AND planned = 1
	`

	var stats WeeklyStats
	var totalSessions, completed, skipped, totalMinutes int
	var avgEffort *float64

	err := db.conn.QueryRow(query, userID, weekStart, weekEnd).Scan(
		&totalSessions, &completed, &skipped, &totalMinutes, &avgEffort,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly stats: %w", err)
	}

	stats.TotalSessions = totalSessions
	stats.CompletedSessions = completed
	stats.SkippedSessions = skipped
	stats.TotalMinutes = totalMinutes
	stats.WeekStart = weekStart
	stats.WeekEnd = weekEnd

	if avgEffort != nil {
		stats.AvgEffort = *avgEffort
	}

	if totalSessions > 0 {
		stats.CompletionRate = float64(completed) / float64(totalSessions) * 100
	}

	return &stats, nil
}

// GetMonthlyStats calculates statistics for the last N days
func (db *DB) GetMonthlyStats(userID int, days int) (*MonthlyStats, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	query := `
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN completed = 1 THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN skipped = 1 THEN 1 ELSE 0 END) as skipped,
			SUM(CASE WHEN completed = 1 THEN duration_minutes ELSE 0 END) as total_minutes,
			AVG(CASE WHEN completed = 1 AND perceived_effort > 0 THEN perceived_effort ELSE NULL END) as avg_effort
		FROM training_sessions
		WHERE user_id = ?
		AND session_date >= ?
		AND session_date <= ?
		AND planned = 1
	`

	var stats MonthlyStats
	var totalSessions, completed, skipped, totalMinutes int
	var avgEffort *float64

	err := db.conn.QueryRow(query, userID, startDate, endDate).Scan(
		&totalSessions, &completed, &skipped, &totalMinutes, &avgEffort,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly stats: %w", err)
	}

	stats.TotalSessions = totalSessions
	stats.CompletedSessions = completed
	stats.SkippedSessions = skipped
	stats.TotalMinutes = totalMinutes
	stats.TotalHours = float64(totalMinutes) / 60.0
	stats.StartDate = startDate
	stats.EndDate = endDate

	if avgEffort != nil {
		stats.AvgEffort = *avgEffort
	}

	if totalSessions > 0 {
		stats.CompletionRate = float64(completed) / float64(totalSessions) * 100
	}

	return &stats, nil
}

// GetActivityStats returns statistics broken down by activity
func (db *DB) GetActivityStats(userID int, days int) ([]*ActivityStats, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	query := `
		SELECT
			COALESCE(ts.sport_id, 0) as activity_id,
			COALESCE(us.sport_name, 'Unspecified') as activity_name,
			COUNT(*) as sessions,
			SUM(ts.duration_minutes) as total_minutes,
			AVG(CASE WHEN ts.perceived_effort > 0 THEN ts.perceived_effort ELSE NULL END) as avg_effort,
			MAX(ts.session_date) as last_session
		FROM training_sessions ts
		LEFT JOIN user_sports us ON ts.sport_id = us.id
		WHERE ts.user_id = ?
		AND ts.session_date >= ?
		AND ts.session_date <= ?
		AND ts.completed = 1
		GROUP BY ts.sport_id, us.sport_name
		ORDER BY sessions DESC
	`

	rows, err := db.conn.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity stats: %w", err)
	}
	defer rows.Close()

	var stats []*ActivityStats
	for rows.Next() {
		var stat ActivityStats
		var avgEffort *float64
		var lastSession *time.Time

		err := rows.Scan(&stat.ActivityID, &stat.ActivityName, &stat.Sessions,
			&stat.TotalMinutes, &avgEffort, &lastSession)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity stats: %w", err)
		}

		stat.TotalHours = float64(stat.TotalMinutes) / 60.0
		if avgEffort != nil {
			stat.AvgEffort = *avgEffort
		}
		stat.LastSession = lastSession

		stats = append(stats, &stat)
	}

	return stats, nil
}

// GetStreakStats calculates current and longest training streaks
func (db *DB) GetStreakStats(userID int) (*StreakStats, error) {
	// Get all completed sessions ordered by date
	query := `
		SELECT DISTINCT DATE(session_date) as date
		FROM training_sessions
		WHERE user_id = ?
		AND completed = 1
		ORDER BY date DESC
	`

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get streak data: %w", err)
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			return nil, fmt.Errorf("failed to scan date: %w", err)
		}
		dates = append(dates, date)
	}

	stats := &StreakStats{
		CurrentStreak: 0,
		LongestStreak: 0,
	}

	if len(dates) == 0 {
		return stats, nil
	}

	stats.LastTraining = &dates[0]

	// Calculate current streak
	currentStreak := 1
	today := time.Now().Truncate(24 * time.Hour)
	lastDate := dates[0].Truncate(24 * time.Hour)

	// Check if last training was today or yesterday
	daysSinceLastTraining := int(today.Sub(lastDate).Hours() / 24)
	if daysSinceLastTraining > 1 {
		currentStreak = 0
	} else {
		for i := 1; i < len(dates); i++ {
			currentDate := dates[i].Truncate(24 * time.Hour)
			prevDate := dates[i-1].Truncate(24 * time.Hour)
			diff := int(prevDate.Sub(currentDate).Hours() / 24)

			if diff == 1 {
				currentStreak++
			} else {
				break
			}
		}
	}

	stats.CurrentStreak = currentStreak

	// Calculate longest streak
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(dates); i++ {
		currentDate := dates[i].Truncate(24 * time.Hour)
		prevDate := dates[i-1].Truncate(24 * time.Hour)
		diff := int(prevDate.Sub(currentDate).Hours() / 24)

		if diff == 1 {
			tempStreak++
			if tempStreak > longestStreak {
				longestStreak = tempStreak
			}
		} else {
			tempStreak = 1
		}
	}

	stats.LongestStreak = longestStreak

	return stats, nil
}

// GetVolumeDataPoints returns volume data for charting
func (db *DB) GetVolumeDataPoints(userID int, days int) ([]*VolumeDataPoint, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	query := `
		SELECT
			DATE(session_date) as date,
			COUNT(*) as sessions,
			SUM(duration_minutes) as minutes
		FROM training_sessions
		WHERE user_id = ?
		AND session_date >= ?
		AND session_date <= ?
		AND completed = 1
		GROUP BY DATE(session_date)
		ORDER BY date ASC
	`

	rows, err := db.conn.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume data: %w", err)
	}
	defer rows.Close()

	var dataPoints []*VolumeDataPoint
	for rows.Next() {
		var dp VolumeDataPoint
		err := rows.Scan(&dp.Date, &dp.Sessions, &dp.Minutes)
		if err != nil {
			return nil, fmt.Errorf("failed to scan volume data: %w", err)
		}
		dataPoints = append(dataPoints, &dp)
	}

	return dataPoints, nil
}

// GetEffortDataPoints returns effort data for charting
func (db *DB) GetEffortDataPoints(userID int, days int) ([]*EffortDataPoint, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	query := `
		SELECT
			DATE(session_date) as date,
			AVG(perceived_effort) as avg_effort
		FROM training_sessions
		WHERE user_id = ?
		AND session_date >= ?
		AND session_date <= ?
		AND completed = 1
		AND perceived_effort > 0
		GROUP BY DATE(session_date)
		ORDER BY date ASC
	`

	rows, err := db.conn.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get effort data: %w", err)
	}
	defer rows.Close()

	var dataPoints []*EffortDataPoint
	for rows.Next() {
		var dp EffortDataPoint
		err := rows.Scan(&dp.Date, &dp.AvgEffort)
		if err != nil {
			return nil, fmt.Errorf("failed to scan effort data: %w", err)
		}
		dataPoints = append(dataPoints, &dp)
	}

	return dataPoints, nil
}
