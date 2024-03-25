package main

import (
	"context"
	"fmt"
	"log"
	floodControl "task/flood-control"
	"time"
)

func main() {
	redisAddr := "localhost:6379"
	N := 10
	K := 5
	control := floodControl.NewFloodControl(redisAddr, N, K)

	userID := int64(123)

	for i := 0; i < 7; i++ {
		result, err := control.Check(context.Background(), userID)
		if err != nil {
			log.Fatal(err)
		}
		if result {
			fmt.Println("Flood control check passed")
		} else {
			fmt.Println("Flood control check failed")
		}
		time.Sleep(1 * time.Second)
	}

	time.Sleep(5 * time.Second)

	for i := 0; i < 5; i++ {
		result, err := control.Check(context.Background(), userID)
		if err != nil {
			log.Fatal(err)
		}
		if result {
			fmt.Println("Flood control check passed")
		} else {
			fmt.Println("Flood control check failed")
		}
		time.Sleep(1 * time.Second)
	}

}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
