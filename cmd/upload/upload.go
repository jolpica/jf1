/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package upload

import (
	"fmt"
	"time"

	"github.com/jolpica/jf1/pkg/uploader"
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
		RunE: runUploadCmd,
		Args: cobra.MaximumNArgs(1),
	}

	cmd.Flags().Bool("dry-run", false, "run command in dry-run mode")
	viper.BindPFlag("upload.dry-run", cmd.Flags().Lookup("dry-run"))

	cmd.Flags().StringP("scanned-file", "s", "scanned.gob", "file name to save previously checked directories")
	viper.BindPFlag("upload.scanned-file", cmd.Flags().Lookup("scanned-file"))

	return cmd
}

func runUploadCmd(cmd *cobra.Command, args []string) error {
	start := time.Now()
	dirsPath := "."
	if len(args) >= 1 {
		dirsPath = args[0]
	}
	fmt.Printf("Scanning Dir: %v\n", dirsPath)
	err := uploader.RunUploader(dirsPath, viper.GetString("upload.scanned-file"), viper.GetBool("upload.dry-run"))

	fmt.Printf("End of program. err: %v\nTook: %v\n", err, time.Since(start))
	return err
}
