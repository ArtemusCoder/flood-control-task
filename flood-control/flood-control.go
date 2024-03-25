package flood_control

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type FloodControl struct {
	redisClient *redis.Client
	N           int
	K           int
	mu          sync.Mutex
}

func NewFloodControl(redisAddr string, N, K int) *FloodControl {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &FloodControl{
		redisClient: client,
		N:           N,
		K:           K,
	}
}

func (fc *FloodControl) Check(ctx context.Context, userID int64) (bool, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	key := fmt.Sprintf("user:%d", userID)
	currentTime := time.Now()
	err := fc.redisClient.ZAdd(ctx, key, redis.Z{Score: float64(currentTime.Unix()), Member: currentTime.String()}).Err()
	if err != nil {
		return false, err
	}

	minScore := float64(time.Now().Add(-time.Duration(fc.N) * time.Second).Unix())
	_, err = fc.redisClient.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%f", minScore)).Result()
	if err != nil {
		return false, err
	}

	count, err := fc.redisClient.ZCard(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if int(count) > fc.K {
		return false, nil
	}

	return true, nil
}
