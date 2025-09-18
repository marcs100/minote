package main

import (
	"log"

	"os"
	"path/filepath"
	"runtime"

	"github.com/marcs100/minote/config"
	"github.com/marcs100/minote/main_app"
	"github.com/marcs100/minote/minotedb"
	"github.com/marcs100/minote/ui"
)

const VERSION = "0.010"

func main() {
	var err error
	var dir_err error
	var appConfig *config.Config
	const confFileName = "config.toml"
	var confFilePath string
	var homeDir string

	var about = main_app.About{
		Version:     VERSION,
		Licence:     "MIT",
		LicenceLink: "https://mit-license.org/",
		Maintainer:  "marcs100@gmail.com",
		Website:     "https://github.com/marcs100/minote",
	}

	if homeDir, dir_err = os.UserHomeDir(); dir_err != nil {
		log.Panicln(dir_err)
	}

	if runtime.GOOS == "windows" {
		//This needs to be improved but will do for now!!!!
		confFilePath = filepath.Join(homeDir, "minote")
	} else {
		confFilePath = filepath.Join(homeDir, ".config/minote")
	}

	confFile := filepath.Join(confFilePath, confFileName)

	if _, f_err := os.Stat(confFile); f_err != nil {
		//write the default config confFile
		if err = os.MkdirAll(confFilePath, os.ModePerm); err != nil {
			log.Panicln("something went wrong with conf file path!!")
		}

		//create the default config.toml
		newConfig := CreateAppConfig(homeDir)

		if err = config.WriteConfig(confFile, newConfig); err != nil {
			log.Panicf("Error writing config file: %s", err)
		}
	}

	if appConfig, err = config.GetConfig(confFile); err != nil {
		log.Panicln(err)
		return
	}

	//check if the database already exists
	if _, dbf_err := os.Stat(appConfig.Settings.Database); dbf_err != nil {
		dbName := filepath.Base(appConfig.Settings.Database)
		dbPath := filepath.Dir(appConfig.Settings.Database)

		//create a new database
		if err = minotedb.CreateNew(dbName, dbPath); err != nil {
			log.Panicf("Something went wrong creating new db: %s", err)
		}
		minotedb.Close()
	}

	err = minotedb.Open(appConfig.Settings.Database)
	defer minotedb.Close()
	if err != nil {
		log.Panicln(err)
	}

	//validate config
	if err = main_app.ValidateConfig(appConfig); err != nil {
		log.Panicln(err)
		return
	}
	ui.StartUI(appConfig, confFile, about)
}

func CreateAppConfig(homeDir string) config.Config {
	dbFile := ""
	dbDir := ""
	var err error
	if runtime.GOOS != "windows" {
		dbFile = filepath.Join(homeDir, ".minote", "minote.db")
		dbDir = filepath.Join(homeDir, ".minote")
	} else {
		dbFile = filepath.Join(homeDir, "MinoteData", "minote.db")
		dbDir = filepath.Join(homeDir, "MinoteData")
	}
	if _, f_err := os.Stat(dbDir); f_err != nil {
		//create config path
		if err = os.MkdirAll(dbDir, os.ModePerm); err != nil {
			log.Panicf("Something went wrong creating config path: %s", err)
		}
	}
	appSettings := config.AppSettings{
		Database:          dbFile,
		InitialLayout:     "grid",
		InitialView:       "pinned",
		NoteHeight:        350,
		NoteWidth:         500,
		RecentNotesLimit:  50,
		GridMaxPages:      500,
		FontSize:          12,
		DateFormat:        "02-01-2006",
		TimeFormat:        "15:04",
		DateTimeFormat:    "[02-01-2006 at 15:04]",
		ThemeVariant:      "auto",
		DarkColourNote:    "#242424",
		LightColourNote:   "#e2e2e2",
		DarkColourBg:      "#1e1e1e",
		LightColourBg:     "#efefef",
		DarkColourFg:      "#e7e7e7",
		LightColourFg:     "#666666",
		DarkColourCtBg:    "#1e1e1e",
		LightColourCtBg:   "#efefef",
		DarkColourAccent:  "#d6701c",
		LightColourAccent: "#2684ff",
		DarkColourButton:  "#2f2f2f",
		LightColourButton: "#c3c3c3",
	}
	newConfig := config.Config{
		Title:    "Minote config",
		Settings: appSettings,
	}

	return newConfig
}
