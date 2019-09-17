package main

import (
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os"

  "github.com/sirupsen/logrus"

  "github.com/idlephysicist/cave-logger/internal/pkg/db"
  "github.com/idlephysicist/cave-logger/internal/gui"
  "github.com/idlephysicist/cave-logger/internal/pkg/model"
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

  cfg := func (_yamlFile string) *model.Config {
    var _cfg model.Config
    
    if _yamlFile == `` {
      _yamlFile = `./config/config.yml`
    }
		
		yamlFile, err := ioutil.ReadFile(_yamlFile)
		if err != nil {
			log.Fatalf("main.readfile: %v", err)
		}

		err = yaml.Unmarshal(yamlFile, &_cfg)
		if err != nil {
			log.Fatalf("main.unmarshalYAML: %v", err)
		}

		return &_cfg
	}(``)


  db := db.New(log, cfg.Database.Filename)

  gui := gui.New(db, log)

  if err := gui.Start(); err != nil {
    log.Fatalf("main: Cannot start tui: %s", err)
  }
}
