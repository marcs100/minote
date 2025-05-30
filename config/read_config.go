package config

import(
	"fmt"
	"github.com/BurntSushi/toml"
)

func GetConfig(toml_file string)(*Config, error){
	var config Config
	if _, err := toml.DecodeFile(toml_file, &config); err != nil {
		fmt.Println(err)
		return nil,err
	}
	//fmt.Printf("Title: %s\n", config.Title)
	//fmt.Printf("database file: %s\n", config.AppSettings.Database)
	//fmt.Printf("recent notes limit: %d\n", config.AppSettings.RecentNotesLimit)
	return &config,nil
}



