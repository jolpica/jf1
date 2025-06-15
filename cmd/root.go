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
	"os"

	"github.com/jolpica/jf1/cmd/input"
	"github.com/jolpica/jf1/cmd/upload"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:          "jf1",
	Short:        "CLI tools for the jolpica-f1 project",
	SilenceUsage: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&input.ConfigFile, "config", "", "override configuration file to use")

	cobra.OnInitialize(input.InitConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "increase the verbosity of output")
	viper.BindPFlag("upload.verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().StringP("base-url", "u", "https://localhost:8000", "base url for jolpica-f1 api requests")
	viper.BindPFlag("upload.base-url", rootCmd.PersistentFlags().Lookup("base-url"))

	rootCmd.AddCommand(upload.NewUploadCmd())
	rootCmd.AddCommand(stressCmd)
}
