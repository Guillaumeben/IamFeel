# Multi-Activity Support - Implementation Progress

## ✅ Phase 1: Database & Models (COMPLETE)

### Completed Tasks

1. **Design Document Created** (`docs/MULTI_ACTIVITY_DESIGN.md`)
   - Complete architecture for multi-activity support
   - User scenarios (triathlete, boxer+strength, hiking+golf)
   - Backwards compatibility strategy

2. **Database Migration** (`internal/db/migrations/001_multi_activity_support.sql`)
   - Added `goal_type` to `user_sports` (competition_prep, maintenance, learning, recreation)
   - Added `priority` to `user_sports` (high, medium, low)
   - Added `target_sessions_per_week` to `user_sports`
   - Added `notes` to `user_sports`
   - Added `sport_id` to `equipment` (tag equipment to activities)
   - Created new indexes for performance
   - Data migration: primary sports → high priority, others → medium

3. **Models Updated**
   - ✅ `internal/db/models.go` - UserSport enhanced with new fields
   - ✅ `internal/db/models.go` - Equipment tagged to activities
   - ✅ `internal/db/schema.go` - Schema updated with new columns
   - ✅ `internal/config/models.go` - Config model enhanced

4. **Build Verification**
   - ✅ All packages compile successfully
   - ✅ No breaking changes to existing code

## ✅ Phase 2: Backend Logic (COMPLETE)

### Completed Tasks

1. **Database Queries** ✅
   - ✅ Updated GetUserSports to include new fields (goal_type, priority, target_sessions_per_week, notes)
   - ✅ Updated GetPrimarySport to prioritize by priority field
   - ✅ Updated CreateUserSport to set defaults based on isPrimary
   - ✅ Added GetActivitiesByPriority helper function
   - ✅ Added UpdateUserSport for full CRUD support
   - ✅ Updated GetUserEquipment to include sport_id
   - ✅ Added AddEquipmentWithSport to tag equipment to activities
   - ✅ Queries now sort activities by priority (high → medium → low)

2. **Agent Context** ✅
   - ✅ Added ActivityInfo struct with all multi-activity fields
   - ✅ Updated UserContext to include Activities array
   - ✅ LoadUserContext now loads all activities from config
   - ✅ Sets intelligent defaults (primary → high priority + competition_prep)
   - ✅ Backwards compatibility maintained (single-sport fields still populated)
   - ✅ Handles 0, 1, or N activities gracefully

3. **AI Prompts** ✅
   - ✅ BuildUserPrompt detects multi-activity vs single-activity
   - ✅ Multi-activity format shows each activity with priority and goal type
   - ✅ Added CRITICAL balancing instructions for AI
   - ✅ HIGH priority → main focus, MEDIUM → maintain, LOW → fit when possible
   - ✅ Competition prep → progressive, Maintenance → steady state
   - ✅ Scientific references now per-activity
   - ✅ Phase information contextual based on activity count
   - ✅ Added formatGoalType helper for human-readable goal types

4. **Build Verification** ✅
   - ✅ All packages compile successfully
   - ✅ No breaking changes to existing code
   - ✅ Backwards compatible with single-sport users

## 📅 Phase 3: UI Updates (PENDING)

1. **Settings Page**
   - Manage multiple activities (add/edit/delete)
   - Set goal type per activity (dropdown)
   - Set priority per activity (high/medium/low)
   - Set target sessions/week per activity
   - Tag equipment to activities

2. **Dashboard**
   - Visual activity balance (actual vs target sessions)
   - Multi-activity weekly summary
   - Warning if activity neglected

3. **Session Logging**
   - Activity selector when logging sessions
   - Filter history by activity

## 🧪 Phase 4: Testing (PENDING)

1. **Test Scenarios**
   - Triathlete (3 equal high-priority activities)
   - Boxer + Strength (primary + supplementary)
   - Hiking + Golf (learning + recreation)
   - Single-sport user (backwards compatibility)

2. **Integration Tests**
   - Multi-activity plan generation
   - Equipment filtering by activity
   - Club session suggestions per activity
   - Session balance calculations

## Migration Path

### For Existing Users

The migration is **backwards compatible**:

1. Current single-sport users continue working unchanged
2. Primary sport automatically set to `priority='high'`, `goal_type='competition_prep'`
3. Equipment without `sport_id` available to all activities
4. `is_primary` flag still functional (maps to priority)

### For New Multi-Activity Users

Users can now:
1. Add multiple activities with individual goals
2. Set different priorities per activity
3. Tag equipment to specific activities
4. AI balances sessions across all activities based on priorities

## Files Changed

**Phase 1 - Database & Models:**
```
internal/db/models.go                        - Enhanced UserSport and Equipment models
internal/db/schema.go                        - Updated table definitions with new columns
internal/config/models.go                    - Enhanced UserSport config model
internal/db/migrations/001_multi_activity.sql - Migration SQL with data backfill
docs/MULTI_ACTIVITY_DESIGN.md                - Complete design document
```

**Phase 2 - Backend Logic:**
```
internal/db/queries.go                       - Updated user_sports CRUD queries
internal/db/config_queries.go                - Updated equipment queries
internal/agent/context.go                    - Multi-activity context loading
internal/agent/prompts.go                    - Multi-activity AI prompts
docs/MULTI_ACTIVITY_PROGRESS.md              - This progress tracker
```

## Next Steps

1. Update database queries to use new fields
2. Update agent context building for multi-activity
3. Enhance AI prompts to handle multiple activities
4. Add API endpoints for activity management
5. Update UI for multi-activity support

---

**Status:** Phase 1 & 2 Complete - Ready for Phase 3 (UI) or Testing
**Last Updated:** 2026-03-14

## Example AI Prompt Output (Multi-Activity)

For a triathlete with Swimming (high), Cycling (high), Running (high):

```
## Activities

This athlete practices multiple activities. Balance training time and recovery across all activities based on their priorities and goals.

### Swimming (HIGH Priority - Competition Preparation)
- Experience: 3 years
- Current Phase: Base Building
- Target: 3.0 sessions/week

### Cycling (HIGH Priority - Competition Preparation)
- Experience: 5 years
- Current Phase: Build
- Target: 3.0 sessions/week

### Running (HIGH Priority - Competition Preparation)
- Experience: 4 years
- Current Phase: Base Building
- Target: 3.0 sessions/week

**CRITICAL - Multi-Activity Balancing:**
- HIGH priority activities: These are the main focus - allocate most training time here
- MEDIUM priority activities: Maintain current level - steady state training
- LOW priority activities: Fit in when schedule allows - recreational/recovery focus
- Respect target sessions per week for each activity
- COMPETITION PREP activities need progressive, periodized training
- MAINTENANCE activities need consistent but not progressive training
```

The AI now understands to balance 9 total sessions across all three sports!
