/*
 * @Author: Tomato
 * @Date: 2026-05-10 21:29:53
 */
package main

import (
	"fmt"

	"github.com/compassa/tomatomq/internal/logger"
)

func main() {
	defer logger.Sync()

	if err := logger.InitBroker(); err != nil {
		panic(err)
	}

	fmt.Printf("broker")
}
