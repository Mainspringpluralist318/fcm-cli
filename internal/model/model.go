package model

// CLIResult represents the final output of the CLI operation.
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

// MulticastItem represents the result for a single token in a multicast send.
type MulticastItem struct {
	Token     string `json:"token"`
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	Code      int    `json:"code,omitempty"`
	Error     string `json:"error,omitempty"`
}

// FCMMessage is the top-level structure for the FCM HTTP v1 API.
type FCMMessage struct {
	Message MessageBody `json:"message"`
}

// MessageBody contains the actual message content.
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

// Notification contains the visual notification content.
type Notification struct {
	Title string `json:"title" yaml:"title"`
	Body  string `json:"body" yaml:"body"`
}

// LogEntry represents a single log line in JSON format.
type LogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

// LogLevel type for logging.
type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	DEBUG LogLevel = "DEBUG"
)
