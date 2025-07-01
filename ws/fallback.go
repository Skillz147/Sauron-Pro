package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"o365/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitWSRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func EnqueueWSFallback(userID string, msgType string, payload interface{}) {
	data := SocketMessage{Type: msgType, Data: payload}
	b, err := json.Marshal(data)
	if err != nil {
		utils.SystemLogger.Error().Err(err).Msg("❌ Failed to marshal WS fallback payload")
		return
	}
	key := fmt.Sprintf("ws_queue:%s", userID)
	redisClient.RPush(context.Background(), key, b)
	redisClient.Expire(context.Background(), key, 24*time.Hour)
}

func DeliverPendingMessages(userID string) {
	key := fmt.Sprintf("ws_queue:%s", userID)
	ctx := context.Background()

	for {
		raw, err := redisClient.LPop(ctx, key).Result()
		if err != nil || raw == "" {
			break
		}
		var msg SocketMessage
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			continue
		}
		SendToUser(userID, msg.Type, msg.Data)
	}
}
