package uploader

import (
	"context"
	"crypto/md5"
	"fmt"
	"maps"
)

type UploadConfig struct {
	Verbose               bool
	BaseUrl               string `mapstructure:"base-url"`
	DryRun                bool   `mapstructure:"dry-run"`
	MaxConcurrentRequests int    `mapstructure:"max-concurrent-requests"`

	UploadedFile       string `mapstructure:"uploaded-file"`
	OnlyUpdateUploaded bool   `mapstructure:"only-update-uploaded"`
}

func RunUploader(dirPaths []string, config UploadConfig, token string) error {
	ctx := context.Background()

	knownHashes := readKnownHashesFromFile(config.UploadedFile)

	dirsc, errc := getDirs(ctx, dirPaths)

	updatedDirsc := findUpdatedDirs(ctx, knownHashes, dirsc)

	dirLoadResultc := loadDataFromDirectories(ctx, updatedDirsc)

	requestResultc := sendDataLoadRequests(ctx, dirLoadResultc, config, token)

	if err := saveAndDisplayResults(requestResultc, knownHashes, config); err != nil {
		return err
	}

	return <-errc
}

func saveAndDisplayResults(requestResultc <-chan RequestResult, knownHashes map[string][md5.Size]byte, config UploadConfig) error {
	newHashes := make(map[string][md5.Size]byte)
	maps.Copy(newHashes, knownHashes)
	for result := range requestResultc {
		if result.Err != nil {
			fmt.Printf("\nFailed to make a request for %s: %v\n", result.ProcessedDir.SourceDirPath, result.Err)
			continue
		}

		if result.StatusCode >= 300 {
			fmt.Printf("\nFAILURE (%v) %q: %+v\n", result.StatusCode, result.RequestData.Description, result.ResponseData.Errors)
			continue
		}

		requestData := result.RequestData
		responseData := result.ResponseData
		if !config.OnlyUpdateUploaded {
			fmt.Printf("\nSUCCESS (dry_run:%v) uploading %v\n", requestData.DryRun, requestData.Description)
			fmt.Printf("Total: %v, Created: %v, Updated %v\n", responseData.TotalCount, responseData.CreatedCount, responseData.UpdatedCount)
		}

		maps.Copy(newHashes, result.ProcessedDir.FileChecksums)
	}
	fmt.Println()

	if !config.DryRun {
		err := writeKnownHashesToFile(config.UploadedFile, newHashes)
		if err != nil {
			return err
		}
		if config.OnlyUpdateUploaded {
			fmt.Printf("Successfully updated %s with the current directory contents\n", config.UploadedFile)
		}
	} else {
		fmt.Println("Skipped saving scanned files as dry-run is enabled")
	}
	return nil
}
