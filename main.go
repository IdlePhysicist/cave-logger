package main

import (
  "os"

  "github.com/sirupsen/logrus"

  "github.com/idlephysicist/service"
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

  _http  := os.Getenv("HTTP")
  if _http == "" {
    _http = ":8000"
  }

  service.NewService(log,  _http).Run()
}

