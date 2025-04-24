package uploader

import (
	"crypto/md5"
	"fmt"
	"maps"
)

type UploadConfig struct {
	BaseUrl               string `mapstructure:"base-url"`
	DryRun                bool   `mapstructure:"dry-run"`
	ScannedFile           string `mapstructure:"scanned-file"`
	MaxConcurrentRequests int    `mapstructure:"max-concurrent-requests"`
}

func RunUploader(dirsPath string, config UploadConfig, token string) error {
	done := make(chan struct{})

	knownHashes := readKnownHashesFromFile(config.ScannedFile)

	dirsc, errc := getDirs(done, dirsPath)

	updatedDirsc := findUpdatedDirs(done, knownHashes, dirsc)

	dirLoadResultc := loadDataFromDirectories(done, updatedDirsc)

	requestResultc := sendDataLoadRequests(done, dirLoadResultc, config, token)

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
			fmt.Printf("Failed to make a request for %s: %v\n\n", result.ProcessedDir.SourceDirPath, result.Err)
			continue
		}

		if result.StatusCode >= 300 {
			fmt.Printf("FAILURE %q (%v): %+v\n\n", result.RequestData.Description, result.StatusCode, result.ResponseData)
			continue
		}

		requestData := result.RequestData
		responseData := result.ResponseData
		fmt.Printf("SUCCESS (dry_run:%v) uploading %v\n", requestData.DryRun, requestData.Description)
		fmt.Printf("Total: %v, Created: %v, Updated %v\n\n", responseData.TotalCount, responseData.CreatedCount, responseData.UpdatedCount)

		maps.Copy(newHashes, result.ProcessedDir.FileChecksums)
	}

	if !config.DryRun {
		err := writeKnownHashesToFile(config.ScannedFile, newHashes)
		if err != nil {
			return err
		}
	}
	return nil
}
