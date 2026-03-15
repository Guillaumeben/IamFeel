package agent

import (
    "fmt"
    "strings"
)

const baseSystemPrompt = `You are an expert training coach and programming assistant. Your role is to create personalized, effective training plans based on the athlete's profile, goals, training history, and current phase.

## Your Coaching Philosophy

1. **Evidence-Based**: Base recommendations on proven training principles and scientific research. Your training recommendations should be backed by professional articles, scientific studies, and knowledgeable sources in the respective sport (e.g., Boxing Science for boxing, TrainingPeaks or scientific journals for endurance sports, respected strength & conditioning authorities for other disciplines).
2. **Personalized**: Adapt to the individual's experience, recovery, and goals
3. **Progressive**: Build systematically toward goals with proper progression
4. **Balanced**: Balance intensity, volume, recovery, and variety
5. **Practical**: Create realistic plans that fit the athlete's schedule and resources

## Key Principles

- **Periodization**: Align sessions with the current training phase
- **Recovery**: Respect recovery needs between high-intensity sessions
- **Variety**: Mix session types to prevent monotony and overuse
- **Specificity**: Focus on sport-specific adaptations
- **Individual Response**: Consider how the athlete has responded to previous training

## Output Format

You must respond with a structured weekly plan in the following format:

---
WEEKLY TRAINING PLAN
Week of: [dates]

RATIONALE:
[MAXIMUM 2-3 sentences explaining the overall approach for this week - be concise and focused on the key theme]

DAILY BREAKDOWN:

Monday:
  Session: [Session name/type]
  Duration: [X minutes]
  Focus: [Primary focus of the session]
  Details: [Structured breakdown - MINIMAL for gym/club sessions, DETAILED for custom workouts - see format below]
  Why: [Brief rationale - one sentence for gym/club sessions, 1-2 sentences for custom workouts]

Tuesday:
  [Same format]

[Continue for all 7 days - include rest/recovery days]

NOTES:
- [Any important considerations for the week]
- [Flexibility options if needed]
- [What to watch for (signs of fatigue, etc.)]

NEXT WEEK PREVIEW:
[Brief note on where training should progress next week]
---

## CRITICAL: Details Section Format

**IMPORTANT - Two Different Detail Levels:**

### For GYM/CLUB SESSIONS (sessions from the athlete's gym/club):
- **DO NOT** provide detailed warm-up, main work, cool-down breakdowns
- **DO NOT** add post-workout routines or additional exercises
- Keep it MINIMAL - the gym/club instructor will provide the specifics
- Format:
  Details:
  Attend [Session Name] at [Gym Name]. [Brief 1-sentence guidance on what to focus on during the session, if helpful]

- Example: "Details: Attend Boxing Sparring at EKO Boxing Club. Focus on implementing defensive techniques practiced this week."

### For CUSTOM/HOME WORKOUTS (athlete trains independently):
- Provide FULL detailed breakdown with timing
- Structure with warm-up, main work, cool-down
- Be specific and actionable

**Format for CUSTOM workouts:**

**WARM-UP (5-10 min):**
- List specific exercises with sets/reps/time
- Example: "Jump rope 3 min, arm circles 20 reps, shadow boxing 2 min"

**MAIN WORK (breakdown by segment with timing):**
- Segment 1 (X min): Specific exercises with sets/reps/rest periods
- Segment 2 (X min): Specific exercises with sets/reps/rest periods
- Continue with clear time allocations

**COOL-DOWN (5-10 min):**
- List specific cool-down activities
- Example: "Light cardio 3 min, static stretching 5 min (hamstrings, quads, shoulders)"

**TOTAL TIME CHECK:** The sum of all segments MUST match the stated Duration. If Duration is 45 minutes, your breakdown must add up to approximately 45 minutes. Account for rest periods and transitions.

**Example of GOOD structure for a 45-minute CUSTOM session:**

Details:
**WARM-UP (8 min):**
- Jump rope: 3 min
- Dynamic stretching: 5 min (leg swings, arm circles, hip rotations)

**MAIN WORK (32 min):**
- **Block 1 — Explosive Power (12 min):** 3 rounds, 60sec rest between rounds
  * Slamball slams: 12 reps
  * Ab wheel rollouts: 10 reps
  * Pallof press with band: 10 reps/side
- **Block 2 — Core Stability (10 min):** 3 rounds, 45sec rest between rounds
  * Dead bugs: 12 reps
  * Plank with shoulder taps: 45sec
  * Russian twists: 20 reps
- **Block 3 — Finisher (10 min):** Tabata format (20sec work, 10sec rest, 8 rounds)
  * Burpees and mountain climbers alternating

**COOL-DOWN (5 min):**
- Walking/light movement: 2 min
- Static stretching: 3 min (focus on core, shoulders, hips)

TOTAL: 8 + 32 + 5 = 45 minutes ✓
`

// BuildSystemPrompt creates the complete system prompt with coaching preferences
func BuildSystemPrompt(coachingStyle string, explanationDetail string) string {
    var sb strings.Builder

    sb.WriteString(baseSystemPrompt)
    sb.WriteString("\n\n")

    // Add coaching style adjustment
    sb.WriteString("## Coaching Style\n\n")
    switch coachingStyle {
    case "aggressive":
        sb.WriteString("Your coaching style is AGGRESSIVE: Push hard, high volume, challenge the athlete. Assume they can handle more.\n")
    case "supportive":
        sb.WriteString("Your coaching style is SUPPORTIVE: Be encouraging, focus on building confidence, ensure adequate recovery.\n")
    default: // balanced
        sb.WriteString("Your coaching style is BALANCED: Find the middle ground between pushing hard and ensuring recovery.\n")
    }
    sb.WriteString("\n")

    // Add explanation detail level
    sb.WriteString("## Explanation Detail\n\n")
    switch explanationDetail {
    case "brief":
        sb.WriteString("Keep explanations concise and to-the-point. Focus on what to do.\n")
    case "detailed":
        sb.WriteString("Provide detailed explanations with scientific rationale and progression logic.\n")
    default: // moderate
        sb.WriteString("Provide moderate detail - explain the why but stay focused.\n")
    }

    return sb.String()
}

// BuildUserPrompt creates the user prompt with all context
func BuildUserPrompt(userContext *UserContext) string {
    var sb strings.Builder

    sb.WriteString("Please create a weekly training plan for the following athlete:\n\n")

    // User profile
    sb.WriteString("## Athlete Profile\n\n")
    sb.WriteString(fmt.Sprintf("- Name: %s\n", userContext.Name))
    sb.WriteString(fmt.Sprintf("- Age: %d\n", userContext.Age))
    sb.WriteString(fmt.Sprintf("- General Fitness Level: %s\n", userContext.ExperienceLevel))
    sb.WriteString("\n")

    // Activities section - supports multi-sport
    if len(userContext.Activities) > 1 {
        sb.WriteString("## Activities\n\n")
        sb.WriteString("This athlete practices multiple activities. Balance training time and recovery across all activities based on their priorities and goals.\n\n")

        for _, activity := range userContext.Activities {
            sb.WriteString(fmt.Sprintf("### %s (%s Priority - %s)\n", activity.Name,
                strings.ToUpper(activity.Priority), formatGoalType(activity.GoalType)))
            sb.WriteString(fmt.Sprintf("- Experience: %d years", activity.ExperienceYears))
            if activity.ExperienceYears == 0 {
                sb.WriteString(" (NEW to this activity - focus on fundamentals)")
            }
            sb.WriteString("\n")
            if activity.CurrentPhase != "" {
                sb.WriteString(fmt.Sprintf("- Current Phase: %s\n", activity.CurrentPhase))
                if activity.PhaseStartDate != "" && activity.PhaseEndDate != "" {
                    sb.WriteString(fmt.Sprintf("- Phase Duration: %s to %s\n", activity.PhaseStartDate, activity.PhaseEndDate))
                }
            }
            if activity.TargetSessionsPerWeek > 0 {
                sb.WriteString(fmt.Sprintf("- Target: %.1f sessions/week\n", activity.TargetSessionsPerWeek))
            }
            if activity.Notes != "" {
                sb.WriteString(fmt.Sprintf("- Notes: %s\n", activity.Notes))
            }
            sb.WriteString("\n")
        }

        sb.WriteString("**CRITICAL - Multi-Activity Balancing:**\n")
        sb.WriteString("- HIGH priority activities: These are the main focus - allocate most training time here\n")
        sb.WriteString("- MEDIUM priority activities: Maintain current level - steady state training\n")
        sb.WriteString("- LOW priority activities: Fit in when schedule allows - recreational/recovery focus\n")
        sb.WriteString("- Respect target sessions per week for each activity\n")
        sb.WriteString("- COMPETITION PREP activities need progressive, periodized training\n")
        sb.WriteString("- MAINTENANCE activities need consistent but not progressive training\n")
        sb.WriteString("- LEARNING activities need gradual skill building and technique focus\n")
        sb.WriteString("- RECREATION activities are for enjoyment and active recovery\n\n")
    } else if len(userContext.Activities) == 1 {
        // Single activity - simpler format
        activity := userContext.Activities[0]
        sb.WriteString(fmt.Sprintf("- Primary Activity: %s\n", activity.Name))
        sb.WriteString(fmt.Sprintf("- Experience: %d years\n", activity.ExperienceYears))
        if activity.ExperienceYears == 0 {
            sb.WriteString("  **Note**: This athlete is NEW to this activity. Focus on teaching fundamentals, proper technique, and gradual progression.\n")
        }
        sb.WriteString("\n")
    } else {
        // Fallback to deprecated single-sport fields
        sb.WriteString(fmt.Sprintf("- Experience in %s: %d years\n", userContext.SportName, userContext.SportExperienceYears))
        if userContext.SportExperienceYears == 0 {
            sb.WriteString("  **Note**: This athlete is NEW to this specific sport, despite having general fitness experience. Focus on teaching fundamentals, proper technique, and gradual progression.\n")
        }
        sb.WriteString("\n")
    }

    // Sport-specific scientific references
    sb.WriteString("## Scientific Basis for Training\n\n")
    if len(userContext.Activities) > 1 {
        sb.WriteString("Draw from professional sources for each activity:\n")
        seen := make(map[string]bool)
        for _, activity := range userContext.Activities {
            if !seen[activity.Name] {
                sb.WriteString(fmt.Sprintf("- **%s**: %s\n", activity.Name, getSportReferences(activity.Name)))
                seen[activity.Name] = true
            }
        }
        sb.WriteString("\n")
    } else if len(userContext.Activities) == 1 {
        sb.WriteString(fmt.Sprintf("For %s training, draw from these professional sources:\n", userContext.Activities[0].Name))
        sb.WriteString(fmt.Sprintf("- %s\n\n", getSportReferences(userContext.Activities[0].Name)))
    } else {
        // Fallback
        sb.WriteString(fmt.Sprintf("For %s training, draw from these professional sources:\n", userContext.SportName))
        sb.WriteString(fmt.Sprintf("- %s\n\n", getSportReferences(userContext.SportName)))
    }
    sb.WriteString("Your recommendations should be grounded in these evidence-based resources and current best practices in the field.\n\n")

    // Current phase (only show if single activity or primary sport)
    if len(userContext.Activities) == 1 {
        activity := userContext.Activities[0]
        if activity.CurrentPhase != "" {
            sb.WriteString("## Current Training Phase\n\n")
            sb.WriteString(fmt.Sprintf("- Phase: %s\n", activity.CurrentPhase))
            if activity.PhaseStartDate != "" {
                sb.WriteString(fmt.Sprintf("- Started: %s\n", activity.PhaseStartDate))
            }
            if activity.PhaseEndDate != "" {
                sb.WriteString(fmt.Sprintf("- Ends: %s\n", activity.PhaseEndDate))
            }
            sb.WriteString("\n")
        }
    } else if len(userContext.Activities) == 0 && userContext.CurrentPhaseName != "" {
        // Fallback to deprecated fields
        sb.WriteString("## Current Training Phase\n\n")
        sb.WriteString(fmt.Sprintf("- Phase: %s\n", userContext.CurrentPhaseName))
        sb.WriteString(fmt.Sprintf("- Started: %s\n", userContext.PhaseStartDate))
        sb.WriteString(fmt.Sprintf("- Ends: %s\n", userContext.PhaseEndDate))
        sb.WriteString("\n")
    }

    // Goals
    if len(userContext.ShortTermGoals) > 0 || len(userContext.MediumTermGoals) > 0 || len(userContext.LongTermGoals) > 0 {
        sb.WriteString("## Goals\n\n")
        if len(userContext.ShortTermGoals) > 0 {
            sb.WriteString("**Short-term (1-3 months):**\n")
            for _, goal := range userContext.ShortTermGoals {
                sb.WriteString(fmt.Sprintf("- %s\n", goal))
            }
            sb.WriteString("\n")
        }
        if len(userContext.MediumTermGoals) > 0 {
            sb.WriteString("**Medium-term (3-6 months):**\n")
            for _, goal := range userContext.MediumTermGoals {
                sb.WriteString(fmt.Sprintf("- %s\n", goal))
            }
            sb.WriteString("\n")
        }
        if len(userContext.LongTermGoals) > 0 {
            sb.WriteString("**Long-term (6-12 months):**\n")
            for _, goal := range userContext.LongTermGoals {
                sb.WriteString(fmt.Sprintf("- %s\n", goal))
            }
            sb.WriteString("\n")
        }
    }

    // Recent training history
    if len(userContext.RecentSessions) > 0 {
        sb.WriteString("## Recent Training History (Last 14 Days)\n\n")
        for _, session := range userContext.RecentSessions {
            sb.WriteString(fmt.Sprintf("- %s: %s (%d min, effort %d/10)\n",
                session.Date, session.SessionType, session.DurationMinutes, session.PerceivedEffort))
            if session.Notes != "" {
                sb.WriteString(fmt.Sprintf("  Notes: %s\n", session.Notes))
            }
        }
        sb.WriteString("\n")
    } else {
        sb.WriteString("## Recent Training History\n\nNo recent training sessions logged.\n\n")
    }

    // Scheduled club sessions
    if len(userContext.ClubSessions) > 0 {
        sb.WriteString("## Available Club/Gym Sessions\n\n")
        sb.WriteString("The athlete has access to these sessions at their gym(s)/club(s). Consider incorporating them when appropriate:\n\n")
        for _, session := range userContext.ClubSessions {
            sb.WriteString(fmt.Sprintf("**%s** at %s\n", session.Name, session.GymName))
            if session.Description != "" {
                sb.WriteString(fmt.Sprintf("- Description: %s\n", session.Description))
            }
            if session.Occurrences != "" {
                sb.WriteString(fmt.Sprintf("- Schedule: %s\n", session.Occurrences))
            }
            if session.Cost != "" {
                sb.WriteString(fmt.Sprintf("- Cost: %s\n", session.Cost))
            }
            sb.WriteString("\n")
        }

        sb.WriteString("**IMPORTANT**: When scheduling these gym/club sessions in the plan:\n")
        sb.WriteString("- Use MINIMAL detail format (see Details Section Format above)\n")
        sb.WriteString("- Simply reference the session and gym/club name\n")
        sb.WriteString("- Optionally add one brief sentence of guidance on what to focus on\n")
        sb.WriteString("- DO NOT create detailed workout breakdowns - the instructor will handle that\n")
        sb.WriteString("- DO NOT add pre/post-workout routines to these sessions\n")
        sb.WriteString("- Consider cost when suggesting optional sessions\n")
        sb.WriteString("\n")
    }

    // Weekly availability
    sb.WriteString("## Weekly Availability\n\n")
    for day, avail := range userContext.Availability {
        times := []string{}
        if avail.Morning {
            times = append(times, "morning")
        }
        if avail.Lunch {
            times = append(times, "lunch")
        }
        if avail.Evening {
            times = append(times, "evening")
        }
        if len(times) > 0 {
            sb.WriteString(fmt.Sprintf("- %s: %s", day, strings.Join(times, ", ")))
        } else {
            sb.WriteString(fmt.Sprintf("- %s: not available", day))
        }
        if avail.PreferredLocation != "" {
            sb.WriteString(fmt.Sprintf(" [prefers: %s]", avail.PreferredLocation))
        }
        if avail.Notes != "" {
            sb.WriteString(fmt.Sprintf(" (%s)", avail.Notes))
        }
        sb.WriteString("\n")
    }
    sb.WriteString("\n")
    sb.WriteString("**Location Preferences:** Use the preferred location hints when suggesting workouts. Home = use home equipment, Gym = suggest gym/club sessions, Flexible = either works.\n\n")

    // Equipment
    if len(userContext.AvailableEquipment) > 0 {
        sb.WriteString("## Available Equipment\n\n")
        for _, eq := range userContext.AvailableEquipment {
            sb.WriteString(fmt.Sprintf("- %s\n", eq))
        }
        sb.WriteString("\n")
    }

    // Training preferences
    sb.WriteString("## Training Preferences\n\n")
    if userContext.PrimaryGoal != "" {
        sb.WriteString(fmt.Sprintf("- Primary Goal: %s\n", userContext.PrimaryGoal))
    }
    if userContext.SessionsPerWeek > 0 {
        sb.WriteString(fmt.Sprintf("- Target Sessions Per Week: %.1f\n", userContext.SessionsPerWeek))
    }
    sb.WriteString(fmt.Sprintf("- Intensity Preference: %s\n", userContext.IntensityPreference))
    sb.WriteString(fmt.Sprintf("- Recovery Priority: %s\n", userContext.RecoveryPriority))
    if userContext.SessionDurationPreference != "" {
        sb.WriteString(fmt.Sprintf("- Session Duration: %s\n", userContext.SessionDurationPreference))
    }
    if userContext.MaxSessionsPerDay > 0 {
        sb.WriteString(fmt.Sprintf("- **Maximum Sessions Per Day: %d**", userContext.MaxSessionsPerDay))
        if !userContext.AllowShortSessions {
            sb.WriteString(" (strict - no additional short exercise sessions)")
        } else {
            sb.WriteString(" (main coached sessions only - additional short supplementary exercise sessions are allowed if beneficial)")
        }
        sb.WriteString("\n")
    }
    sb.WriteString("\n")
    sb.WriteString("**CRITICAL CONSTRAINT - Maximum Sessions Per Day:**\n")
    if userContext.MaxSessionsPerDay == 1 {
        sb.WriteString("- Schedule ONLY ONE session per day in the DAILY BREAKDOWN\n")
        sb.WriteString("- Each day can have AT MOST ONE entry (Monday/Tuesday/etc.) with ONE session\n")
        sb.WriteString("- Do NOT use '+' or list multiple session names on the same day\n")
        sb.WriteString("- INCORRECT: 'Thursday: Bag Work + Ringcraft'\n")
        sb.WriteString("- CORRECT: 'Thursday: Bag Work' (choose the most important session)\n")
    } else {
        sb.WriteString(fmt.Sprintf("- Schedule AT MOST %d sessions per day in the DAILY BREAKDOWN\n", userContext.MaxSessionsPerDay))
        sb.WriteString(fmt.Sprintf("- Each day can have AT MOST %d session entries listed\n", userContext.MaxSessionsPerDay))
    }

    if userContext.AllowShortSessions {
        sb.WriteString("- Rest days can include a short supplementary exercise session (10-20 minutes) if beneficial for recovery, mobility, or skill maintenance\n")
    } else {
        sb.WriteString("- Rest days should explicitly state 'Rest' or 'Recovery' with no additional sessions\n")
    }
    sb.WriteString("\n")

    // Special constraints/requests from user (if provided)
    if userContext.SpecialConstraints != "" {
        sb.WriteString("## Special Requests & Constraints\n\n")
        sb.WriteString(userContext.SpecialConstraints)
        sb.WriteString("\n\n")
        sb.WriteString("**IMPORTANT**: Please take the above constraints into account when creating this week's plan.\n\n")
    }

    // Request
    sb.WriteString("## Request\n\n")
    sb.WriteString(fmt.Sprintf("Create a weekly training plan for the week of %s.\n", userContext.WeekStart))
    sb.WriteString("Follow the output format specified in the system prompt exactly.\n")

    return sb.String()
}

// BuildDayModificationPrompt creates a prompt for modifying a single day
func BuildDayModificationPrompt(userContext *UserContext, currentPlan string, dayName string, modificationReason string) string {
    var sb strings.Builder

    sb.WriteString("I need to modify the training plan for a specific day in the current week's plan.\n\n")

    // User profile (condensed)
    sb.WriteString("## Athlete Profile\n\n")
    sb.WriteString(fmt.Sprintf("- Name: %s\n", userContext.Name))
    sb.WriteString(fmt.Sprintf("- Experience Level: %s\n", userContext.ExperienceLevel))
    sb.WriteString(fmt.Sprintf("- Current Phase: %s\n", userContext.CurrentPhaseName))
    sb.WriteString("\n")

    // Current plan
    sb.WriteString("## Current Week's Plan\n\n")
    sb.WriteString("```\n")
    sb.WriteString(currentPlan)
    sb.WriteString("\n```\n\n")

    // Modification request
    sb.WriteString("## Modification Request\n\n")
    sb.WriteString(fmt.Sprintf("**Day to modify**: %s\n", dayName))
    sb.WriteString(fmt.Sprintf("**Reason for change**: %s\n\n", modificationReason))

    // Recent training context
    if len(userContext.RecentSessions) > 0 {
        sb.WriteString("## Recent Training (Last 14 Days)\n\n")
        for _, session := range userContext.RecentSessions {
            sb.WriteString(fmt.Sprintf("- %s: %s (%d min, effort %d/10)\n",
                session.Date, session.SessionType, session.DurationMinutes, session.PerceivedEffort))
        }
        sb.WriteString("\n")
    }

    // Availability for that day
    if avail, ok := userContext.Availability[dayName]; ok {
        sb.WriteString(fmt.Sprintf("## %s Availability\n\n", dayName))
        times := []string{}
        if avail.Morning {
            times = append(times, "morning")
        }
        if avail.Lunch {
            times = append(times, "lunch")
        }
        if avail.Evening {
            times = append(times, "evening")
        }
        if len(times) > 0 {
            sb.WriteString(fmt.Sprintf("- Available: %s\n", strings.Join(times, ", ")))
        }
        if avail.PreferredLocation != "" {
            sb.WriteString(fmt.Sprintf("- Preferred location: %s\n", avail.PreferredLocation))
        }
        if avail.Notes != "" {
            sb.WriteString(fmt.Sprintf("- Notes: %s\n", avail.Notes))
        }
        sb.WriteString("\n")
    }

    // Request
    sb.WriteString("## Your Task\n\n")
    sb.WriteString(fmt.Sprintf("Please regenerate ONLY the %s session in the plan, taking into account:\n", dayName))
    sb.WriteString("1. The modification reason provided above\n")
    sb.WriteString("2. The overall structure and goals of the current week's plan\n")
    sb.WriteString("3. The athlete's recent training load\n")
    sb.WriteString("4. Proper recovery and progression principles\n\n")

    sb.WriteString("Provide the modified day in the same format as the original plan:\n\n")
    sb.WriteString(fmt.Sprintf("%s:\n", dayName))
    sb.WriteString("  Session: [Session name/type]\n")
    sb.WriteString("  Duration: [X minutes]\n")
    sb.WriteString("  Focus: [Primary focus]\n")
    sb.WriteString("  Details: [Specific guidance]\n")
    sb.WriteString("  Why: [Brief rationale for this modified session]\n")

    return sb.String()
}

// BuildPlanAdjustmentPrompt creates a prompt for adjusting a weekly plan
func BuildPlanAdjustmentPrompt(userContext *UserContext, previousPlan string, adjustmentNotes string) string {
    var sb strings.Builder

    sb.WriteString("I need to create a new weekly training plan based on a previous week's plan with some adjustments.\n\n")

    // User profile (condensed)
    sb.WriteString("## Athlete Profile\n\n")
    sb.WriteString(fmt.Sprintf("- Name: %s\n", userContext.Name))
    sb.WriteString(fmt.Sprintf("- Experience Level: %s\n", userContext.ExperienceLevel))
    sb.WriteString(fmt.Sprintf("- Current Phase: %s\n", userContext.CurrentPhaseName))
    sb.WriteString("\n")

    // Previous plan
    sb.WriteString("## Previous Week's Plan\n\n")
    sb.WriteString("This is the plan from the previous week that should be used as a baseline:\n\n")
    sb.WriteString("```\n")
    sb.WriteString(previousPlan)
    sb.WriteString("\n```\n\n")

    // Adjustment request
    sb.WriteString("## Adjustment Request\n\n")
    sb.WriteString(fmt.Sprintf("**Requested adjustments**: %s\n\n", adjustmentNotes))

    // Recent training context
    if len(userContext.RecentSessions) > 0 {
        sb.WriteString("## Recent Training (Last 14 Days)\n\n")
        for _, session := range userContext.RecentSessions {
            sb.WriteString(fmt.Sprintf("- %s: %s (%d min, effort %d/10)\n",
                session.Date, session.SessionType, session.DurationMinutes, session.PerceivedEffort))
            if session.Notes != "" {
                sb.WriteString(fmt.Sprintf("  Notes: %s\n", session.Notes))
            }
        }
        sb.WriteString("\n")
    }

    // Weekly availability
    sb.WriteString("## Weekly Availability\n\n")
    for day, avail := range userContext.Availability {
        times := []string{}
        if avail.Morning {
            times = append(times, "morning")
        }
        if avail.Lunch {
            times = append(times, "lunch")
        }
        if avail.Evening {
            times = append(times, "evening")
        }
        if len(times) > 0 {
            sb.WriteString(fmt.Sprintf("- %s: %s", day, strings.Join(times, ", ")))
        } else {
            sb.WriteString(fmt.Sprintf("- %s: not available", day))
        }
        if avail.PreferredLocation != "" {
            sb.WriteString(fmt.Sprintf(" [prefers: %s]", avail.PreferredLocation))
        }
        sb.WriteString("\n")
    }
    sb.WriteString("\n")

    // Request
    sb.WriteString("## Your Task\n\n")
    sb.WriteString(fmt.Sprintf("Please create a NEW weekly training plan for the week of %s, based on:\n", userContext.WeekStart))
    sb.WriteString("1. The previous week's plan structure (use it as a template)\n")
    sb.WriteString("2. The adjustment notes provided above\n")
    sb.WriteString("3. The athlete's recent training load and recovery status\n")
    sb.WriteString("4. Proper progression principles (don't simply repeat - adapt and progress)\n")
    sb.WriteString("5. The athlete's current availability\n\n")
    sb.WriteString("Follow the standard weekly plan output format specified in the system prompt.\n")
    sb.WriteString("Make sure to incorporate the requested adjustments while maintaining a balanced and progressive training approach.\n")

    return sb.String()
}

// formatGoalType converts goal type to human-readable format
func formatGoalType(goalType string) string {
    switch goalType {
    case "competition_prep":
        return "Competition Preparation"
    case "maintenance":
        return "Maintenance"
    case "learning":
        return "Learning/Development"
    case "recreation":
        return "Recreation"
    default:
        return goalType
    }
}

// getSportReferences returns professional/scientific source references for a sport
func getSportReferences(sportName string) string {
    switch strings.ToLower(sportName) {
    case "boxing":
        return "Boxing Science (boxingscience.co.uk), peer-reviewed sports science journals on combat sports training"
    case "running":
        return "Running science research (e.g., Journal of Applied Physiology), TrainingPeaks resources, coaching authorities like Jack Daniels' Running Formula"
    case "cycling":
        return "TrainingPeaks, British Cycling resources, peer-reviewed cycling performance research"
    case "swimming":
        return "Swimming science research, USA Swimming resources, peer-reviewed aquatic sports physiology"
    case "bjj", "brazilian jiu-jitsu":
        return "Scientific research on grappling sports, strength & conditioning for combat sports, BJJ-specific performance resources"
    case "crossfit":
        return "CrossFit Journal, sports science research on high-intensity functional training, NSCA guidelines"
    case "weightlifting", "strength training", "fitness":
        return "NSCA (National Strength & Conditioning Association), peer-reviewed strength training research, evidence-based programming resources"
    default:
        return "Peer-reviewed sports science research, respected coaching authorities in the field, evidence-based training methodologies"
    }
}
