/*
 * @Author: Tomato
 * @Date: 2026-05-10 21:29:53
 */
package main

import (
	"fmt"
	"log/slog"

	"github.com/compassa/tomatomq/internal/broker"
)

func main() {
	env, cfg := broker.LoadConfig()

	broker.LoadLogger(env, &cfg.Log)

	broker.AppLogger.Info("info", slog.Int("key", 124))
	broker.AppLogger.Error("error", slog.String("key", "value"))

	fmt.Printf("broker")
}
