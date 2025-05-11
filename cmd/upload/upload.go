/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package upload

import (
	"fmt"
	"time"

	"github.com/jolpica/jf1/cmd/input"
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
		Args: cobra.MinimumNArgs(1),
	}

	cmd.Flags().Bool("dry-run", false, "run command in dry-run mode")
	viper.BindPFlag("upload.dry-run", cmd.Flags().Lookup("dry-run"))

	cmd.Flags().StringP("scanned-file", "s", "scanned.gob", "file name to save previously checked directories")
	viper.BindPFlag("upload.scanned-file", cmd.Flags().Lookup("scanned-file"))

	cmd.Flags().Bool("only-update-scanned", false, "only update the contexts of scanned-file, do not make any requests")
	viper.BindPFlag("upload.only-update-scanned", cmd.Flags().Lookup("only-update-scanned"))

	cmd.Flags().IntP("max-concurrent-requests", "m", 3, "maximum number of requests to jolpica-f1 at once")
	viper.BindPFlag("upload.max-concurrent-requests", cmd.Flags().Lookup("max-concurrent-requests"))

	return cmd
}

func runUploadCmd(cmd *cobra.Command, args []string) error {
	start := time.Now()
	fmt.Printf("Scanning Dirs: %v\n", args)
	err := uploader.RunUploader(args, input.I.Upload, input.I.Secret.Token)

	if input.I.Upload.Verbose {
		fmt.Printf("End of program.\nerr: %v\nTook: %v\n", err, time.Since(start))
	}
	return err
}
