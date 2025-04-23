/*
Copyright © 2025 Jessica Perry

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jolpica/jf1/cmd/upload"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Input Jf1Input
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jf1",
	Short: "CLI tools for the jolpica-f1 project",
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./jf1.toml)")

	rootCmd.AddCommand(upload.NewUploadCmd())
	rootCmd.AddCommand(stressCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("toml")
		viper.SetConfigName("jf1.toml")
	}

	viper.SetEnvPrefix("JF1")
	viper.AutomaticEnv()
	viper.BindEnv("token")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	err := viper.UnmarshalExact(&Input)
	cobra.CheckErr(err)

	fmt.Printf("Config: %+v", Input)
}

type Jf1Input struct {
	Upload UploadInput

	Token string
}

type UploadInput struct {
	DryRun bool `mapstructure:"dry-run"`
}
