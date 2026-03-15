package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tuxnam/iamfeel/internal/db"
)

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// ActivityRequest represents an activity create/update request
type ActivityRequest struct {
	ID                    int     `json:"id"`
	SportName             string  `json:"sport_name"`
	ExperienceYears       int     `json:"experience_years"`
	CurrentPhase          string  `json:"current_phase"`
	PhaseStartDate        string  `json:"phase_start_date"`
	PhaseEndDate          string  `json:"phase_end_date"`
	GoalType              string  `json:"goal_type"`
	Priority              string  `json:"priority"`
	TargetSessionsPerWeek float64 `json:"target_sessions_per_week"`
	Notes                 string  `json:"notes"`
	IsPrimary             bool    `json:"is_primary"`
}

// ActivityResponse represents an activity in API responses
type ActivityResponse struct {
	ID                    int     `json:"id"`
	SportName             string  `json:"sport_name"`
	ExperienceYears       int     `json:"experience_years"`
	CurrentPhase          string  `json:"current_phase"`
	PhaseStartDate        string  `json:"phase_start_date"`
	PhaseEndDate          string  `json:"phase_end_date"`
	GoalType              string  `json:"goal_type"`
	Priority              string  `json:"priority"`
	TargetSessionsPerWeek float64 `json:"target_sessions_per_week"`
	Notes                 string  `json:"notes"`
	IsPrimary             bool    `json:"is_primary"`
}

// HandleGetActivities returns all activities for the current user
func (s *Server) HandleGetActivities(w http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to load user", http.StatusInternalServerError)
		return
	}

	activities, err := s.db.GetUserSports(user.ID)
	if err != nil {
		http.Error(w, "Failed to get activities", http.StatusInternalServerError)
		log.Printf("Error getting activities: %v", err)
		return
	}

	// Convert to response format
	response := make([]ActivityResponse, 0, len(activities))
	for _, activity := range activities {
		resp := ActivityResponse{
			ID:                    activity.ID,
			SportName:             activity.SportName,
			ExperienceYears:       activity.ExperienceYears,
			CurrentPhase:          activity.CurrentPhase,
			GoalType:              activity.GoalType,
			Priority:              activity.Priority,
			TargetSessionsPerWeek: activity.TargetSessionsPerWeek,
			Notes:                 activity.Notes,
			IsPrimary:             activity.IsPrimary,
		}
		if activity.PhaseStartDate != nil {
			resp.PhaseStartDate = activity.PhaseStartDate.Format("2006-01-02")
		}
		if activity.PhaseEndDate != nil {
			resp.PhaseEndDate = activity.PhaseEndDate.Format("2006-01-02")
		}
		response = append(response, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleAddActivity adds a new activity for the user
func (s *Server) HandleAddActivity(w http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to load user", http.StatusInternalServerError)
		return
	}

	var req ActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.SportName == "" {
		http.Error(w, "Sport name is required", http.StatusBadRequest)
		return
	}
	if req.CurrentPhase == "" {
		req.CurrentPhase = "foundation"
	}

	// Set defaults if not provided
	if req.GoalType == "" {
		if req.IsPrimary {
			req.GoalType = "competition_prep"
		} else {
			req.GoalType = "maintenance"
		}
	}
	if req.Priority == "" {
		if req.IsPrimary {
			req.Priority = "high"
		} else {
			req.Priority = "medium"
		}
	}

	// Create UserSport
	sport := &db.UserSport{
		UserID:                user.ID,
		SportName:             req.SportName,
		ConfigPath:            fmt.Sprintf("configs/%s.yaml", req.SportName),
		IsPrimary:             req.IsPrimary,
		ExperienceYears:       req.ExperienceYears,
		CurrentPhase:          req.CurrentPhase,
		GoalType:              req.GoalType,
		Priority:              req.Priority,
		TargetSessionsPerWeek: req.TargetSessionsPerWeek,
		Notes:                 req.Notes,
	}

	// Parse dates if provided
	if req.PhaseStartDate != "" {
		if date, err := parseDate(req.PhaseStartDate); err == nil {
			sport.PhaseStartDate = &date
		}
	}
	if req.PhaseEndDate != "" {
		if date, err := parseDate(req.PhaseEndDate); err == nil {
			sport.PhaseEndDate = &date
		}
	}

	// Insert into database
	err = s.db.UpdateUserSport(sport)
	if err != nil {
		// If update fails, try creating
		err = s.db.CreateUserSport(user.ID, req.SportName, req.IsPrimary)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create activity: %v", err), http.StatusInternalServerError)
			log.Printf("Error creating activity: %v", err)
			return
		}

		// Get the created sport to return its ID
		activities, err := s.db.GetUserSports(user.ID)
		if err == nil {
			for _, a := range activities {
				if a.SportName == req.SportName {
					sport.ID = a.ID
					break
				}
			}
		}

		// Now update with full details
		err = s.db.UpdateUserSport(sport)
		if err != nil {
			log.Printf("Warning: Failed to update activity details: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      sport.ID,
	})
}

// HandleUpdateActivity updates an existing activity
func (s *Server) HandleUpdateActivity(w http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to load user", http.StatusInternalServerError)
		return
	}

	// Get activity ID from URL
	activityIDStr := r.URL.Query().Get("id")
	if activityIDStr == "" {
		http.Error(w, "Activity ID is required", http.StatusBadRequest)
		return
	}
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	var req ActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create UserSport with updated data
	sport := &db.UserSport{
		ID:                    activityID,
		UserID:                user.ID,
		SportName:             req.SportName,
		ConfigPath:            fmt.Sprintf("configs/%s.yaml", req.SportName),
		IsPrimary:             req.IsPrimary,
		ExperienceYears:       req.ExperienceYears,
		CurrentPhase:          req.CurrentPhase,
		GoalType:              req.GoalType,
		Priority:              req.Priority,
		TargetSessionsPerWeek: req.TargetSessionsPerWeek,
		Notes:                 req.Notes,
	}

	// Parse dates if provided
	if req.PhaseStartDate != "" {
		if date, err := parseDate(req.PhaseStartDate); err == nil {
			sport.PhaseStartDate = &date
		}
	}
	if req.PhaseEndDate != "" {
		if date, err := parseDate(req.PhaseEndDate); err == nil {
			sport.PhaseEndDate = &date
		}
	}

	// Update in database
	err = s.db.UpdateUserSport(sport)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update activity: %v", err), http.StatusInternalServerError)
		log.Printf("Error updating activity: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// HandleDeleteActivity deletes an activity
func (s *Server) HandleDeleteActivity(w http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Failed to load user", http.StatusInternalServerError)
		return
	}

	// Get activity ID from URL
	activityIDStr := r.URL.Query().Get("id")
	if activityIDStr == "" {
		http.Error(w, "Activity ID is required", http.StatusBadRequest)
		return
	}
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	// Verify the activity belongs to the user
	activities, err := s.db.GetUserSports(user.ID)
	if err != nil {
		http.Error(w, "Failed to get activities", http.StatusInternalServerError)
		return
	}

	found := false
	for _, activity := range activities {
		if activity.ID == activityID {
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	// Delete the activity
	_, err = s.db.Conn().Exec("DELETE FROM user_sports WHERE id = ? AND user_id = ?", activityID, user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete activity: %v", err), http.StatusInternalServerError)
		log.Printf("Error deleting activity: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
