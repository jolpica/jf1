package uploader

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"
)

type DirectoryLoadResult struct {
	Path   string
	Result *ProcessedDirectory
	Err    error
}

type ProcessedDirectory struct {
	Data          []map[string]any
	SourceDirPath string
	FileChecksums map[string][md5.Size]byte
}

func (data *ProcessedDirectory) Description() string {
	if data == nil {
		return "<nil ImportData>"
	}
	return fmt.Sprintf("%s (%d files, %d data entries)", data.SourceDirPath, len(data.FileChecksums), len(data.Data))
}

func loadDataFromDirectories(ctx context.Context, changedDirsc <-chan DirWithChangedFiles) <-chan DirectoryLoadResult {
	dirLoadResultc := make(chan DirectoryLoadResult, 5)
	maxConcurrentDirs := 10
	sem := make(chan struct{}, maxConcurrentDirs)
	go func() {
		var wg sync.WaitGroup
		defer func() { wg.Wait(); close(dirLoadResultc) }()
		for updatedDir := range changedDirsc {
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				fmt.Printf("Cancelled while waiting for semaphore, stopping dir imports")
				return
			}
			wg.Add(1)
			go func(dirWithChanges DirWithChangedFiles) {
				defer wg.Done()
				defer func() { <-sem }()
				fmt.Printf("Found changes: %q\n", dirWithChanges.Path)
				processedDir, err := processDirectoryFiles(dirWithChanges)
				select {
				case dirLoadResultc <- DirectoryLoadResult{Path: dirWithChanges.Path, Result: processedDir, Err: err}:
				case <-ctx.Done():
					fmt.Printf("Cancelled sending directory data for %s", dirWithChanges.Path)
				}
			}(updatedDir)
		}
	}()

	return dirLoadResultc
}

func processDirectoryFiles(dirWithChanges DirWithChangedFiles) (processedDir *ProcessedDirectory, err error) {
	data := []map[string]any{}
	fileChecksums := make(map[string][md5.Size]byte, len(dirWithChanges.ChangedFiles))

	for _, file := range dirWithChanges.ChangedFiles {
		var fileData []map[string]any
		if err := json.Unmarshal(file.Contents, &fileData); err != nil {
			return nil, fmt.Errorf("umarshalling directory %q failed: %v", dirWithChanges.Path, err)
		}
		data = append(data, fileData...)
		fileChecksums[file.Path] = file.Hash
	}
	return &ProcessedDirectory{SourceDirPath: dirWithChanges.Path, Data: data, FileChecksums: fileChecksums}, nil
}
