package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
)

type DependencyRetrieveResult struct {
	ExitCode int
	Message  string
}

func main() {
	username, packageID, installationKey := getAndValidateOpts()

	dependencies := getDependencies(username, packageID, installationKey)

	fmt.Printf("Found %d dependencies.", len(dependencies))

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

func getDependencies(username string, packageID string, installationKey string) []string {
	loader := startSpinner(" Retrieving dependencies...", "Dependencies retrieved!\n")

	soqlQuery := fmt.Sprintf("\"SELECT Dependencies FROM SubscriberPackageVersion WHERE Id='%s'\"", packageID)

	// Use tooling API to execute query.
	args := []string{
		"force:data:soql:query",
		"-u",
		username,
		"-t",
		"-q",
		soqlQuery,
		"--json"}

	retrieveResults, err := runSfCliCommand(args)
	dependencies := parseDependencyResponse(retrieveResults, err)

	loader.Stop()

	return dependencies
}

func runSfCliCommand(args []string) (string, error) {
	retrieveCommand := exec.Command("sfdx", args...)
	out, err := retrieveCommand.CombinedOutput()
	return string(out), err
}

func parseDependencyResponse(rawJSON string, err error) []string {
	var response DependencyRetrieveResult
	json.Unmarshal([]byte(rawJSON), &response)

	if err != nil {
		panic(response.Message)
	}

	return []string{}
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
