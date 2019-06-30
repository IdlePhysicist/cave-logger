package main

import (
  "os"

  "github.com/sirupsen/logrus"

  "github.com/idlephysicist/cave-logger/pkg/service"
)

func main() {
  log := logrus.New()
  log.SetFormatter(&logrus.TextFormatter{})

  _level := os.Getenv("LOG_LEVEL")
  if _level == "" {
    _level = "debug"
  }
  level, err := logrus.ParseLevel(_level)
  if err != nil {
    level = logrus.InfoLevel
  }
  log.SetLevel(level)

  port  := os.Getenv("HTTP")
  if _http == "" {
    port = ":8000"
  }

  service.NewService(log,  port).Run()
}
