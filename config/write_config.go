package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

func WriteConfig(tomlFile string, newConfig Config) error {
	var f *os.File
	var f_err error
	if f, f_err = os.Create(tomlFile); f_err != nil {
		return f_err
	}
	defer f.Close()
	err := toml.NewEncoder(f).Encode(newConfig)
	return err
}
