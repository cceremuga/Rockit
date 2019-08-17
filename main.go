package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

func main() {
	username, packageID, installationKey := getAndValidateOpts()

	getDependencies(username, packageID, installationKey)

	fmt.Println("We didn't actually do anything, this is a work in progress.")
}

func getAndValidateOpts() (string, string, string) {
	username, packageID, installationKey := getOpts()
	validateOpts(username, packageID)

	return username, packageID, installationKey
}

func getOpts() (string, string, string) {
	username := flag.String("u", "", "The target org username to install packages to.")
	packageID := flag.String("p", "", "The Id of the top-level package to install.")
	installationKey := flag.String("k", "", "An optional installation key for packages.")
	flag.Parse()

	return *username, *packageID, *installationKey
}

func validateOpts(username string, packageID string) {
	if username == "" {
		panic("Target org username must be specified with the -u command-line flag.")
	}

	if packageID == "" {
		panic("Top-level package Id must be specified with the -p command-line flag.")
	}
}

func getDependencies(username string, packageID string, installationKey string) {
	loader := startSpinner(" Retrieving dependencies...", "Dependencies retrieved!\n")
	time.Sleep(4 * time.Second)
	loader.Stop()
}

func startSpinner(suffix string, finalMsg string) *spinner.Spinner {
	loader := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	loader.Suffix = suffix
	loader.FinalMSG = finalMsg
	loader.HideCursor = true
	loader.Color("cyan")
	loader.Start()

	return loader
}
