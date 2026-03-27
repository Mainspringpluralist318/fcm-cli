package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fcm/internal/log"
	"fcm/internal/model"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type fcmSuccessResponse struct {
	Name string `json:"name"`
}

func SendWithRetry(ctx context.Context, url, accessToken string, msg model.FCMMessage, retries int) (string, int, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return "", 0, fmt.Errorf("unable to marshal message: %w", err)
	}

	var lastCode int

	for i := 0; i <= retries; i++ {
		log.Log(model.DEBUG, "Sending request: %s", string(body))

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			return "", 0, fmt.Errorf("unable to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Log(model.ERROR, "Request error: %v", err)
			if i < retries {
				backoff := time.Duration(1<<i) * time.Second
				log.Log(model.DEBUG, "Retrying in %v...", backoff)
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
		log.Log(model.ERROR, "%v", err)

		if i < retries {
			backoff := time.Duration(1<<i) * time.Second
			log.Log(model.DEBUG, "Retrying in %v...", backoff)
			time.Sleep(backoff)
		} else {
			return "", lastCode, err
		}
	}

	return "", lastCode, fmt.Errorf("failed after retries")
}

func SendMulticast(ctx context.Context, url, accessToken string, base model.MessageBody, tokens []string) model.CLIResult {
	var wg sync.WaitGroup
	var completed int64
	var success int64

	results := make([]model.MulticastItem, len(tokens))
	total := int64(len(tokens))

	for i, t := range tokens {
		wg.Add(1)

		go func(idx int, tok string) {
			defer wg.Done()

			msg := model.FCMMessage{Message: base}
			msg.Message.Token = tok

			messageID, code, err := SendWithRetry(ctx, url, accessToken, msg, 3)
			item := model.MulticastItem{
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
				log.Log(model.ERROR, "Failed for token %s: %v", tok, err)
			}

			results[idx] = item

			done := atomic.AddInt64(&completed, 1)
			log.RenderProgress(done, total)
		}(i, t)
	}

	wg.Wait()
	if !log.OutputJSON {
		fmt.Println()
		log.Log(model.INFO, "Success: %d/%d", success, total)
	}

	return model.CLIResult{
		Success:      atomic.LoadInt64(&success) == total,
		SuccessCount: int(success),
		FailureCount: len(tokens) - int(success),
		Results:      results,
	}
}
