# IamFeel - Backlog

Future features, ideas, and enhancements. Items here are not prioritized - just captured for consideration.

---

## Calendar Integration
- [ ] Google Calendar sync (read/write)
- [ ] Automatically block training time
- [ ] Detect schedule conflicts
- [ ] Sync club session attendance
- [ ] Send workout reminders
- [ ] iCal export for other calendar apps

## Advanced Analytics
- [ ] Training load calculation (acute/chronic workload ratio)
- [ ] Fatigue monitoring
- [ ] Recovery predictions
- [ ] Performance trends over time
- [ ] Identify plateaus or overtraining
- [ ] Injury risk indicators
- [ ] Correlate nutrition with performance

## Nutrition Enhancements
- [ ] Meal planning suggestions
- [ ] Recipe database
- [ ] Macro calculator based on training phase
- [ ] Supplement timing recommendations
- [ ] Hydration tracking
- [ ] Pre/post workout nutrition guidance
- [ ] Integration with MyFitnessPal or similar

## Social & Sharing
- [ ] Share workouts with training partners
- [ ] Coach mode (trainer can manage multiple athletes)
- [ ] Export training log for coaches
- [ ] Community sport configs (marketplace)
- [ ] Achievement badges/milestones

## Mobile Experience
- [ ] Progressive Web App (PWA)
- [ ] Native mobile app
- [ ] Workout timer with audio cues
- [ ] Quick session logging
- [ ] Offline mode
- [ ] Apple Watch/Wear OS integration

## Wearables & IoT
- [ ] Import data from Garmin/Fitbit/Whoop
- [ ] Heart rate zone training
- [ ] Sleep quality integration
- [ ] HRV monitoring for recovery
- [ ] Auto-log workouts from fitness tracker

## AI Enhancements
- [ ] Voice interface ("Hey Coach, what's today's workout?")
- [ ] Computer vision for form analysis (upload videos)
- [ ] Automatic periodization adjustments
- [ ] Predict optimal deload weeks
- [ ] Learn from user preferences over time
- [ ] Multi-agent system (specialized agents per domain)

## Knowledge Base & Learning
- [ ] Article ingestion system
  - [ ] Feed training articles to enhance AI expertise
  - [ ] Parse and store article content (PDF, web articles)
  - [ ] Use articles to influence scheduling decisions
  - [ ] Cite sources in plan recommendations
  - [ ] Tag articles by topic (periodization, nutrition, technique, etc.)
  - [ ] Search knowledge base for relevant info
- [ ] Video content integration (future)
  - [ ] Transcript extraction from YouTube videos
  - [ ] Add coaching videos to knowledge base
  - [ ] Use video insights in recommendations
- [ ] Expert knowledge curation
  - [ ] Pre-loaded knowledge packs per sport (boxing coaches, S&C experts)
  - [ ] User can add custom articles/papers
  - [ ] Update knowledge base over time
- [ ] Context-aware recommendations
  - [ ] Reference scientific principles in plans
  - [ ] Explain rationale using learned expertise

## UI/UX Enhancements
- [ ] Dynamic theming based on primary sport
  - [ ] Boxing: red/black theme with glove iconography
  - [ ] Running: blue/green with trail imagery
  - [ ] Strength: iron/steel theme
  - [ ] Customizable color schemes
  - [ ] Sport-specific background imagery
- [ ] Dark mode / Light mode toggle
- [ ] Customizable dashboard layout
- [ ] Quick action shortcuts
- [ ] Drag-and-drop workout reordering
- [ ] Calendar view with color-coded sessions

## Sport-Specific Features

### Boxing
- [ ] Sparring partner matching
- [ ] Competition prep mode
- [ ] Weight class management
- [ ] Fight camp templates
- [ ] Hand/wrist injury tracking
- [ ] Shadow boxing drill library

### Endurance Sports
- [ ] Race calendar integration
- [ ] Taper period automation
- [ ] Course-specific training
- [ ] Brick workout planning (triathlon)
- [ ] Altitude training adjustments

### Strength Sports
- [ ] 1RM calculator
- [ ] Powerlifting meet prep
- [ ] Deload week automation
- [ ] Exercise library with videos
- [ ] Plate calculator (what to load on bar)

## Developer Experience
- [ ] MCP server implementation (use with Claude Desktop)
- [ ] API documentation (OpenAPI/Swagger)
- [ ] Docker containerization
- [ ] One-click deploy (Railway, Fly.io)
- [ ] Plugin system for custom integrations
- [ ] Webhook support for external tools

## Data & Privacy
- [ ] End-to-end encryption option
- [ ] Self-hosted mode (no API calls)
- [ ] Local LLM support (Ollama, LM Studio)
- [ ] Data export (full backup)
- [ ] GDPR compliance features
- [ ] Anonymous usage analytics (opt-in)

## Gamification
- [ ] Streak tracking
- [ ] Training consistency score
- [ ] Challenge system (30-day challenges)
- [ ] Visual progress graphs
- [ ] Goal completion celebrations
- [ ] Training journal with photos

## Coaching Tools
- [ ] Session template library
- [ ] Drill database
- [ ] Exercise database (local DB with alternatives/progressions)
  - [ ] Reduce LLM calls for exercise substitutions
  - [ ] Exercise variations (beginner/intermediate/advanced)
  - [ ] Equipment-based alternatives (barbell → dumbbell → bodyweight)
  - [ ] Muscle group categorization
  - [ ] Exercise progression pathways
- [ ] Video embedding for exercise demos
  - [ ] YouTube video integration for technique demos
  - [ ] Embed videos directly in workout plans
  - [ ] Searchable video library by exercise/drill
  - [ ] Curated playlists per sport
- [ ] Warmup/cooldown generators
- [ ] Mobility work suggestions
- [ ] Prehab/rehab exercise library

## Smart Scheduling
- [ ] Weather-aware planning (outdoor sessions)
- [ ] Travel mode (hotel gym, bodyweight)
- [ ] Injury mode (work around limitations)
- [ ] Sick day adjustments
- [ ] Time zone handling for travel
- [ ] Automatic plan adjustments based on adherence

## Multi-Sport Support
- [ ] Cross-training coordination (e.g., boxing + running)
- [ ] Sport priority weighting
- [ ] Conflict resolution (don't schedule heavy legs before sparring)
- [ ] Off-season vs in-season modes
- [ ] Sport switching (triathlon → marathon)

## Reporting
- [ ] Weekly summary emails
- [ ] Monthly progress reports
- [ ] Custom report builder
- [ ] PDF export of training logs
- [ ] Share reports with healthcare providers

## Community Configs
- [ ] Sport config repository (GitHub-based)
- [ ] Rating/review system for configs
- [ ] Fork and customize existing configs
- [ ] Version control for personal configs
- [ ] Import/export configs

## Performance Optimization
- [ ] Caching layer for LLM responses
  - [ ] Cache exercise alternatives/substitutions
  - [ ] Cache common queries (e.g., "what's a good warmup?")
  - [ ] Cache plan templates
  - [ ] TTL-based cache invalidation
  - [ ] Redis integration (optional, for advanced users)
- [ ] Batch plan generation (generate month at once)
- [ ] Background job processing
- [ ] Database optimization for large histories
  - [ ] Indexed queries for common lookups
  - [ ] Archive old data (>12 months)
  - [ ] Query optimization for dashboard
- [ ] CDN for static assets
- [ ] Lazy loading for UI components
- [ ] Preload common data on dashboard

## Accessibility
- [ ] Screen reader support
- [ ] Keyboard navigation
- [ ] High contrast mode
- [ ] Text size adjustments
- [ ] Multi-language support

## Testing & Quality
- [ ] E2E testing with Playwright
- [ ] Load testing
- [ ] Security audit
- [ ] Performance benchmarks
- [ ] Fuzzing for edge cases

---

## Ideas to Explore

### Hybrid Training Models
Instead of pure LLM reasoning, combine:
- Rule-based periodization (proven principles)
- LLM for personalization and adaptation
- Machine learning for pattern recognition in user data

### Physiological Modeling
- Simple fatigue accumulation model
- Recovery time predictions based on session type
- Training readiness score

### Coach Personas
Let users pick coaching style:
- "Aggressive" - push hard, high volume
- "Conservative" - prioritize recovery
- "Balanced" - middle ground
- "Adaptive" - learn from user feedback

### Session Library
Pre-built workout templates:
- Quick 30-min sessions
- Full 90-min sessions
- Active recovery
- Technique-focused
- Conditioning-focused

### Motivation System
- Motivational quotes/messages
- Progress photos timeline
- "Why" reminders (your goals)
- Accountability partner notifications
- Celebration of milestones

---

## Technical Debt / Refactoring

Track technical improvements here:
- [ ] TODO: Add proper error handling patterns
- [ ] TODO: Implement structured logging
- [ ] TODO: Add configuration validation
- [ ] TODO: Create mock LLM for testing
- [ ] TODO: Optimize database indexes

---

## Research Topics

Things to investigate before implementing:
- [ ] Best practices for training periodization
- [ ] LLM prompt engineering for coaching
- [ ] SQLite performance with large datasets
- [ ] Go best practices for CLI tools
- [ ] Effective workout plan formats

---

## Won't Do (For Now)

Features explicitly out of scope:
- ❌ Meal delivery service integration
- ❌ E-commerce for supplements
- ❌ Live video coaching
- ❌ Complex team management
- ❌ Blockchain/crypto anything

---

## Decision Log

Track important decisions:

**2026-02-20: Go vs Python**
- Decision: Use Go
- Reason: Better performance, single binary, user prefers Go

**2026-02-20: SQLite vs PostgreSQL**
- Decision: SQLite
- Reason: Local-first, simpler setup, sufficient for use case

**2026-02-20: Weekly vs Daily Planning**
- Decision: Weekly default, configurable to monthly
- Reason: Balance between flexibility and planning overhead

**2026-02-20: Calendar Integration**
- Decision: Defer to post-MVP
- Reason: Focus on core functionality first

**2026-02-20: Additional Backlog Items**
- Added: Exercise database for offline alternatives
- Added: Knowledge base system for articles/videos
- Added: Sport-specific UI theming
- Added: Enhanced caching strategy
- Reason: Reduce LLM costs, enhance expertise, improve UX

---

## User Feedback

Capture user feedback here as it comes in:
- (No feedback yet - project just started!)
