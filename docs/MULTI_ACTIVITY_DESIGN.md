# Multi-Activity Support - Design Document

## Overview
Transform IamFeel from single-sport to multi-activity support, enabling users to train for multiple sports/activities with different goals and priorities.

## Current State Analysis

### Existing Tables (Good Foundation)
- ✅ `user_sports` - Already tracks multiple sports per user
- ✅ `training_sessions.sport_id` - Sessions can be tagged to sports
- ✅ `club_sessions.sport_id` - Club sessions can be tagged to sports

### Gaps to Fill
- ❌ No goal type per activity (competition prep vs maintenance)
- ❌ No priority levels per activity
- ❌ No target sessions/week per activity
- ❌ Equipment not tagged to specific activities
- ❌ Single "primary sport" model in config/prompts
- ❌ AI prompts assume one sport focus

## Database Schema Changes

### 1. Enhance `user_sports` Table
Add columns to support activity-specific goals and priorities:

```sql
ALTER TABLE user_sports ADD COLUMN goal_type TEXT
  CHECK(goal_type IN ('competition_prep', 'maintenance', 'learning', 'recreation'))
  DEFAULT 'maintenance';

ALTER TABLE user_sports ADD COLUMN priority TEXT
  CHECK(priority IN ('high', 'medium', 'low'))
  DEFAULT 'medium';

ALTER TABLE user_sports ADD COLUMN target_sessions_per_week REAL DEFAULT 0;

ALTER TABLE user_sports ADD COLUMN notes TEXT;
```

**Migration Strategy:**
- Set existing primary sport to `goal_type='competition_prep'`, `priority='high'`
- Set other sports to `goal_type='maintenance'`, `priority='medium'`

### 2. Tag Equipment to Activities
Add sport/activity reference to equipment:

```sql
ALTER TABLE equipment ADD COLUMN sport_id INTEGER
  REFERENCES user_sports(id) ON DELETE SET NULL;

CREATE INDEX idx_equipment_sport ON equipment(sport_id);
```

**Migration Strategy:**
- Leave existing equipment with `sport_id=NULL` (available for all activities)
- New equipment can be tagged to specific activities

### 3. Update Indexes
```sql
CREATE INDEX IF NOT EXISTS idx_club_sessions_sport ON club_sessions(sport_id);
CREATE INDEX IF NOT EXISTS idx_user_sports_priority ON user_sports(user_id, priority);
```

## Data Model Updates

### `db/models.go` - UserSport Enhancement
```go
type UserSport struct {
    ID                    int       `json:"id"`
    UserID                int       `json:"user_id"`
    SportName             string    `json:"sport_name"`
    ConfigPath            string    `json:"config_path"`
    IsPrimary             bool      `json:"is_primary"` // Deprecated, use Priority instead
    ExperienceYears       int       `json:"experience_years"`
    CurrentPhase          string    `json:"current_phase"`
    PhaseStartDate        *time.Time `json:"phase_start_date"`
    PhaseEndDate          *time.Time `json:"phase_end_date"`

    // New fields for multi-activity support
    GoalType              string    `json:"goal_type"` // competition_prep, maintenance, learning, recreation
    Priority              string    `json:"priority"`  // high, medium, low
    TargetSessionsPerWeek float64   `json:"target_sessions_per_week"`
    Notes                 string    `json:"notes"`

    CreatedAt             time.Time `json:"created_at"`
    UpdatedAt             time.Time `json:"updated_at"`
}
```

### `db/models.go` - Equipment Enhancement
```go
type Equipment struct {
    ID            int       `json:"id"`
    UserID        int       `json:"user_id"`
    Location      string    `json:"location"`
    EquipmentName string    `json:"equipment_name"`
    SportID       *int      `json:"sport_id"` // NEW: link to specific activity (NULL = available for all)
    CreatedAt     time.Time `json:"created_at"`
}
```

### `config/models.go` - UserSport Enhancement
```go
type UserSport struct {
    Name                  string `yaml:"name"`
    ConfigFile            string `yaml:"config_file"`
    Primary               bool   `yaml:"primary"` // Deprecated, use Priority
    ExperienceYears       int    `yaml:"experience_years"`
    CurrentPhase          string `yaml:"current_phase"`
    PhaseStartDate        string `yaml:"phase_start_date"`
    PhaseEndDate          string `yaml:"phase_end_date"`

    // New fields
    GoalType              string  `yaml:"goal_type,omitempty"` // competition_prep, maintenance, learning, recreation
    Priority              string  `yaml:"priority,omitempty"`  // high, medium, low
    TargetSessionsPerWeek float64 `yaml:"target_sessions_per_week,omitempty"`
    Notes                 string  `yaml:"notes,omitempty"`
}
```

## AI Prompt Changes

### Context Building Updates
The AI prompts need to handle multiple activities with different training modes:

**Current (Single Sport):**
```
Experience in Boxing: 2 years
Current Phase: Base Building
Primary Goal: Improve conditioning
```

**New (Multi-Activity):**
```
## Activities

### Boxing (HIGH PRIORITY - Competition Prep)
- Experience: 2 years
- Current Phase: Base Building
- Goal: Prepare for amateur bout in June
- Target: 4 sessions/week

### Running (MEDIUM PRIORITY - Maintenance)
- Experience: 5 years
- Training Mode: Maintain fitness, no specific event
- Target: 2 sessions/week

### Yoga (LOW PRIORITY - Recreation)
- Experience: 1 year
- Training Mode: Recovery and flexibility support
- Target: 1 session/week
```

### Prompt Template Changes

1. **System Prompt**: Update to handle multiple activities with different focus levels
2. **User Context**: Build activity-specific sections
3. **Balancing Logic**: AI should:
   - Prioritize high-priority activities
   - Maintain maintenance-mode activities at steady state
   - Fit recreation activities when schedule allows
   - Balance total weekly volume across all activities

## Implementation Phases

### Phase 1: Database & Models (Current)
- [x] Design schema changes
- [ ] Write migration SQL
- [ ] Update db/models.go
- [ ] Update db/schema.go
- [ ] Update db queries for new fields
- [ ] Update config/models.go

### Phase 2: Backend Logic
- [ ] Update agent/context.go to build multi-activity context
- [ ] Update agent/prompts.go for multi-activity support
- [ ] Update API handlers for activity CRUD
- [ ] Add activity balance calculations

### Phase 3: UI Updates
- [ ] Settings page: Manage multiple activities
- [ ] Activity priority/goal type selectors
- [ ] Equipment tagging to activities
- [ ] Dashboard: Activity balance visualization
- [ ] Session logging: Activity selector

### Phase 4: Testing
- [ ] Test triathlon scenario (3 equal activities)
- [ ] Test primary+supplementary (boxing + strength)
- [ ] Test recreational multi-sport (hiking + golf)
- [ ] Verify backwards compatibility (single sport users)

## Backwards Compatibility

**Critical**: Must not break existing single-sport users

**Strategy:**
1. Make new columns nullable with defaults
2. If user has only one sport, AI behaves as before
3. `is_primary` flag remains functional (maps to high priority)
4. Equipment with `sport_id=NULL` available to all activities
5. Existing sessions with `sport_id=NULL` still work

## Example User Scenarios

### Triathlete (Equal Focus)
```yaml
sports:
  - name: Swimming
    goal_type: competition_prep
    priority: high
    target_sessions_per_week: 3
  - name: Cycling
    goal_type: competition_prep
    priority: high
    target_sessions_per_week: 3
  - name: Running
    goal_type: competition_prep
    priority: high
    target_sessions_per_week: 3
```

### Boxer + Strength (Primary + Supplementary)
```yaml
sports:
  - name: Boxing
    goal_type: competition_prep
    priority: high
    target_sessions_per_week: 4
    current_phase: Peak
  - name: Strength Training
    goal_type: maintenance
    priority: medium
    target_sessions_per_week: 2
```

### Recreational (Hiking + Golf)
```yaml
sports:
  - name: Hiking
    goal_type: learning
    priority: high
    target_sessions_per_week: 2
    notes: "Preparing for multi-day trek in 4 months"
  - name: Golf
    goal_type: recreation
    priority: low
    target_sessions_per_week: 1
```

## Success Metrics

- ✅ Users can add N activities (not just 1 primary)
- ✅ Each activity has independent goals and priorities
- ✅ AI balances sessions across activities appropriately
- ✅ Equipment can be tagged to specific activities
- ✅ Dashboard shows activity balance
- ✅ Existing single-sport users unaffected

---

**Status:** Design Complete - Ready for Implementation
**Next Step:** Create migration SQL and update models
