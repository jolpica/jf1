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
		Short: "Upload new data to the jolpica-f1 api",
		Long: `Scans the given directories for new data and uploads new data to the jolpica-f1 api.
		
		Previously uploaded directories are saved in the file specified by --uploaded-file.`,
		RunE: runUploadCmd,
		Args: cobra.MinimumNArgs(1),
	}

	cmd.Flags().Bool("dry-run", false, "run command in dry-run mode")
	viper.BindPFlag("upload.dry-run", cmd.Flags().Lookup("dry-run"))

	cmd.Flags().StringP("uploaded-file", "s", "uploaded.gob", "file name to save previously uploaded directories")
	viper.BindPFlag("upload.uploaded-file", cmd.Flags().Lookup("uploaded-file"))

	cmd.Flags().Bool("only-update-uploaded", false, "only update the contents of uploaded-file, do not make any requests")
	viper.BindPFlag("upload.only-update-uploaded", cmd.Flags().Lookup("only-update-uploaded"))

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
