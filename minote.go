package main

import (
	"log"

	"os"
	"path/filepath"
	"runtime"

	"github.com/marcs100/minote/config"
	"github.com/marcs100/minote/minotedb"
	"github.com/marcs100/minote/ui"
)

const VERSION = "0.007"

func main() {
	var err error
	var dir_err error
	var appConfig *config.Config
	const confFileName = "config.toml"
	var confFilePath string
	var homeDir string

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

	appConfig, err = config.GetConfig(confFile)
	if err != nil {
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

	ui.StartUI(appConfig, confFile, VERSION)
}

func CreateAppConfig(homeDir string) config.Config {
	appSettings := config.AppSettings{
		Database: filepath.Join(homeDir, "minoteData", "minote.db"), //this one for release
		//Database: filepath.Join(homeDir,"sync","scribe","scribeNB.db"), //temp one for dev
		InitialLayout:    "grid",
		InitialView:      "pinned",
		NoteHeight:       350,
		NoteWidth:        500,
		RecentNotesLimit: 50,
		GridMaxPages:     500,
		FontSize:         12,
		ThemeVariant:     "system",
		DarkColourNote:   "#242424",
		LightColourNote:  "#e2e2e2",
		DarkColourBg:     "#1e1e1e",
		LightColourBg:    "#efefef",
		DarkColourCtBg:   "#1e1e1e",
		LightColourCtBg:  "#efefef",
	}
	newConfig := config.Config{
		Title:    "Minote config",
		Settings: appSettings,
	}

	return newConfig
}
