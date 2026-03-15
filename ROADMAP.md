# IamFeel Roadmap

## Quick Wins (Current Focus)

### 1. Search/Filter on History Page
- [ ] Filter by date range (start/end date pickers)
- [ ] Filter by session type
- [ ] Filter by effort level range
- [ ] Clear filters button

### 2. Copy Previous Week's Plan
- [x] "Copy Last Week" button on settings page
- [x] Option to request AI adjustments when copying
- [ ] Preview before saving

### 3. Rest Day Notes
- [x] Add rest day logging functionality
- [x] Track recovery activities (stretching, foam rolling, etc.)
- [x] Simple wellness check-in (how you feel)

### 4. Equipment Quick-Add
- [ ] Pre-populated templates by sport
  - Boxing: Heavy bag, speed bag, jump rope, wraps, gloves
  - Fitness: Barbell, dumbbells, resistance bands, bench, squat rack
  - Running: Running shoes, GPS watch, heart rate monitor
  - BJJ: Gi, rashguard, mouthguard, knee pads
  - Cycling: Road bike, helmet, cycling shoes, bike computer
  - Swimming: Goggles, swim cap, pull buoy, kickboard
  - CrossFit: Barbell, plates, kettlebells, pull-up bar, rowing machine
- [ ] "Add Common Equipment" button on settings

### 5. Session Templates
- [ ] Save completed sessions as templates
- [ ] Quick-log from template
- [ ] Edit templates
- [ ] Template library by sport

## UI/UX Improvements (High Priority)

### Settings Page Restructuring
**Priority: URGENT** - Current settings page is 75KB and overwhelming
- [ ] Break settings into tabbed sections:
  - Profile tab (name, age, experience)
  - Activities tab (sports, goals, phases)
  - Schedule tab (availability, preferences)
  - Equipment tab (home equipment, gyms/clubs)
  - Coach tab (AI settings, coaching style)
- [ ] Add progressive disclosure (collapsible sections)
- [ ] Visual feedback on unsaved changes
- [ ] Inline validation
- [ ] Better help text and tooltips
- [ ] Save confirmation toasts

### Progressive Web App (PWA)
- [x] Add manifest.json for installability
- [x] Service worker for offline support
- [x] Add to home screen functionality
- [ ] Offline data sync strategy
- [x] App-like navigation (no browser chrome)
- [ ] Push notification support

### Empty States & Onboarding
- [ ] Engaging empty states with actionable guidance
- [ ] First-time user onboarding flow
  - Welcome screen
  - Quick profile setup wizard
  - Sample plan generation
  - Feature highlights tour
- [ ] Contextual help throughout app
- [ ] Progress indicators during onboarding

### Progress Visualization
- [ ] Weekly completion progress bars
- [ ] Goal progress trackers with visual indicators
- [ ] Trend arrows (↑↓) with percentage changes
- [ ] Comparison views (this week vs last week)
- [ ] Heat maps for training frequency
- [ ] Calendar view showing training days
- [ ] Milestone celebrations

### Mobile Optimization
- [ ] Responsive breakpoints review
- [ ] Touch-friendly tap targets (44px minimum)
- [ ] Hamburger navigation for mobile
- [ ] Full-screen modals on mobile
- [ ] Swipe gestures (tabs, history cards)
- [ ] Mobile form optimization
- [ ] Bottom navigation bar option

### Notification System
- [x] Toast notifications instead of page reloads
- [ ] Loading states with skeleton screens
- [x] Success/error feedback
- [ ] Workout reminders
- [ ] Supplement timing notifications
- [ ] Streak notifications
- [ ] Non-intrusive notification center

### Accessibility (a11y)
- [ ] ARIA labels on interactive elements
- [ ] Keyboard navigation support
- [ ] Focus management in modals
- [ ] Color contrast audit (WCAG AA)
- [ ] Reduced motion preferences
- [ ] Screen reader testing
- [ ] Skip to main content link

### Achievement System & Gamification
- [ ] Streak badges (7-day, 30-day, 100-day)
- [ ] Milestone achievements (first session, 50 sessions, etc.)
- [ ] Personal records celebration
- [ ] Goal completion badges
- [ ] Achievement showcase on dashboard
- [ ] Share achievements feature
- [ ] Level/XP system based on consistency

## High Impact Features

### Multi-Activity Support (Core Enhancement)
**Goal:** Support users with multiple sports/activities (triathletes, recreational multi-sport, etc.)

**MUST-HAVE (Foundational):**
- [ ] Multiple activities per user (not just "primary sport")
  - [ ] Database schema: activities table with user_id, activity_name, priority, goal_type
  - [ ] Per-activity goal types: Competition prep, Maintenance, Learning/beginner, Recreation
  - [ ] Per-activity priority: High/Medium/Low for time balancing
- [ ] Activity-specific session types
  - [ ] Define session type taxonomies per activity (e.g., Running: intervals/tempo/long run/recovery)
  - [ ] Update AI prompts to understand activity-appropriate session types
  - [ ] Session logging includes activity_id reference
- [ ] Activity-specific context
  - [ ] Tag equipment to specific activities
  - [ ] Tag gym/club sessions to activities
  - [ ] Venue/facility availability per activity

**SHOULD-HAVE (Significantly Better UX):**
- [ ] Event/milestone calendar
  - [ ] Support multiple event types: races, trips, tournaments, competitions
  - [ ] Event date + activity + importance level
  - [ ] AI auto-adjusts training phases (taper, peak) based on event proximity
- [ ] Flexible periodization per activity
  - [ ] Per-activity training mode: Event prep (periodized), Maintenance, Learning
  - [ ] Independent phases per activity (Activity A in Build, Activity B in Maintenance)
- [ ] Activity balance dashboard
  - [ ] Weekly time/volume distribution across activities
  - [ ] Target vs actual per activity
  - [ ] Visual warnings for neglected activities

**NICE-TO-HAVE (Polish):**
- [ ] Cross-training intelligence
  - [ ] AI suggests complementary activities (e.g., "Yoga helps golf flexibility + hiking recovery")
  - [ ] Recognize activities that support multiple goals
- [ ] Multi-activity sessions
  - [ ] Support brick workouts (triathlon: bike-to-run transitions)
  - [ ] Combined sessions (hike + trail run, bike commute to gym)
  - [ ] Single session with multiple activity components
- [ ] Seasonal/contextual awareness
  - [ ] Activity seasonality (ski season, golf season)
  - [ ] Weather constraints per activity type
  - [ ] Automatic activity prioritization based on season

### Analytics Dashboard
- [ ] Progress charts (sessions completed vs planned)
- [ ] Effort trends visualization
- [ ] Streak tracking
- [ ] Weekly/monthly summary stats
- [ ] Completion rate metrics

### AI Coach Recommendations (Active)
- [ ] Analyze skip patterns for injury prevention
- [ ] Recovery suggestions based on effort levels
- [ ] Goal progress assessment
- [ ] Adaptive plan suggestions
- [ ] Detect overtraining patterns

### Performance Metrics Tracking
- [ ] PR (Personal Record) tracking per exercise
- [ ] Volume calculations (total weight lifted per week)
- [ ] Exercise library with personal bests
- [ ] Auto-parse performance notes for PRs
- [ ] Progress photos integration

### Recovery & Wellness Tracking
- [ ] Sleep logging (hours, quality)
- [ ] Soreness/fatigue scale (1-10)
- [ ] Recovery score calculation
- [ ] Auto-suggest deload weeks
- [ ] Hydration tracking

### Nutrition Planning
- [ ] Meal planning tied to training days
- [ ] Simple macro tracking
- [ ] Pre/post workout nutrition suggestions
- [ ] Nutrition goals by training phase
- [ ] Integration with existing supplement tracking

### Mobile Experience Improvements
- [ ] Responsive design enhancements
- [ ] PWA capabilities (offline access)
- [ ] Installable app
- [ ] Quick log widget
- [ ] Mobile-optimized forms

### Notifications System
- [ ] Workout reminders (configurable times)
- [ ] Supplement reminders based on timing
- [ ] Weekly plan generation reminder
- [ ] Streak notifications
- [ ] Browser notifications support

### Export Features
- [ ] PDF workout plans for gym printing
- [ ] CSV export of session history
- [ ] Backup/restore functionality
- [ ] Share plans with coach/trainer
- [ ] Export to popular fitness apps

## Future Considerations

### Social Features
- [ ] Share achievements
- [ ] Compare with training partners
- [ ] Coach/athlete relationship support
- [ ] Community challenges

### Advanced Analytics
- [ ] Predictive analytics for performance
- [ ] Injury risk assessment
- [ ] Optimal rest day recommendations
- [ ] Training load management

### Integrations
- [ ] Wearables (Garmin, Apple Watch, Whoop)
- [ ] Calendar sync (Google, Apple)
- [ ] Nutrition apps (MyFitnessPal)
- [ ] Video analysis integration

## Technical Improvements

- [ ] Add test coverage
- [ ] API documentation
- [ ] Rate limiting
- [ ] Caching layer
- [ ] Anthropic prompt caching (when static content reaches 1,024+ tokens)
  - [ ] Restructure prompt building to separate static vs dynamic content
  - [ ] Implement cache_control breakpoints in system messages
  - [ ] Add anthropic-beta header for prompt caching
  - [ ] Expected savings: ~43% cost reduction, improved latency
- [ ] Database migrations system
- [ ] CI/CD pipeline
- [ ] Docker containerization
- [ ] Multi-language support (i18n)

---

**Last Updated:** 2026-03-14

## Notes

### Recent Changes (2026-03-14)
- Updated AI prompts to distinguish between gym/club sessions and custom workouts
  - Gym/club sessions: Minimal guidance only (attendance + optional focus tip)
  - Custom workouts: Full detailed breakdown (warm-up, main work, cool-down)
  - Prevents AI from adding unnecessary details/post-workout routines to instructor-led sessions

### Multi-Activity Support Context
The current app is optimized for single-sport, progression-focused training. Multi-activity support is a foundational enhancement needed to serve:
- Multi-sport athletes (triathletes, duathletes, adventure racers)
- Recreational multi-sport users (hiking + golf, yoga + running, climbing + cycling)
- Athletes with primary sport + supplementary training (boxing + strength training)

This enhancement requires data model changes and AI prompt restructuring to support N activities with individual goals, priorities, and training modes.
