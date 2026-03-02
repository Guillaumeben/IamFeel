package agent

import (
    "fmt"
    "strings"

    "github.com/tuxnam/iamfeel/internal/config"
)

const baseSystemPrompt = `You are an expert training coach and programming assistant. Your role is to create personalized, effective training plans based on the athlete's profile, goals, training history, and current phase.

## Your Coaching Philosophy

1. **Evidence-Based**: Base recommendations on proven training principles
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
[2-3 paragraphs explaining the overall approach for this week, considering the athlete's current phase, recent training, and goals]

WEEKLY OVERVIEW:
[Brief summary of the week's focus and session distribution]

DAILY BREAKDOWN:

Monday:
  Session: [Session name/type]
  Duration: [X minutes]
  Focus: [Primary focus of the session]
  Details: [Specific guidance - what to work on, intensity level, key exercises/drills]
  Why: [Brief rationale for this session on this day]

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

Be specific and actionable. The athlete should know exactly what to do each day.
`

// BuildSystemPrompt creates the complete system prompt with sport-specific context
func BuildSystemPrompt(sportConfig *config.SportConfig, coachingStyle string, explanationDetail string) string {
    var sb strings.Builder

    sb.WriteString(baseSystemPrompt)
    sb.WriteString("\n\n")

    // Add sport-specific context
    sb.WriteString("## Sport-Specific Context\n\n")
    sb.WriteString(fmt.Sprintf("**Sport**: %s (%s)\n\n", sportConfig.SportName, sportConfig.SportType))

    // Add agent context from sport config
    if sportConfig.AgentContext != "" {
        sb.WriteString("### Coaching Guidelines for This Sport\n\n")
        sb.WriteString(sportConfig.AgentContext)
        sb.WriteString("\n\n")
    }

    // Add available session types
    sb.WriteString("### Available Session Types\n\n")
    for _, st := range sportConfig.SessionTypes {
        sb.WriteString(fmt.Sprintf("**%s** (%s)\n", st.Name, st.ID))
        sb.WriteString(fmt.Sprintf("- Description: %s\n", st.Description))
        sb.WriteString(fmt.Sprintf("- Typical Duration: %s\n", st.TypicalDurationMinutes))
        sb.WriteString(fmt.Sprintf("- Intensity: %s\n", st.Intensity))
        sb.WriteString(fmt.Sprintf("- Primary Adaptation: %s\n", st.PrimaryAdaptation))
        sb.WriteString(fmt.Sprintf("- Recovery Impact: %s\n", st.RecoveryImpact))
        if st.Notes != "" {
            sb.WriteString(fmt.Sprintf("- Notes: %s\n", st.Notes))
        }
        sb.WriteString("\n")
    }

    // Add recovery guidelines
    sb.WriteString("### Recovery Guidelines\n\n")
    sb.WriteString("Consider these recovery requirements when scheduling sessions:\n\n")
    for intensity, recovery := range sportConfig.RecoveryGuidelines {
        sb.WriteString(fmt.Sprintf("**%s**: %d days before next hard session, %d days before same type",
            intensity, recovery.RestBeforeNextHardSession, recovery.RestBeforeSameType))
        if recovery.Notes != "" {
            sb.WriteString(fmt.Sprintf(" (%s)", recovery.Notes))
        }
        sb.WriteString("\n")
    }
    sb.WriteString("\n")

    // Add training phases
    sb.WriteString("### Training Phases\n\n")
    for _, phase := range sportConfig.Phases {
        sb.WriteString(fmt.Sprintf("**%s**\n", phase.DisplayName))
        sb.WriteString(fmt.Sprintf("- Focus: %s\n", phase.Focus))
        sb.WriteString("- Priorities:\n")
        for _, priority := range phase.Priorities {
            sb.WriteString(fmt.Sprintf("  - %s\n", priority))
        }
        sb.WriteString("\n")
    }

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
    sb.WriteString(fmt.Sprintf("- Experience Level: %s\n", userContext.ExperienceLevel))
    sb.WriteString(fmt.Sprintf("- Sport Experience: %d years\n", userContext.SportExperienceYears))
    sb.WriteString("\n")

    // Current phase
    sb.WriteString("## Current Training Phase\n\n")
    sb.WriteString(fmt.Sprintf("- Phase: %s\n", userContext.CurrentPhaseName))
    sb.WriteString(fmt.Sprintf("- Started: %s\n", userContext.PhaseStartDate))
    sb.WriteString(fmt.Sprintf("- Ends: %s\n", userContext.PhaseEndDate))
    sb.WriteString("\n")

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

        sb.WriteString("Use this information to understand what sessions are available and when they occur. Incorporate them into the plan when they align with training goals. Consider cost when suggesting optional sessions.\n")
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
    sb.WriteString(fmt.Sprintf("- Intensity Preference: %s\n", userContext.IntensityPreference))
    sb.WriteString(fmt.Sprintf("- Recovery Priority: %s\n", userContext.RecoveryPriority))
    if userContext.SessionDurationPreference != "" {
        sb.WriteString(fmt.Sprintf("- Session Duration: %s\n", userContext.SessionDurationPreference))
    }
    sb.WriteString("\n")

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
