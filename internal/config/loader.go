package config

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

// LoadSportConfig loads a sport configuration from a YAML file
func LoadSportConfig(path string) (*SportConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read sport config file: %w", err)
    }

    var config SportConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse sport config: %w", err)
    }

    // Basic validation
    if config.SportName == "" {
        return nil, fmt.Errorf("sport_name is required")
    }
    if len(config.Phases) == 0 {
        return nil, fmt.Errorf("at least one phase is required")
    }
    if len(config.SessionTypes) == 0 {
        return nil, fmt.Errorf("at least one session type is required")
    }

    return &config, nil
}


// GetSessionType finds a session type by ID in a sport config
func (sc *SportConfig) GetSessionType(id string) *SessionType {
    for i := range sc.SessionTypes {
        if sc.SessionTypes[i].ID == id {
            return &sc.SessionTypes[i]
        }
    }
    return nil
}

// GetPhase finds a phase by name in a sport config
func (sc *SportConfig) GetPhase(name string) *Phase {
    for i := range sc.Phases {
        if sc.Phases[i].Name == name {
            return &sc.Phases[i]
        }
    }
    return nil
}

// GetPrimarySport returns the primary sport configuration for the user
func (uc *UserConfig) GetPrimarySport() *UserSport {
    for i := range uc.Sports {
        if uc.Sports[i].Primary {
            return &uc.Sports[i]
        }
    }
    // If no primary is set, return the first sport
    if len(uc.Sports) > 0 {
        return &uc.Sports[0]
    }
    return nil
}

