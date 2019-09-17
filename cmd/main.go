package main

import (
  "os"

  "github.com/sirupsen/logrus"

  "github.com/idlephysicist/cave-logger/internal/pkg/keeper"
  "github.com/idlephysicist/cave-logger/internal/tui"
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

  dbFileName := `xyz.db`

  keeper.New(log, dbFileName)

  tui := tui.New()

  if err := tui.Start(); err != nil {
    log.Fatalf("main: Cannot start tui: %s", err)
  }
}
