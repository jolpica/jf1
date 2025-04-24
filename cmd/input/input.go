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
	Upload uploader.UploadConfig
	Secret Jf1Secret
}

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName("jf1.toml")

	viper.SetEnvPrefix("JF1")
	viper.BindEnv("secret.token")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	err := viper.UnmarshalExact(&I)
	cobra.CheckErr(err)

	fmt.Printf("Config: %+v\n", I)
}
