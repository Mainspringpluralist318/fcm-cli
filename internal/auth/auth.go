package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
)

const MessagingScope = "https://www.googleapis.com/auth/firebase.messaging"

func GetAccessToken(ctx context.Context, filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("unable to read key file: %w", err)
	}

	config, err := google.JWTConfigFromJSON(data, MessagingScope)
	if err != nil {
		return "", fmt.Errorf("unable to parse JWT config: %w", err)
	}

	token, err := config.TokenSource(ctx).Token()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve token: %w", err)
	}

	return token.AccessToken, nil
}

func GetProjectID(filename string) (string, error) {
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
