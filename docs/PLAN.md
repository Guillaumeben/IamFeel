# IamFeel - Development Plan

## Project Vision
AI-powered training assistant that helps plan sport sessions based on schedule, training history, current phase (strength, technique, etc.), and personal goals. Sport-agnostic design allows adaptation for boxing, ironman, golf, or any sport.

## Core Principles
- Sport-agnostic architecture (configurable per sport)
- Agent-native design (file-based state, transparent)
- Simple, local-first UI
- Privacy-focused (runs locally, configurable LLM)

---

## Phase 1: Foundation
**Goal:** Basic project structure and data layer

### Project Setup
- [x] Create directory structure
- [x] Initialize Go module
- [x] Create planning documents
- [x] Add `.gitignore` (ignore `data/`, `*.db`, `user_config.yaml`)
- [x] Set up initial dependencies in `go.mod`

### Database Layer
- [x] Design SQLite schema
  - [x] Users table (profile, current phase, goals)
  - [x] Sports config table
  - [x] Training sessions table
  - [x] Nutrition log table
  - [x] Weekly plans table
- [x] Create migration system
- [x] Implement database models (`internal/db/models.go`)
- [x] Implement basic queries (`internal/db/queries.go`)
- [ ] Write database tests

### Configuration System
- [x] Create sport config structure (`configs/sports/template.yaml`)
- [x] Create boxing config (`configs/sports/boxing.yaml`)
- [x] Implement config loader (`internal/config/loader.go`)
- [x] Create user profile structure (`user_config.yaml` template)
- [x] Write config validation

---

## Phase 2: Onboarding
**Goal:** Interactive CLI to build user profile

### Onboarding Wizard
- [x] Create CLI framework (`cmd/cli/main.go`)
- [x] Implement onboarding command
  - [x] Sport selection (with option for multiple sports)
  - [x] Equipment/gym membership questions
  - [x] Club sessions configuration
  - [x] Fitness level & experience
  - [x] Goal setting (short/medium/long-term)
  - [x] Nutrition preferences
  - [x] Weekly availability
- [x] Save user config to `data/user_config.yaml`
- [x] Seed database with initial data
- [ ] Add ability to update profile later

---

## Phase 3: Agent Core
**Goal:** Claude integration for plan generation

### Agent Implementation
- [x] Set up Anthropic API integration (HTTP client)
- [x] Create configurable model settings
- [x] Design system prompt template (`internal/agent/prompts.go`)
- [x] Implement sport-specific context loading
- [x] Create plan generation function
  - [x] Load user profile
  - [x] Load training history (last 14 days)
  - [x] Load current phase
  - [x] Query Claude for weekly plan
- [x] Parse and structure agent responses
- [x] Store generated plans in database

### Prompt Engineering
- [x] Create base system prompt
- [x] Create sport-specific prompt sections
- [x] Design plan output format (text-based)
- [ ] Add examples for few-shot learning (future enhancement)
- [ ] Test prompt with various scenarios (ongoing)

---

## Phase 4: Plan Generation
**Goal:** Generate weekly training plans via CLI

### CLI Plan Command
- [x] Implement `plan` command (`cmd/cli/plan.go`)
- [x] Add weekly plan generation
- [ ] Add optional monthly plan generation (future)
- [x] Display generated plan in terminal
- [x] Save plan to database
- [ ] Add flag for plan period (weekly/monthly) (future)

### Plan Logic
- [x] Implement history analyzer (last 14 days)
- [x] Factor in current training phase
- [x] Consider club sessions (paid vs membership)
- [x] Balance different training types (handled by AI)
- [x] Include rest/recovery days (handled by AI)
- [x] Generate rationale for each session (handled by AI)

---

## Phase 5: Web Interface
**Goal:** Simple dashboard to view plans

### Web Server
- [x] Set up HTTP server (`cmd/server/main.go`)
- [x] Add router (chi)
- [x] Serve static files
- [x] Create base HTML template

### Dashboard
- [x] Create dashboard handler (`internal/api/handlers.go`)
- [x] Design dashboard HTML (`web/templates/dashboard.html`)
- [x] Display current week's plan
- [x] Show today's workout highlighted
- [x] Add basic CSS styling (`web/static/style.css`)
- [x] Make responsive (mobile-friendly)

### Session Logging
- [x] Create session log form
- [x] Implement POST handler for logging sessions
- [x] Save to database
- [x] Add perceived effort rating
- [x] Add notes field
- [x] Show completion status on dashboard

---

## Phase 6: Chat Interface
**Goal:** Interactive adjustments and questions

### Chat Implementation
- [ ] Create chat endpoint
- [ ] Design chat UI component
- [ ] Implement streaming responses
- [ ] Load conversation context (current plan, recent history)
- [ ] Handle plan modification requests
- [ ] Update database when plans change
- [ ] Add chat history persistence

### Chat Features
- [ ] "I'm too tired today" → adjust workout
- [ ] "What should I focus on this week?"
- [ ] "Why did you suggest this session?"
- [ ] Reschedule/swap workouts
- [ ] Add unplanned sessions

---

## Phase 7: History & Tracking
**Goal:** View past sessions and progress

### History View
- [ ] Create history page
- [ ] Display last 6 months of sessions
- [ ] Filter by session type
- [ ] Search functionality
- [ ] Export to CSV/JSON

### Nutrition Tracking
- [ ] Create nutrition log form
- [ ] Design nutrition dashboard section
- [ ] Daily/weekly nutrition view
- [ ] Supplement tracking
- [ ] Simple macros tracking (optional)

### Progress Visualization
- [ ] Session frequency chart
- [ ] Training volume over time
- [ ] Phase progression timeline
- [ ] Goal progress indicators

---

## Phase 8: Polish & Documentation
**Goal:** Make it ready for others to use

### Documentation
- [ ] Write comprehensive README
  - [ ] Project overview
  - [ ] Installation instructions
  - [ ] Quick start guide
  - [ ] Configuration guide
- [ ] Create sport config guide (how to add new sports)
- [ ] Document API endpoints
- [ ] Add inline code documentation
- [ ] Create troubleshooting guide

### Developer Experience
- [ ] Add example configs
- [ ] Create demo user data
- [ ] Add Makefile for common tasks
- [ ] Set up linting (golangci-lint)
- [ ] Add pre-commit hooks

### Testing
- [ ] Unit tests for critical paths
- [ ] Integration tests for database
- [ ] Test onboarding flow
- [ ] Test plan generation
- [ ] Test chat modifications

---

## Success Metrics
- [ ] Can onboard new user in < 5 minutes
- [ ] Generates relevant weekly plan based on profile
- [ ] Dashboard loads in < 1 second
- [ ] Chat responses feel natural and helpful
- [ ] Another person can adapt for different sport in < 30 minutes

---

## Notes
- Keep it simple - no over-engineering
- Focus on making it work for boxing first
- Ensure sport-agnostic design as we go
- Iterate based on actual usage
- Document decisions in BACKLOG.md
