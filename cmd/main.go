package main

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os"

  "github.com/sirupsen/logrus"
  flag "github.com/spf13/pflag"

  "github.com/idlephysicist/cave-logger/internal/db"
  "github.com/idlephysicist/cave-logger/internal/gui"
  "github.com/idlephysicist/cave-logger/internal/model"
)

func main() {
  // Parse cfg override
  var cfgOverride string
  flag.StringVarP(&cfgOverride, `config`, `c`, ``, `Config file override`)

  // Set up logger
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


  // Read config file
  cfg := func (_yamlFile string) *model.Config {
    var _cfg model.Config

    if _yamlFile == `` {
      _yamlFile = fmt.Sprintf("%s/.config/cave-logger/config.yml", os.Getenv("HOME"))
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
	}(cfgOverride)

  // Initialise the database connection and handler
  db := db.New(log, cfg.Database.Filename)

  // Initialise the Gui / Tui
  gui := gui.New(db)

  if err := gui.Start(); err != nil {
    log.Fatalf("main: Cannot start tui: %s", err)
  }
}
