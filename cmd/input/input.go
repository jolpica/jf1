package input

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var I Jf1Input

type Jf1Input struct {
	Upload UploadInput

	Token string
}

type UploadInput struct {
	BaseUrl     string `mapstructure:"base-url"`
	DryRun      bool   `mapstructure:"dry-run"`
	ScannedFile string `mapstructure:"scanned-file"`
}

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName("jf1.toml")

	viper.SetEnvPrefix("JF1")
	viper.BindEnv("token")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	err := viper.UnmarshalExact(&I)
	cobra.CheckErr(err)

	fmt.Printf("Config: %+v\n", I)
}
