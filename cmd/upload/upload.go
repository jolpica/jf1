/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package upload

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewUploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("upload called")
			fmt.Printf("Value: %v\n", viper.GetString("token"))
			fmt.Printf("Speed: %v\n", viper.GetString("upload.speed"))
			fmt.Printf("dry run: %v\n", viper.GetBool("upload.dry-run"))
		},
	}

	cmd.Flags().Bool("dry-run", false, "run command in dry-run mode")
	viper.BindPFlag("upload.dry-run", cmd.Flags().Lookup("dry-run"))

	return cmd
}
