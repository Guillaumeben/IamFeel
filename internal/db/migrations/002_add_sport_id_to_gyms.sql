-- Migration: Add sport_id to gyms table
-- This allows linking gyms/clubs to specific activities

-- Add sport_id column to gyms table
ALTER TABLE gyms ADD COLUMN sport_id INTEGER REFERENCES user_sports(id) ON DELETE SET NULL;

-- Create index for gym-sport lookups
CREATE INDEX IF NOT EXISTS idx_gyms_sport ON gyms(sport_id);
