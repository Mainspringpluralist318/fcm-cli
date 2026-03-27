package main

import (
	"context"
	"encoding/json"
	"fcm/internal/auth"
	"fcm/internal/config"
	"fcm/internal/fcm"
	"fcm/internal/log"
	"fcm/internal/model"
	"fcm/internal/util"
	"flag"
	"fmt"
	"os"
)

var version = "dev"

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

	path := util.FirstNonEmpty(*fileLong, *fileShort)
	if path == "" {
		path = "fcm.yaml"
	}

	if err := config.WriteDefaultConfig(path, *force); err != nil {
		if log.OutputJSON {
			log.PrintJSON(model.CLIResult{
				Success: false,
				Error:   err.Error(),
			})
		} else {
			log.Log(model.ERROR, "%v", err)
		}
		os.Exit(1)
	}

	if log.OutputJSON {
		log.PrintJSON(model.CLIResult{
			Success: true,
			Meta: map[string]string{
				"file": path,
			},
		})
		return
	}

	log.Log(model.INFO, "Created %s", path)
}

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--json" {
			log.OutputJSON = true
			break
		}
	}

	if len(os.Args) > 1 && os.Args[1] == "init" {
		runInit(os.Args[2:])
		return
	}

	envFileDefault := os.Getenv("FCM_ENV_FILE")
	config.LoadDotEnv(envFileDefault)

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
		log.OutputJSON = true
	}

	if *envFileFlag != "" {
		config.LoadDotEnv(*envFileFlag)
	}

	if *helpShort || *helpLong {
		flag.Usage()
		return
	}

	if *versionShort || *versionLong {
		if log.OutputJSON {
			log.PrintJSON(model.CLIResult{
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

	configPath := util.FirstNonEmpty(*configShort, *configLong, os.Getenv("FCM_CONFIG"))
	var cfg *config.Config
	var err error
	if configPath != "" {
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			if log.OutputJSON {
				log.PrintJSON(model.CLIResult{
					Success: false,
					Error:   err.Error(),
				})
			} else {
				log.Log(model.ERROR, "%v", err)
			}
			os.Exit(1)
		}
	}

	resolved, err := config.ResolveConfig(cfg, *profileFlag)
	if err != nil {
		if log.OutputJSON {
			log.PrintJSON(model.CLIResult{
				Success: false,
				Error:   err.Error(),
			})
		} else {
			log.Log(model.ERROR, "%v", err)
		}
		os.Exit(1)
	}

	keyFile := util.FirstNonEmpty(*keyShort, *keyLong, resolved.Key, os.Getenv("FCM_KEY"))
	fcmToken := util.FirstNonEmpty(*tokenShort, *tokenLong, resolved.Token)

	csvTokens := util.ParseTokensCSV(*tokensFlag)

	var fileTokens []string
	if *tokensFileFlag != "" {
		fileTokens, err = util.ReadTokensFile(*tokensFileFlag)
		if err != nil {
			if log.OutputJSON {
				log.PrintJSON(model.CLIResult{
					Success: false,
					Error:   err.Error(),
				})
			} else {
				log.Log(model.ERROR, "%v", err)
			}
			os.Exit(1)
		}
	}

	fcmTokens := util.FirstNonEmptySlice(csvTokens, fileTokens, resolved.Tokens)

	topic := util.FirstNonEmpty(*topicFlag, resolved.Topic)
	condition := util.FirstNonEmpty(*conditionShort, *conditionLong, resolved.Condition)

	notifJSON := util.FirstNonEmpty(*notifShort, *notifLong)
	var notif *model.Notification
	if resolved.Notification != nil {
		copyNotif := *resolved.Notification
		notif = &copyNotif
	}
	if notifJSON != "" {
		var parsed model.Notification
		if err := json.Unmarshal([]byte(notifJSON), &parsed); err != nil {
			if log.OutputJSON {
				log.PrintJSON(model.CLIResult{
					Success: false,
					Error:   "invalid notification JSON",
					Meta:    map[string]string{"raw": notifJSON},
				})
			} else {
				log.Log(model.ERROR, "Invalid notification JSON: %v", err)
			}
			os.Exit(1)
		}
		notif = &parsed
	}

	dataJSON := util.FirstNonEmpty(*dataShort, *dataLong)
	data := resolved.Data
	if dataJSON != "" {
		if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
			if log.OutputJSON {
				log.PrintJSON(model.CLIResult{
					Success: false,
					Error:   "invalid data JSON",
					Meta:    map[string]string{"raw": dataJSON},
				})
			} else {
				log.Log(model.ERROR, "Invalid data JSON: %v", err)
			}
			os.Exit(1)
		}
	}

	logLevelStr := util.FirstNonEmpty(*logShort, *logLong, resolved.Log)
	switch logLevelStr {
	case "debug":
		log.CurrentLogLevel = model.DEBUG
	case "json":
		log.JSONLogs = true
	}

	if keyFile == "" {
		err := fmt.Errorf("Firebase key file is required (-k or FCM_KEY)")
		if log.OutputJSON {
			log.PrintJSON(model.CLIResult{
				Success: false,
				Error:   err.Error(),
			})
		} else {
			log.Log(model.ERROR, "%v", err)
		}
		os.Exit(1)
	}

	ctx := context.Background()
	accessToken, err := auth.GetAccessToken(ctx, keyFile)
	if err != nil {
		if log.OutputJSON {
			log.PrintJSON(model.CLIResult{
				Success: false,
				Error:   err.Error(),
			})
		} else {
			log.Log(model.ERROR, "Auth error: %v", err)
		}
		os.Exit(1)
	}

	projectID, err := auth.GetProjectID(keyFile)
	if err != nil {
		if log.OutputJSON {
			log.PrintJSON(model.CLIResult{
				Success: false,
				Error:   err.Error(),
			})
		} else {
			log.Log(model.ERROR, "Project ID error: %v", err)
		}
		os.Exit(1)
	}

	fcmURL := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", projectID)

	baseMsg := model.MessageBody{
		Notification: notif,
		Data:         data,
		Android:      resolved.Android,
		Apns:         resolved.Apns,
		Webpush:      resolved.Webpush,
	}

	var finalResult model.CLIResult

	if len(fcmTokens) > 0 {
		finalResult = fcm.SendMulticast(ctx, fcmURL, accessToken, baseMsg, fcmTokens)
	} else {
		// Single target: token, topic, or condition
		msg := model.FCMMessage{Message: baseMsg}
		msg.Message.Token = fcmToken
		msg.Message.Topic = topic
		msg.Message.Condition = condition

		// Validate that exactly one target is provided
		count := 0
		if fcmToken != "" {
			count++
		}
		if topic != "" {
			count++
		}
		if condition != "" {
			count++
		}
		if count != 1 {
			err := fmt.Errorf("provide exactly one target: token, tokens, topic or condition")
			if log.OutputJSON {
				log.PrintJSON(model.CLIResult{
					Success: false,
					Error:   err.Error(),
				})
			} else {
				log.Log(model.ERROR, "%v", err)
			}
			os.Exit(1)
		}

		messageID, code, err := fcm.SendWithRetry(ctx, fcmURL, accessToken, msg, 3)
		if err != nil {
			finalResult = model.CLIResult{
				Success: false,
				Code:    code,
				Error:   err.Error(),
			}
		} else {
			finalResult = model.CLIResult{
				Success:   true,
				MessageID: messageID,
			}
			log.Log(model.INFO, "Success! Message ID: %s", messageID)
		}
	}

	if log.OutputJSON {
		log.PrintJSON(finalResult)
	}

	if !finalResult.Success {
		os.Exit(1)
	}
}
