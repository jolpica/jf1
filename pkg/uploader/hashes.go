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
func writeKnownHashesToFile(fileName string, knownHashes map[string][md5.Size]byte) {
	file, err := os.Create("test.gob")
	if err != nil {
		log.Fatalf("Failure encoding gob: %v", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(knownHashes); err != nil {
		log.Fatalf("failure encoding gob: %v", err)
	}

	if err := file.Close(); err != nil {
		log.Printf("warning failed to close file: %v", err)
	}
}
