package config

import (
	"fcm/internal/model"
	"fcm/internal/util"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Key          string                 `yaml:"key"`
	Token        string                 `yaml:"token"`
	Tokens       []string               `yaml:"tokens"`
	Topic        string                 `yaml:"topic"`
	Condition    string                 `yaml:"condition"`
	Notification *model.Notification    `yaml:"notification"`
	Data         map[string]string      `yaml:"data"`
	Android      map[string]interface{} `yaml:"android"`
	Apns         map[string]interface{} `yaml:"apns"`
	Webpush      map[string]interface{} `yaml:"webpush"`
	Log          string                 `yaml:"log"`
	Profiles     map[string]Profile     `yaml:"profiles"`
}

type Profile struct {
	Key          string                 `yaml:"key"`
	Token        string                 `yaml:"token"`
	Tokens       []string               `yaml:"tokens"`
	Topic        string                 `yaml:"topic"`
	Condition    string                 `yaml:"condition"`
	Notification *model.Notification    `yaml:"notification"`
	Data         map[string]string      `yaml:"data"`
	Android      map[string]interface{} `yaml:"android"`
	Apns         map[string]interface{} `yaml:"apns"`
	Webpush      map[string]interface{} `yaml:"webpush"`
	Log          string                 `yaml:"log"`
}

type ResolvedConfig struct {
	Key          string
	Token        string
	Tokens       []string
	Topic        string
	Condition    string
	Notification *model.Notification
	Data         map[string]string
	Android      map[string]interface{}
	Apns         map[string]interface{}
	Webpush      map[string]interface{}
	Log          string
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse yaml config: %w", err)
	}

	return &cfg, nil
}

func ResolveConfig(cfg *Config, profileName string) (ResolvedConfig, error) {
	var resolved ResolvedConfig

	if cfg == nil {
		return resolved, nil
	}

	resolved.Key = cfg.Key
	resolved.Token = cfg.Token
	resolved.Tokens = util.FirstNonEmptySlice(cfg.Tokens)
	resolved.Topic = cfg.Topic
	resolved.Condition = cfg.Condition
	resolved.Notification = cfg.Notification
	resolved.Data = util.CloneStringMap(cfg.Data)
	resolved.Android = util.CloneInterfaceMap(cfg.Android)
	resolved.Apns = util.CloneInterfaceMap(cfg.Apns)
	resolved.Webpush = util.CloneInterfaceMap(cfg.Webpush)
	resolved.Log = cfg.Log

	if profileName == "" {
		return resolved, nil
	}

	profile, ok := cfg.Profiles[profileName]
	if !ok {
		return resolved, fmt.Errorf("profile %q not found in config", profileName)
	}

	resolved.Key = util.FirstNonEmpty(profile.Key, resolved.Key)
	resolved.Token = util.FirstNonEmpty(profile.Token, resolved.Token)
	resolved.Tokens = util.FirstNonEmptySlice(profile.Tokens, resolved.Tokens)
	resolved.Topic = util.FirstNonEmpty(profile.Topic, resolved.Topic)
	resolved.Condition = util.FirstNonEmpty(profile.Condition, resolved.Condition)

	if profile.Notification != nil {
		resolved.Notification = profile.Notification
	}
	if profile.Data != nil {
		resolved.Data = util.CloneStringMap(profile.Data)
	}
	if profile.Android != nil {
		resolved.Android = util.CloneInterfaceMap(profile.Android)
	}
	if profile.Apns != nil {
		resolved.Apns = util.CloneInterfaceMap(profile.Apns)
	}
	if profile.Webpush != nil {
		resolved.Webpush = util.CloneInterfaceMap(profile.Webpush)
	}
	resolved.Log = util.FirstNonEmpty(profile.Log, resolved.Log)

	return resolved, nil
}

func WriteDefaultConfig(path string, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("file %s already exists; use --force to overwrite", path)
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return fmt.Errorf("unable to create config directory: %w", err)
	}

	content := `# FCM CLI configuration
# Keep secrets like key.json path in environment variables when possible:
# export FCM_KEY=service-account.json

notification:
  title: Hello
  body: World

data:
  env: dev

log: info

profiles:
  dev:
    token: "YOUR_DEVICE_TOKEN"
    notification:
      title: Dev notification
      body: Sent from dev profile

  prod:
    topic: production
    notification:
      title: Deploy
      body: New version released
    data:
      env: prod
      version: "1.3.0"

  smoke:
    tokens:
      - token1
      - token2
      - token3
    notification:
      title: Smoke test
      body: Batch notification test
`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("unable to write config file: %w", err)
	}

	return nil
}

func LoadDotEnv(defaultPath string) {
	_ = godotenv.Load()
	if defaultPath != "" {
		_ = godotenv.Overload(defaultPath)
	}
}
