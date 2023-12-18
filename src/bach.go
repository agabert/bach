package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var sourceDirectory = "/space/escher/trunk"

	if _, isDebug := os.LookupEnv("DEBUG"); isDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug logging enabled.")
	}

	logrus.Info("main(): starting")

	var fileWalker = func(fullName string, fileInfo os.FileInfo, fileError error) error {
		if fileError != nil {
			logrus.Fatal(fileError)
			return fileError
		}

		if fileInfo.IsDir() {
			logrus.Debug("Skipping directory: ", fullName)
		} else {
			/*
			 * (I) check if the file is a normal file
			 */
			sourceFileStat, err := os.Stat(fullName)
			if err != nil {
				return nil
			}
			if !sourceFileStat.Mode().IsRegular() {
				return nil
			}
			sourceLinkStat, err := os.Lstat(fullName)
			if sourceLinkStat.Mode()&os.ModeSymlink != 0 {
				return nil
			}

			/*
			 * (II) prepare the fingerprint by reading the file once
			 */
			fingerPrint := getFileChecksum(fullName)

			/*
			 * (III) Check if the name on the box matches the content
			 */
			_, fileName := filepath.Split(fullName)

			if strings.TrimSuffix(fileName, filepath.Ext(fileName)) == fingerPrint {
				logrus.Debug("Filename: [", fullName, "] matches fingerprint: [", fingerPrint, "]")
			} else {
				logrus.Fatal("Filename [", fullName, "] has different content than fingerprint [", fingerPrint, "]")
			}
		}

		return nil
	}

	if err := filepath.Walk(sourceDirectory, fileWalker); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("main(): exiting")

	os.Exit(0)
}
