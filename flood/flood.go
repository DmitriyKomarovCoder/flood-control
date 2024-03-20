package flood

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type FloodContr struct {
	maxRequest int
	maxTime    time.Duration
	db         *redis.Client
}

func NewFloodControl(maxReq int, maxTime time.Duration, db *redis.Client) *FloodContr {
	return &FloodContr{maxRequest: maxReq, maxTime: maxTime, db: db}
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

func (fc *FloodContr) Check(ctx context.Context, userID int64) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		key := fmt.Sprintf("%d:%d", userID, time.Now().Unix())
		err := fc.db.Set(ctx, key, ".", 5*time.Second).Err()
		if err != nil {
			return false, err
		}

		var cursor uint64 = 0
		var count int64 = 0

		pattern := fmt.Sprintf("%d:*", userID)

		for {
			keys, nextCursor, err := fc.db.Scan(ctx, cursor, pattern, int64(fc.maxRequest)+1).Result()
			if err != nil {
				return false, err
			}

			count += int64(len(keys))

			if nextCursor == 0 {
				break
			}

			cursor = nextCursor
		}

		if count > int64(fc.maxRequest) {
			return false, nil
		}
	}
	return true, nil
}
