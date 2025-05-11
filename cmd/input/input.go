package input

import (
	"fmt"
	"os"

	"github.com/jolpica/jf1/pkg/uploader"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var I Jf1Input

type Jf1Secret struct {
	Token string
}

type Jf1Input struct {
	Verbose bool `mapstructure:"verbose"`
	Test    string

	Upload uploader.UploadConfig
	Secret Jf1Secret
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName("jf1.toml")

	viper.SetEnvPrefix("JF1")
	viper.AutomaticEnv()
	viper.BindEnv("secret.token", "JF1_SECRET_TOKEN")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Prevent secret.token from being set in config
	if viper.InConfig("secret.token") {
		fmt.Fprintln(os.Stderr, "Error: secret.token must not be set in config file. Use env var JF1_SECRET_TOKEN.")
		os.Exit(1)
	}

	err := viper.UnmarshalExact(&I)
	cobra.CheckErr(err)

	if I.Verbose {
		masked := I
		masked.Secret = Jf1Secret{}
		fmt.Printf("Config: %+v\n", masked)
	}
}
