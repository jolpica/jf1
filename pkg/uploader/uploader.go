package uploader

import (
	"crypto/md5"
	"fmt"
	"maps"
)

func RunUploader(dirsPath string, hashFileName string, baseUrl string, dryRun bool) error {
	done := make(chan struct{})

	knownHashes := readKnownHashesFromFile(hashFileName)

	dirsc, errc := getDirs(done, dirsPath)

	updatedDirsc := findUpdatedDirs(done, knownHashes, dirsc)

	dirLoadResultc := loadDataFromDirectories(done, updatedDirsc)

	requestResultc := sendDataLoadRequests(done, dirLoadResultc, baseUrl, dryRun)

	if err := saveAndDisplayResults(requestResultc, knownHashes, hashFileName, dryRun); err != nil {
		return err
	}

	return <-errc
}

func saveAndDisplayResults(requestResultc <-chan RequestResult, knownHashes map[string][md5.Size]byte, hashFileName string, dryRun bool) error {
	newHashes := make(map[string][md5.Size]byte)
	maps.Copy(newHashes, knownHashes)
	for result := range requestResultc {
		if result.Err != nil {
			fmt.Printf("Failed to make a request for %s: %v\n\n", result.ProcessedDir.SourceDirPath, result.Err)
			continue
		}

		if result.StatusCode >= 300 {
			fmt.Printf("FAILURE Request (%v): %+v\n\n", result.StatusCode, result.ResponseData)
			continue
		}

		requestData := result.RequestData
		responseData := result.ResponseData
		fmt.Printf("SUCCESS (dry_run:%v) uploading %v\n", requestData.DryRun, requestData.Description)
		fmt.Printf("Total: %v, Created: %v, Updated %v\n\n", responseData.TotalCount, responseData.CreatedCount, responseData.UpdatedCount)

		maps.Copy(newHashes, result.ProcessedDir.FileChecksums)
	}

	if !dryRun {
		writeKnownHashesToFile(hashFileName, newHashes)
	}
	return nil
}
