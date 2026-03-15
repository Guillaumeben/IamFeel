-- Migration: Multi-Activity Support
-- Add goal_type, priority, and target sessions to user_sports table
-- Tag equipment to specific activities

-- 1. Enhance user_sports table with activity-specific fields
ALTER TABLE user_sports ADD COLUMN goal_type TEXT
  CHECK(goal_type IN ('competition_prep', 'maintenance', 'learning', 'recreation'))
  DEFAULT 'maintenance';

ALTER TABLE user_sports ADD COLUMN priority TEXT
  CHECK(priority IN ('high', 'medium', 'low'))
  DEFAULT 'medium';

ALTER TABLE user_sports ADD COLUMN target_sessions_per_week REAL DEFAULT 0;

ALTER TABLE user_sports ADD COLUMN notes TEXT;

-- 2. Tag equipment to activities
ALTER TABLE equipment ADD COLUMN sport_id INTEGER
  REFERENCES user_sports(id) ON DELETE SET NULL;

-- 3. Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_equipment_sport ON equipment(sport_id);
CREATE INDEX IF NOT EXISTS idx_user_sports_priority ON user_sports(user_id, priority);
CREATE INDEX IF NOT EXISTS idx_club_sessions_sport_active ON club_sessions(sport_id, active);

-- 4. Data migration: Set sensible defaults for existing data
-- Set primary sport to high priority competition prep
UPDATE user_sports
SET goal_type = 'competition_prep',
    priority = 'high'
WHERE is_primary = 1;

-- Set non-primary sports to medium priority maintenance
UPDATE user_sports
SET goal_type = 'maintenance',
    priority = 'medium'
WHERE is_primary = 0 OR is_primary IS NULL;

-- 5. Migration complete
-- Note: is_primary column remains for backwards compatibility
