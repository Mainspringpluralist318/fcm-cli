//
// SPDX-License-Identifier: MPL-2.0
//

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"gopkg.in/yaml.v3"
)

const messagingScope = "https://www.googleapis.com/auth/firebase.messaging"

var version = "dev"

// ================= LOGGER =================

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	DEBUG LogLevel = "DEBUG"
)

var currentLogLevel = INFO
var jsonLogs = false
var outputJSON = false

func log(level LogLevel, msg string, args ...interface{}) {
	if outputJSON {
		return
	}
	if level == DEBUG && currentLogLevel != DEBUG {
		return
	}

	message := fmt.Sprintf(msg, args...)

	if jsonLogs {
		entry := map[string]interface{}{
			"level":   level,
			"message": message,
			"time":    time.Now().Format(time.RFC3339),
		}
		_ = json.NewEncoder(os.Stdout).Encode(entry)
		return
	}

	color := ""
	reset := "\033[0m"

	switch level {
	case INFO:
		color = "\033[34m"
	case ERROR:
		color = "\033[31m"
	case DEBUG:
		color = "\033[33m"
	}

	fmt.Printf("%s[%s]%s %s\n", color, level, reset, message)
}

// ================= OUTPUT STRUCTS =================

type CLIResult struct {
	Success      bool              `json:"success"`
	MessageID    string            `json:"message_id,omitempty"`
	Code         int               `json:"code,omitempty"`
	Error        string            `json:"error,omitempty"`
	SuccessCount int               `json:"success_count,omitempty"`
	FailureCount int               `json:"failure_count,omitempty"`
	Results      []MulticastItem   `json:"results,omitempty"`
	Meta         map[string]string `json:"meta,omitempty"`
}

type MulticastItem struct {
	Token     string `json:"token"`
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	Code      int    `json:"code,omitempty"`
	Error     string `json:"error,omitempty"`
}

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

// ================= FCM STRUCTS =================

type FCMMessage struct {
	Message MessageBody `json:"message"`
}

type MessageBody struct {
	Token        string                 `json:"token,omitempty"`
	Topic        string                 `json:"topic,omitempty"`
	Condition    string                 `json:"condition,omitempty"`
	Notification *Notification          `json:"notification,omitempty"`
	Data         map[string]string      `json:"data,omitempty"`
	Android      map[string]interface{} `json:"android,omitempty"`
	Apns         map[string]interface{} `json:"apns,omitempty"`
	Webpush      map[string]interface{} `json:"webpush,omitempty"`
}

type Notification struct {
	Title string `json:"title" yaml:"title"`
	Body  string `json:"body" yaml:"body"`
}

type fcmSuccessResponse struct {
	Name string `json:"name"`
}

// ================= YAML CONFIG =================

type Config struct {
	Key          string                 `yaml:"key"`
	Token        string                 `yaml:"token"`
	Tokens       []string               `yaml:"tokens"`
	Topic        string                 `yaml:"topic"`
	Condition    string                 `yaml:"condition"`
	Notification *Notification          `yaml:"notification"`
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
	Notification *Notification          `yaml:"notification"`
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
	Notification *Notification
	Data         map[string]string
	Android      map[string]interface{}
	Apns         map[string]interface{}
	Webpush      map[string]interface{}
	Log          string
}

// ================= UTILS =================

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func firstNonEmptySlice(values ...[]string) []string {
	for _, v := range values {
		if len(v) > 0 {
			return v
		}
	}
	return nil
}

func cloneStringMap(in map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneInterfaceMap(in map[string]interface{}) map[string]interface{} {
	if in == nil {
		return nil
	}
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func parseTokensCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func readTokensFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read tokens file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, line)
	}

	if len(out) == 0 {
		return nil, errors.New("tokens file is empty")
	}

	return out, nil
}

func loadConfig(path string) (*Config, error) {
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

func resolveConfig(cfg *Config, profileName string) (ResolvedConfig, error) {
	var resolved ResolvedConfig

	if cfg == nil {
		return resolved, nil
	}

	resolved.Key = cfg.Key
	resolved.Token = cfg.Token
	resolved.Tokens = firstNonEmptySlice(cfg.Tokens)
	resolved.Topic = cfg.Topic
	resolved.Condition = cfg.Condition
	resolved.Notification = cfg.Notification
	resolved.Data = cloneStringMap(cfg.Data)
	resolved.Android = cloneInterfaceMap(cfg.Android)
	resolved.Apns = cloneInterfaceMap(cfg.Apns)
	resolved.Webpush = cloneInterfaceMap(cfg.Webpush)
	resolved.Log = cfg.Log

	if profileName == "" {
		return resolved, nil
	}

	profile, ok := cfg.Profiles[profileName]
	if !ok {
		return resolved, fmt.Errorf("profile %q not found in config", profileName)
	}

	resolved.Key = firstNonEmpty(profile.Key, resolved.Key)
	resolved.Token = firstNonEmpty(profile.Token, resolved.Token)
	resolved.Tokens = firstNonEmptySlice(profile.Tokens, resolved.Tokens)
	resolved.Topic = firstNonEmpty(profile.Topic, resolved.Topic)
	resolved.Condition = firstNonEmpty(profile.Condition, resolved.Condition)

	if profile.Notification != nil {
		resolved.Notification = profile.Notification
	}
	if profile.Data != nil {
		resolved.Data = cloneStringMap(profile.Data)
	}
	if profile.Android != nil {
		resolved.Android = cloneInterfaceMap(profile.Android)
	}
	if profile.Apns != nil {
		resolved.Apns = cloneInterfaceMap(profile.Apns)
	}
	if profile.Webpush != nil {
		resolved.Webpush = cloneInterfaceMap(profile.Webpush)
	}
	resolved.Log = firstNonEmpty(profile.Log, resolved.Log)

	return resolved, nil
}

func validateTargets(token string, tokens []string, topic string, condition string) error {
	count := 0
	if token != "" {
		count++
	}
	if len(tokens) > 0 {
		count++
	}
	if topic != "" {
		count++
	}
	if condition != "" {
		count++
	}

	if count == 0 {
		return fmt.Errorf("provide exactly one target: token, tokens, topic or condition")
	}
	if count > 1 {
		return fmt.Errorf("only one target may be used: token, tokens, topic or condition")
	}
	return nil
}

func writeDefaultConfig(path string, force bool) error {
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

func loadDotEnv(defaultPath string) {
	_ = godotenv.Load()
	if defaultPath != "" {
		_ = godotenv.Overload(defaultPath)
	}
}

// ================= AUTH =================

func getAccessToken(ctx context.Context, filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("unable to read key file: %w", err)
	}

	config, err := google.JWTConfigFromJSON(data, messagingScope)
	if err != nil {
		return "", fmt.Errorf("unable to parse JWT config: %w", err)
	}

	token, err := config.TokenSource(ctx).Token()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve token: %w", err)
	}

	return token.AccessToken, nil
}

func getProjectID(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("unable to read key file: %w", err)
	}

	var key map[string]interface{}
	if err := json.Unmarshal(data, &key); err != nil {
		return "", fmt.Errorf("unable to parse key file JSON: %w", err)
	}

	projectID, ok := key["project_id"].(string)
	if !ok || projectID == "" {
		return "", fmt.Errorf("project_id not found in key file")
	}

	return projectID, nil
}

// ================= HTTP =================

func sendWithRetry(ctx context.Context, url, accessToken string, msg FCMMessage, retries int) (string, int, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return "", 0, fmt.Errorf("unable to marshal message: %w", err)
	}

	var lastCode int

	for i := 0; i <= retries; i++ {
		log(DEBUG, "Sending request: %s", string(body))

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			return "", 0, fmt.Errorf("unable to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log(ERROR, "Request error: %v", err)
			if i < retries {
				backoff := time.Duration(1<<i) * time.Second
				log(DEBUG, "Retrying in %v...", backoff)
				time.Sleep(backoff)
			}
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		lastCode = resp.StatusCode

		if resp.StatusCode == http.StatusOK {
			var successResp fcmSuccessResponse
			if err := json.Unmarshal(respBody, &successResp); err != nil {
				return "", resp.StatusCode, nil
			}
			return successResp.Name, resp.StatusCode, nil
		}

		err = fmt.Errorf("FCM error %d: %s", resp.StatusCode, string(respBody))
		log(ERROR, "%v", err)

		if i < retries {
			backoff := time.Duration(1<<i) * time.Second
			log(DEBUG, "Retrying in %v...", backoff)
			time.Sleep(backoff)
		} else {
			return "", lastCode, err
		}
	}

	return "", lastCode, fmt.Errorf("failed after retries")
}

// ================= PROGRESS =================

func renderProgress(done, total int64) {
	if outputJSON {
		return
	}
	percent := float64(done) / float64(total) * 100
	fmt.Printf("\rProgress: %d/%d (%.0f%%)", done, total, percent)
}

// ================= MULTICAST =================

func sendMulticast(ctx context.Context, url, accessToken string, base MessageBody, tokens []string) CLIResult {
	var wg sync.WaitGroup
	var completed int64
	var success int64

	results := make([]MulticastItem, len(tokens))
	total := int64(len(tokens))

	for i, t := range tokens {
		wg.Add(1)

		go func(idx int, tok string) {
			defer wg.Done()

			msg := FCMMessage{Message: base}
			msg.Message.Token = tok

			messageID, code, err := sendWithRetry(ctx, url, accessToken, msg, 3)
			item := MulticastItem{
				Token: tok,
			}

			if err == nil {
				item.Success = true
				item.MessageID = messageID
				atomic.AddInt64(&success, 1)
			} else {
				item.Success = false
				item.Code = code
				item.Error = err.Error()
				log(ERROR, "Failed for token %s: %v", tok, err)
			}

			results[idx] = item

			done := atomic.AddInt64(&completed, 1)
			renderProgress(done, total)
		}(i, t)
	}

	wg.Wait()
	if !outputJSON {
		fmt.Println()
		log(INFO, "Success: %d/%d", success, total)
	}

	return CLIResult{
		Success:      atomic.LoadInt64(&success) == total,
		SuccessCount: int(success),
		FailureCount: len(tokens) - int(success),
		Results:      results,
	}
}

// ================= INIT COMMAND =================

func runInit(args []string) {
	initFlags := flag.NewFlagSet("init", flag.ExitOnError)
	fileShort := initFlags.String("f", "fcm.yaml", "")
	fileLong := initFlags.String("file", "", "")
	force := initFlags.Bool("force", false, "")

	initFlags.Usage = func() {
		fmt.Println("Usage: fcm init [options]")
		fmt.Println("Options:")
		fmt.Println("  -f, --file <file>   Output config file path (default: fcm.yaml)")
		fmt.Println("  --force             Overwrite existing file")
	}

	_ = initFlags.Parse(args)

	path := firstNonEmpty(*fileLong, *fileShort)
	if path == "" {
		path = "fcm.yaml"
	}

	if err := writeDefaultConfig(path, *force); err != nil {
		if outputJSON {
			printJSON(CLIResult{
				Success: false,
				Error:   err.Error(),
			})
		} else {
			log(ERROR, "%v", err)
		}
		os.Exit(1)
	}

	if outputJSON {
		printJSON(CLIResult{
			Success: true,
			Meta: map[string]string{
				"file": path,
			},
		})
		return
	}

	log(INFO, "Created %s", path)
}

// ================= MAIN =================

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--json" {
			outputJSON = true
			break
		}
	}

	if len(os.Args) > 1 && os.Args[1] == "init" {
		runInit(os.Args[2:])
		return
	}

	envFileDefault := os.Getenv("FCM_ENV_FILE")
	loadDotEnv(envFileDefault)

	keyShort := flag.String("k", "", "")
	keyLong := flag.String("key", "", "")

	tokenShort := flag.String("t", "", "")
	tokenLong := flag.String("token", "", "")

	tokensFlag := flag.String("tokens", "", "")
	tokensFileFlag := flag.String("tokens-file", "", "")

	notifShort := flag.String("n", "", "")
	notifLong := flag.String("notification", "", "")

	dataShort := flag.String("d", "", "")
	dataLong := flag.String("data", "", "")

	topicFlag := flag.String("topic", "", "")

	conditionShort := flag.String("c", "", "")
	conditionLong := flag.String("condition", "", "")

	logShort := flag.String("l", "", "")
	logLong := flag.String("log", "", "")

	configShort := flag.String("f", "", "")
	configLong := flag.String("config", "", "")

	profileFlag := flag.String("profile", "", "")

	envFileFlag := flag.String("env-file", "", "")

	jsonFlag := flag.Bool("json", false, "")

	versionShort := flag.Bool("v", false, "")
	versionLong := flag.Bool("version", false, "")

	helpShort := flag.Bool("h", false, "")
	helpLong := flag.Bool("help", false, "")

	flag.Usage = func() {
		fmt.Println("Usage: fcm [options]")
		fmt.Println("Commands:")
		fmt.Println("  init                         Generate default fcm.yaml")
		fmt.Println("Options:")
		fmt.Println("  -k, --key <file>             Firebase key file")
		fmt.Println("  -t, --token <token>          Single FCM token")
		fmt.Println("  --tokens <t1,t2,t3>          Comma-separated token list")
		fmt.Println("  --tokens-file <file>         File with one token per line")
		fmt.Println("  -n, --notification <json>    Notification JSON")
		fmt.Println("  -d, --data <json>            Data JSON")
		fmt.Println("  -topic <topic>               Topic")
		fmt.Println("  -c, --condition <expr>       Condition")
		fmt.Println("  -f, --config <file>          YAML config file")
		fmt.Println("  --profile <name>             Profile inside config")
		fmt.Println("  --env-file <file>            Load additional .env file")
		fmt.Println("  -l, --log <level>            Log level: info|debug|json")
		fmt.Println("  --json                       Machine-readable output")
		fmt.Println("  -v, --version                Version")
		fmt.Println("  -h, --help                   Help")
	}

	flag.Parse()

	if *jsonFlag {
		outputJSON = true
	}

	if *envFileFlag != "" {
		_ = godotenv.Overload(*envFileFlag)
	}

	if *helpShort || *helpLong {
		flag.Usage()
		return
	}

	if *versionShort || *versionLong {
		if outputJSON {
			printJSON(CLIResult{
				Success: true,
				Meta: map[string]string{
					"version": version,
				},
			})
			return
		}
		fmt.Println(version)
		return
	}

	configPath := firstNonEmpty(*configShort, *configLong, os.Getenv("FCM_CONFIG"))
	var cfg *Config
	var err error
	if configPath != "" {
		cfg, err = loadConfig(configPath)
		if err != nil {
			if outputJSON {
				printJSON(CLIResult{
					Success: false,
					Error:   err.Error(),
				})
			} else {
				log(ERROR, "%v", err)
			}
			os.Exit(1)
		}
	}

	resolved, err := resolveConfig(cfg, *profileFlag)
	if err != nil {
		if outputJSON {
			printJSON(CLIResult{
				Success: false,
				Error:   err.Error(),
			})
		} else {
			log(ERROR, "%v", err)
		}
		os.Exit(1)
	}

	keyFile := firstNonEmpty(*keyShort, *keyLong, resolved.Key, os.Getenv("FCM_KEY"))
	fcmToken := firstNonEmpty(*tokenShort, *tokenLong, resolved.Token)

	csvTokens := parseTokensCSV(*tokensFlag)

	var fileTokens []string
	if *tokensFileFlag != "" {
		fileTokens, err = readTokensFile(*tokensFileFlag)
		if err != nil {
			if outputJSON {
				printJSON(CLIResult{
					Success: false,
					Error:   err.Error(),
				})
			} else {
				log(ERROR, "%v", err)
			}
			os.Exit(1)
		}
	}

	fcmTokens := firstNonEmptySlice(csvTokens, fileTokens, resolved.Tokens)

	topic := firstNonEmpty(*topicFlag, resolved.Topic)
	condition := firstNonEmpty(*conditionShort, *conditionLong, resolved.Condition)

	notifJSON := firstNonEmpty(*notifShort, *notifLong)
	var notif *Notification
	if resolved.Notification != nil {
		copyNotif := *resolved.Notification
		notif = &copyNotif
	}
	if notifJSON != "" {
		var parsed Notification
		if err := json.Unmarshal([]byte(notifJSON), &parsed); err != nil {
			if outputJSON {
				printJSON(CLIResult{
					Success: false,
					Error:   fmt.Sprintf("invalid notification JSON: %v", err),
				})
			} else {
				log(ERROR, "invalid notification JSON: %v", err)
			}
			os.Exit(1)
		}
		notif = &parsed
	}

	dataJSON := firstNonEmpty(*dataShort, *dataLong)
	data := cloneStringMap(resolved.Data)
	if dataJSON != "" {
		var parsed map[string]string
		if err := json.Unmarshal([]byte(dataJSON), &parsed); err != nil {
			if outputJSON {
				printJSON(CLIResult{
					Success: false,
					Error:   fmt.Sprintf("invalid data JSON: %v", err),
				})
			} else {
				log(ERROR, "invalid data JSON: %v", err)
			}
			os.Exit(1)
		}
		data = parsed
	}

	logMode := firstNonEmpty(*logShort, *logLong, resolved.Log, os.Getenv("FCM_LOG"), "info")
	switch logMode {
	case "debug":
		currentLogLevel = DEBUG
	case "json":
		jsonLogs = true
	case "info", "":
	default:
		if outputJSON {
			printJSON(CLIResult{
				Success: false,
				Error:   fmt.Sprintf("invalid log level %q; use info, debug, or json", logMode),
			})
		} else {
			log(ERROR, "invalid log level %q; use info, debug, or json", logMode)
		}
		os.Exit(1)
	}

	if keyFile == "" {
		msg := "missing key file; provide --key, config key, or FCM_KEY env"
		if outputJSON {
			printJSON(CLIResult{
				Success: false,
				Error:   msg,
			})
		} else {
			log(ERROR, "%s", msg)
		}
		os.Exit(1)
	}

	if err := validateTargets(fcmToken, fcmTokens, topic, condition); err != nil {
		if outputJSON {
			printJSON(CLIResult{
				Success: false,
				Error:   err.Error(),
			})
		} else {
			log(ERROR, "%v", err)
		}
		os.Exit(1)
	}

	if notif == nil && len(data) == 0 {
		msg := "message payload is empty; provide notification or data"
		if outputJSON {
			printJSON(CLIResult{
				Success: false,
				Error:   msg,
			})
		} else {
			log(ERROR, "%s", msg)
		}
		os.Exit(1)
	}

	ctx := context.Background()

	projectID, err := getProjectID(keyFile)
	if err != nil {
		if outputJSON {
			printJSON(CLIResult{
				Success: false,
				Error:   fmt.Sprintf("project id error: %v", err),
			})
		} else {
			log(ERROR, "project id error: %v", err)
		}
		os.Exit(1)
	}

	accessToken, err := getAccessToken(ctx, keyFile)
	if err != nil {
		if outputJSON {
			printJSON(CLIResult{
				Success: false,
				Error:   fmt.Sprintf("access token error: %v", err),
			})
		} else {
			log(ERROR, "access token error: %v", err)
		}
		os.Exit(1)
	}

	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", projectID)

	base := MessageBody{
		Topic:        topic,
		Condition:    condition,
		Notification: notif,
		Data:         data,
		Android:      cloneInterfaceMap(resolved.Android),
		Apns:         cloneInterfaceMap(resolved.Apns),
		Webpush:      cloneInterfaceMap(resolved.Webpush),
	}

	if len(fcmTokens) > 0 {
		result := sendMulticast(ctx, url, accessToken, base, fcmTokens)
		if outputJSON {
			printJSON(result)
		}
		if result.FailureCount > 0 {
			os.Exit(1)
		}
		return
	}

	if fcmToken != "" {
		base.Token = fcmToken
	}

	msg := FCMMessage{Message: base}

	messageID, code, err := sendWithRetry(ctx, url, accessToken, msg, 3)
	if err != nil {
		if outputJSON {
			printJSON(CLIResult{
				Success: false,
				Code:    code,
				Error:   err.Error(),
			})
		} else {
			log(ERROR, "Failed: %v", err)
		}
		os.Exit(1)
	}

	if outputJSON {
		printJSON(CLIResult{
			Success:   true,
			MessageID: messageID,
			Code:      code,
		})
		return
	}

	log(INFO, "Delivered")
	if messageID != "" {
		log(INFO, "Message ID: %s", messageID)
	}
	log(INFO, "Done")
}
