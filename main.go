package main

import (
    "github.com/kintohub/common-go/logger"
)

func main() {
    log := logger.NewSimpleLogger()
    logger.SetLogger(log)
    logger.Debug("test")
}
