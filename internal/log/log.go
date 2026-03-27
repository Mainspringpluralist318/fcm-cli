package log

import (
	"encoding/json"
	"fcm/internal/model"
	"fmt"
	"os"
	"time"
)

var (
	CurrentLogLevel = model.INFO
	JSONLogs        = false
	OutputJSON      = false
)

func Log(level model.LogLevel, msg string, args ...interface{}) {
	if OutputJSON {
		return
	}
	if level == model.DEBUG && CurrentLogLevel != model.DEBUG {
		return
	}

	message := fmt.Sprintf(msg, args...)

	if JSONLogs {
		entry := model.LogEntry{
			Level:   string(level),
			Message: message,
			Time:    time.Now().Format(time.RFC3339),
		}
		_ = json.NewEncoder(os.Stdout).Encode(entry)
		return
	}

	color := ""
	reset := "\033[0m"

	switch level {
	case model.INFO:
		color = "\033[34m"
	case model.ERROR:
		color = "\033[31m"
	case model.DEBUG:
		color = "\033[33m"
	}

	fmt.Printf("%s[%s]%s %s\n", color, level, reset, message)
}

func PrintJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func RenderProgress(done, total int64) {
	if OutputJSON {
		return
	}
	percent := float64(done) / float64(total) * 100
	fmt.Printf("\rProgress: %d/%d (%.0f%%)", done, total, percent)
}
