package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func getFileChecksum(fPath string) string {
	fHandle, err := os.Open(fPath)

	if err != nil {
		logrus.Fatal(err)
	}

	defer fHandle.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, fHandle); err != nil {
		logrus.Fatal(err)
	}

	return fmt.Sprintf("%x", hasher.Sum(nil))
}
