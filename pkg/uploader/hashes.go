package uploader

import (
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

func readKnownHashesFromFile(fileName string) (knownHashes map[string][md5.Size]byte) {
	knownHashes = make(map[string][16]byte)
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		fmt.Printf("Running with no known files as known hashes file not found (%v)\n", fileName)
		return knownHashes
	}
	if err != nil {
		log.Fatalf("Failure opening gob: %v", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&knownHashes); err != nil {
		log.Fatalf("Fail to decode gob: %v", err)
	}

	return knownHashes
}
func writeKnownHashesToFile(fileName string, knownHashes map[string][md5.Size]byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failure encoding gob: %v", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(knownHashes); err != nil {
		return fmt.Errorf("failure encoding gob: %v", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failure to close file: %v", err)
	}
	return nil
}
