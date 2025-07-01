// telegram_worker.go
package telegramq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func StartTelegramWorker() {
	ticker := time.NewTicker(300 * time.Millisecond)
	go func() {
		for range ticker.C {
			if err := processNext(); err != nil {
				log.Warn().Err(err).Msg("❌ Telegram queue failed")
			}
		}
	}()
}

func processNext() error {
	ctx := context.Background()
	raw, err := Redis.LPop(ctx, queueName).Result()
	if err != nil || raw == "" {
		return nil // nothing to process
	}

	var job TelegramJob
	if err := json.Unmarshal([]byte(raw), &job); err != nil {
		return fmt.Errorf("invalid job format: %w", err)
	}

	if job.File.Content == "" {
		return sendTelegramText(job)
	}
	return sendTelegramDocument(job)
}

func sendTelegramText(job TelegramJob) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", job.BotToken)

	payload := map[string]string{
		"chat_id":    job.ChatID,
		"text":       job.Text,
		"parse_mode": "HTML",
	}
	data, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Info().Str("chat_id", job.ChatID).Str("slug", job.Slug).Msg("📤 Telegram text sent")
	} else {
		log.Warn().Str("chat_id", job.ChatID).Int("status", resp.StatusCode).Msg("⚠️ Telegram rejected (text)")
	}
	return nil
}

func sendTelegramDocument(job TelegramJob) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", job.BotToken)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("document", job.File.Filename)
	if err != nil {
		return err
	}
	if _, err := part.Write([]byte(job.File.Content)); err != nil {
		return err
	}

	_ = writer.WriteField("chat_id", job.ChatID)
	_ = writer.WriteField("caption", job.Text)
	_ = writer.WriteField("parse_mode", "HTML")
	writer.Close()

	req, _ := http.NewRequest("POST", url, &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Info().Str("chat_id", job.ChatID).Str("slug", job.Slug).Msg("📤 Telegram document sent")
	} else {
		log.Warn().Str("chat_id", job.ChatID).Int("status", resp.StatusCode).Msg("⚠️ Telegram rejected (doc)")
	}
	return nil
}
