package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/IdlePhysicist/cave-logger/internal/db"
	"github.com/IdlePhysicist/cave-logger/internal/gui"
	"github.com/IdlePhysicist/cave-logger/internal/model"
)

var commit, version, date string

func main() {
	// Parse cfg override
	var (
		cfgOverride string
		versionCall bool
	)
	flag.StringVarP(&cfgOverride, `config`, `c`, ``, `Config file override`)
	flag.BoolVarP(&versionCall, `version`, `v`, false, `Print version info`)
	flag.Parse()

	if versionCall {
		fmt.Printf("cave-logger %s (commit: %s) (built: %s)\n", version, commit, date)
		os.Exit(0)
	}

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
	cfg := func(_file string) *model.Config {
		var _cfg model.Config

		if _file == `` {
			_file = fmt.Sprintf("%s/.config/cave-logger/config.json", os.Getenv("HOME"))
		}

		file, err := ioutil.ReadFile(_file)
		if err != nil {
			log.Fatalf("main.readfile: %v", err)
		}

		err = json.Unmarshal(file, &_cfg)
		if err != nil {
			log.Fatalf("main.unmarshal: %v", err)
		}

		return &_cfg
	}(cfgOverride)

	cfg.Database.Filename = strings.Join(
		[]string{os.Getenv("HOME"), cfg.Database.Filename},
		"/",
	)

	// Initialise the database connection and handler
	db := db.New(log, cfg.Database.Filename)

	// Initialise the Gui / Tui
	gui := gui.New(db)
	gui.ProcessColors(cfg.Colors)

	if err := gui.Start(); err != nil {
		log.Fatalf("main: Cannot start tui: %s", err)
	}
}
