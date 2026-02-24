package config

import (
    "fmt"
    "os"
    "path/filepath"

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

// LoadUserConfig loads a user configuration from a YAML file
func LoadUserConfig(path string) (*UserConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read user config file: %w", err)
    }

    var config UserConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse user config: %w", err)
    }

    // Basic validation
    if config.User.Name == "" {
        return nil, fmt.Errorf("user name is required")
    }
    if len(config.Sports) == 0 {
        return nil, fmt.Errorf("at least one sport is required")
    }

    // Validate experience level
    validLevels := map[string]bool{
        "beginner":     true,
        "intermediate": true,
        "advanced":     true,
    }
    if !validLevels[config.User.ExperienceLevel] {
        return nil, fmt.Errorf("invalid experience_level: %s (must be beginner, intermediate, or advanced)", config.User.ExperienceLevel)
    }

    return &config, nil
}

// SaveUserConfig saves a user configuration to a YAML file
func SaveUserConfig(path string, config *UserConfig) error {
    data, err := yaml.Marshal(config)
    if err != nil {
        return fmt.Errorf("failed to marshal user config: %w", err)
    }

    if err := os.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("failed to write user config file: %w", err)
    }

    return nil
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

// GetUserConfigPath returns the path for a user's config file
func GetUserConfigPath(userID int) string {
    return filepath.Join("data", fmt.Sprintf("user_%d_config.yaml", userID))
}

// LoadUserConfigByID loads a user configuration by user ID
func LoadUserConfigByID(userID int) (*UserConfig, error) {
    path := GetUserConfigPath(userID)
    return LoadUserConfig(path)
}

// SaveUserConfigByID saves a user configuration by user ID
func SaveUserConfigByID(userID int, config *UserConfig) error {
    path := GetUserConfigPath(userID)
    return SaveUserConfig(path, config)
}

// UserConfigExists checks if a config file exists for a user
func UserConfigExists(userID int) bool {
    path := GetUserConfigPath(userID)
    _, err := os.Stat(path)
    return err == nil
}
