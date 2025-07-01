// telegram_queue.go
package telegramq

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

const queueName = "telegram_queue"

func InitTelegramRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // change if needed
		Password: "",
		DB:       0,
	})
}

type TelegramJob struct {
	Slug     string      `json:"slug"`
	BotToken string      `json:"bot_token"`
	ChatID   string      `json:"chat_id"`
	Text     string      `json:"text"` // required caption
	File     FilePayload `json:"file"` // required file
}

type FilePayload struct {
	Filename string `json:"filename"`
	Content  string `json:"content"` // plain text cookie string
}

func EnqueueTelegramMessage(job TelegramJob) error {
	ctx := context.Background()
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return Redis.RPush(ctx, queueName, data).Err()
}
