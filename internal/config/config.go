package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	TicketPatterns []string `json:"ticket_patterns"`
	compiled       []*regexp.Regexp
}

func Default() *Config {
	cfg := &Config{
		TicketPatterns: []string{
			`^#\d+$`,       // GitHub issues: #123
			`^[A-Z]+-\d+$`, // Jira/Linear style: PIP-1234, INFRA-124
			`^[A-Z]+_\d+$`, // Underscore variant: PIP_1234
		},
	}
	cfg.compile()
	return cfg
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Default(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Merge with defaults if no patterns specified
	if len(cfg.TicketPatterns) == 0 {
		cfg.TicketPatterns = Default().TicketPatterns
	}

	cfg.compile()
	return &cfg, nil
}

func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func (c *Config) IsTicket(s string) bool {
	if c.compiled == nil {
		c.compile()
	}

	for _, re := range c.compiled {
		if re.MatchString(s) {
			return true
		}
	}
	return false
}

func (c *Config) compile() {
	c.compiled = make([]*regexp.Regexp, 0, len(c.TicketPatterns))
	for _, pattern := range c.TicketPatterns {
		if re, err := regexp.Compile(pattern); err == nil {
			c.compiled = append(c.compiled, re)
		}
	}
}

func getConfigPath() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "branch", "config.json"), nil
}
