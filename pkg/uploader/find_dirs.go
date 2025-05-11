package uploader

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ChangedFile struct {
	Path     string
	Hash     [md5.Size]byte
	Contents []byte
}
type DirWithChangedFiles struct {
	Path         string
	ChangedFiles []ChangedFile
}

type checkedDirectory struct {
	Path     string
	DirEntry fs.DirEntry
}

// Return a channel of top level directories in the given path
func getDirs(ctx context.Context, rootPath string) (<-chan checkedDirectory, chan error) {
	dirsc := make(chan checkedDirectory)
	errc := make(chan error, 1)

	go func() {
		defer close(dirsc)
		errc <- filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() || path == rootPath {
				// Only Walk directories
				return nil
			}

			select {
			case dirsc <- checkedDirectory{Path: path, DirEntry: d}:
				// Don't search nested directories
				return fs.SkipDir
			case <-ctx.Done():
				return errors.New("walk cancelled")
			}
		})
	}()

	return dirsc, errc
}

// Return a channel of directories which have changes from the known hashes
func findUpdatedDirs(ctx context.Context, knownHashes map[string][md5.Size]byte, dirsc <-chan checkedDirectory) <-chan DirWithChangedFiles {
	wg := sync.WaitGroup{}
	updatedDirsc := make(chan DirWithChangedFiles)
	numWorkers := 10

	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			for dir := range dirsc {
				err := checkDirForUpdates(ctx, knownHashes, dir, updatedDirsc)
				if err != nil {
					fmt.Printf("err: %v\n", err)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(updatedDirsc)
	}()

	return updatedDirsc
}

// For a given directory, compare its contents against the known hashes map
func checkDirForUpdates(ctx context.Context, knownHashes map[string][md5.Size]byte, dir checkedDirectory, updatedDirsc chan<- DirWithChangedFiles) error {
	// TODO: Should a changed dir include all files, or just the changed ones? (currently only changed ones)
	changedFiles := []ChangedFile{}
	err := filepath.WalkDir(dir.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		checksum := md5.Sum(contents)

		if checksum != knownHashes[path] {
			changedFiles = append(changedFiles, ChangedFile{Path: path, Hash: checksum, Contents: contents})
		}

		return nil
	})
	if len(changedFiles) > 0 {
		select {
		case updatedDirsc <- DirWithChangedFiles{Path: dir.Path, ChangedFiles: changedFiles}:
		case <-ctx.Done():
			return errors.New("canceled checking dir")
		}
	}

	return err

}
