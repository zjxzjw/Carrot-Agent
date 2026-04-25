package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Agent    AgentConfig    `yaml:"agent"`
	Model    ModelConfig    `yaml:"model"`
	Storage  StorageConfig  `yaml:"storage"`
	Server   ServerConfig   `yaml:"server"`
	Security SecurityConfig `yaml:"security"`
	Auth     AuthConfig     `yaml:"auth"`
}

type AgentConfig struct {
	Name          string `yaml:"name"`
	Version       string `yaml:"version"`
	DataDir       string `yaml:"data_dir"`
	LogLevel      string `yaml:"log_level"`
	SkillNudgeInt int    `yaml:"skill_nudge_interval"`
}

type ModelConfig struct {
	Provider    string  `yaml:"provider"`
	APIKey      string  `yaml:"api_key"`
	ModelName   string  `yaml:"model_name"`
	BaseURL     string  `yaml:"base_url"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
}

type StorageConfig struct {
	DBPath     string `yaml:"db_path"`
	SkillDir   string `yaml:"skill_dir"`
	MemoryDir  string `yaml:"memory_dir"`
	SessionDir string `yaml:"session_dir"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type SecurityConfig struct {
	AllowedPaths []string `yaml:"allowed_paths"`
	BlockedCmds  []string `yaml:"blocked_cmds"`
}

type AuthConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (c *Config) Validate() []error {
	var errors []error

	if c.Model.APIKey == "" {
		errors = append(errors, fmt.Errorf("model API key is required"))
	}

	if c.Model.ModelName == "" {
		errors = append(errors, fmt.Errorf("model name is required"))
	}

	if c.Model.Temperature < 0 || c.Model.Temperature > 2 {
		errors = append(errors, fmt.Errorf("temperature must be between 0 and 2, got %f", c.Model.Temperature))
	}

	if c.Model.MaxTokens <= 0 || c.Model.MaxTokens > 128000 {
		errors = append(errors, fmt.Errorf("max_tokens must be between 1 and 128000, got %d", c.Model.MaxTokens))
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		errors = append(errors, fmt.Errorf("server port must be between 1 and 65535, got %d", c.Server.Port))
	}

	if c.Agent.SkillNudgeInt < 0 {
		errors = append(errors, fmt.Errorf("skill_nudge_interval cannot be negative, got %d", c.Agent.SkillNudgeInt))
	}

	if c.Auth.Username == "" {
		errors = append(errors, fmt.Errorf("auth username is required"))
	}

	if c.Auth.Password == "" {
		errors = append(errors, fmt.Errorf("auth password is required"))
	}

	return errors
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if errs := cfg.Validate(); len(errs) > 0 {
		return nil, fmt.Errorf("configuration validation failed: %v", errs)
	}

	return &cfg, nil
}

func Default() *Config {
	return &Config{
		Agent: AgentConfig{
			Name:          "carrot-agent",
			Version:       "0.1.0",
			DataDir:       "~/.carrot",
			LogLevel:      "info",
			SkillNudgeInt: 10,
		},
		Model: ModelConfig{
			Provider:    "openai",
			ModelName:   "gpt-4",
			Temperature: 0.7,
			MaxTokens:   4096,
		},
		Storage: StorageConfig{
			DBPath:     "~/.carrot/data/carrot.db",
			SkillDir:   "~/.carrot/skills",
			MemoryDir:  "~/.carrot/memories",
			SessionDir: "~/.carrot/sessions",
		},
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
			Mode: "cli",
		},
		Security: SecurityConfig{
			AllowedPaths: []string{"~/.carrot"},
			BlockedCmds:  []string{"rm -rf /", ":(){ :|:& };:"},
		},
		Auth: AuthConfig{
			Username: "admin",
			Password: "admin123",
		},
	}
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}