package config

import (
	"os"
	"log"
	"runtime"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/abdfnx/tran/dfs"
)

// TranConfig struct represents the config for the config.
type TranConfig struct {
	StartDir         string `mapstructure:"start_dir"`
	Borderless       bool   `mapstructure:"borderless"`
	Editor		     string `mapstructure:"editor"`
	EnableMouseWheel bool   `mapstructure:"enable_mousewheel"`
	ShowUpdates	     bool   `mapstructure:"show_updates"`
}

// Config represents the main config for the application.
type Config struct {
	Tran TranConfig `mapstructure:"config"`
}

func defualtEditor() string {
	if runtime.GOOS == "windows" {
		return "notepad.exe"
	}

	return "vim"
}

// LoadConfig loads a users config and creates the config if it does not exist
// located at `~/.config/tran.yml`
func LoadConfig(startDir *pflag.Flag) {
	var err error

	if runtime.GOOS != "windows" {
		homeDir, err := dfs.GetHomeDirectory()

		if err != nil {
			log.Fatal(err)
		}

		err = dfs.CreateDirectory(filepath.Join(homeDir, ".config", "tran"))
		if err != nil {
			log.Fatal(err)
		}

		viper.AddConfigPath("$HOME/.config/tran")
	} else {
		viper.AddConfigPath("$HOME")
	}

	viper.SetConfigName("tran")
	viper.SetConfigType("yml")

	// Setup config defaults.
	viper.SetDefault("config.start_dir", ".")
	viper.SetDefault("config.enable_mousewheel", true)
	viper.SetDefault("config.borderless", false)
	viper.SetDefault("config.editor", defualtEditor())
	viper.SetDefault("config.show_updates", true)

	if err := viper.SafeWriteConfig(); err != nil {
		if os.IsNotExist(err) {
			err = viper.WriteConfig()

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal(err)
		}
	}

	// Setup flags.
	err = viper.BindPFlag("start-dir", startDir)
	if err != nil {
		log.Fatal(err)
	}

	// Setup flag defaults.
	viper.SetDefault("start-dir", "")
}

// GetConfig returns the users config.
func GetConfig() (config Config) {
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error parsing config", err)
	}

	return
}
