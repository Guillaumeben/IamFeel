package agent

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tuxnam/iamfeel/internal/db"
)

// PlannedSession represents a session extracted from a plan
type PlannedSession struct {
	DayOfWeek       string
	SessionType     string
	DurationMinutes int
	Notes           string
}

// ParsePlanForSessions extracts planned sessions from a weekly plan text
func ParsePlanForSessions(planText string, weekStart time.Time) []PlannedSession {
	sessions := []PlannedSession{}

	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	for _, day := range days {
		session := extractDaySession(planText, day)
		if session != nil {
			session.DayOfWeek = day
			sessions = append(sessions, *session)
		}
	}

	return sessions
}

// extractDaySession extracts session info for a specific day
func extractDaySession(planText string, day string) *PlannedSession {
	// Find the day section
	dayMarker := day + ":"
	startIdx := strings.Index(planText, dayMarker)
	if startIdx == -1 {
		return nil
	}

	// Find the end of this day's section (next day or end of plan)
	endIdx := len(planText)
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday", "NOTES:", "NEXT WEEK"}
	for _, nextDay := range days {
		if nextDay == day {
			continue
		}
		idx := strings.Index(planText[startIdx+len(dayMarker):], nextDay+":")
		if idx != -1 {
			actualIdx := startIdx + len(dayMarker) + idx
			if actualIdx < endIdx {
				endIdx = actualIdx
			}
		}
	}

	daySection := planText[startIdx:endIdx]

	// Check if it's a rest day
	lowerSection := strings.ToLower(daySection)
	if strings.Contains(lowerSection, "rest day") || strings.Contains(lowerSection, "rest/recovery") {
		return nil // Skip rest days
	}

	session := &PlannedSession{}

	// Extract session type
	sessionRegex := regexp.MustCompile(`Session:\s*(.+)`)
	if matches := sessionRegex.FindStringSubmatch(daySection); len(matches) > 1 {
		session.SessionType = strings.TrimSpace(matches[1])
	}

	// Extract duration
	durationRegex := regexp.MustCompile(`Duration:\s*(\d+)\s*min(?:utes)?`)
	if matches := durationRegex.FindStringSubmatch(daySection); len(matches) > 1 {
		duration, _ := strconv.Atoi(matches[1])
		session.DurationMinutes = duration
	}

	// Extract focus and details as notes
	focusRegex := regexp.MustCompile(`Focus:\s*(.+)`)
	var focus string
	if matches := focusRegex.FindStringSubmatch(daySection); len(matches) > 1 {
		focus = strings.TrimSpace(matches[1])
	}

	detailsRegex := regexp.MustCompile(`(?s)Details:\s*(.*?)(?:\s*Why:|$)`)
	var details string
	if matches := detailsRegex.FindStringSubmatch(daySection); len(matches) > 1 {
		details = strings.TrimSpace(matches[1])
	}

	if focus != "" {
		session.Notes = "Focus: " + focus
	}
	if details != "" {
		if session.Notes != "" {
			session.Notes += "\n\n"
		}
		session.Notes += "Details: " + details
	}

	// Only return if we found at least a session type
	if session.SessionType == "" {
		return nil
	}

	return session
}

// CreatePlannedSessions creates planned training sessions in the database
func CreatePlannedSessions(database *db.DB, userID int, weekStart time.Time, planText string) error {
	sessions := ParsePlanForSessions(planText, weekStart)

	fmt.Printf("📅 Parsed %d planned sessions from plan\n", len(sessions))

	for _, ps := range sessions {
		fmt.Printf("   • %s: %s (%d min)\n", ps.DayOfWeek, ps.SessionType, ps.DurationMinutes)

		// Calculate the date for this day
		dayOffset := getDayOffset(ps.DayOfWeek)
		sessionDate := weekStart.AddDate(0, 0, dayOffset)

		// Create planned session
		session := &db.TrainingSession{
			UserID:          userID,
			SessionDate:     sessionDate,
			SessionType:     ps.SessionType,
			DurationMinutes: ps.DurationMinutes,
			Notes:           ps.Notes,
			Completed:       false,
			Planned:         true,
		}

		if err := database.CreateTrainingSession(session); err != nil {
			fmt.Printf("   ⚠️  Failed to create session for %s: %v\n", ps.DayOfWeek, err)
			continue
		}
	}

	fmt.Printf("✓ Created %d planned sessions in database\n", len(sessions))
	return nil
}

// getDayOffset returns the number of days from Monday (0) for a given day name
func getDayOffset(dayName string) int {
	days := map[string]int{
		"Monday":    0,
		"Tuesday":   1,
		"Wednesday": 2,
		"Thursday":  3,
		"Friday":    4,
		"Saturday":  5,
		"Sunday":    6,
	}
	return days[dayName]
}
