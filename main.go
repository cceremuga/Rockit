package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
)

// DependencyRetrieveResult wraps the JSON representation of dependency results in a SOQL query.
type DependencyRetrieveResult struct {
	Status  int
	Message string
}

// Manages all aspects of the package installation process.
func main() {
	username, packageID, installationKey := getAndValidateOpts()

	dependencies := getDependencies(username, packageID, installationKey)

	fmt.Printf("Found %d dependencies to install.\n", len(dependencies))

	fmt.Println("We didn't actually do anything, this is a work in progress.")
}

// Pulls all expected command-line flags, validates them.
func getAndValidateOpts() (string, string, string) {
	username, packageID, installationKey := getOpts()
	validateOpts(username, packageID)

	return username, packageID, installationKey
}

// Gets all command-line flags.
func getOpts() (string, string, string) {
	username := flag.String("u", "", "The target org username to install packages to.")
	packageID := flag.String("p", "", "The Id of the top-level package to install.")
	installationKey := flag.String("k", "", "An optional installation key for packages.")
	flag.Parse()

	return *username, *packageID, *installationKey
}

// Validates expected command-line flags.
func validateOpts(username string, packageID string) {
	if username == "" {
		panic("Target org username must be specified with the -u command-line flag.")
	}

	if packageID == "" {
		panic("Top-level package Id must be specified with the -p command-line flag.")
	}
}

// Retrieves all dependency packages for the top-level package Id.
func getDependencies(username string, packageID string, installationKey string) []string {
	loader := startSpinner(" Retrieving dependencies...", "")

	soqlQuery := fmt.Sprintf("SELECT Dependencies FROM SubscriberPackageVersion WHERE Id = '%s'", packageID)
	args := []string{
		"force:data:soql:query",
		"-u",
		username,
		"-t",
		"-q",
		soqlQuery,
		"--json"}

	// Use tooling API to execute query.
	retrieveResults, err := runSfCliCommand(args)
	dependencies := parseDependencyResponse(retrieveResults, err)

	loader.Stop()

	return dependencies
}

// Runs a Salesforce CLI command with the specified arguments.
func runSfCliCommand(args []string) (string, error) {
	retrieveCommand := exec.Command("sfdx", args...)
	out, err := retrieveCommand.CombinedOutput()
	return string(out), err
}

// Parses the JSON response from the dependency check.
func parseDependencyResponse(rawJSON string, err error) []string {
	var response DependencyRetrieveResult
	json.Unmarshal([]byte(rawJSON), &response)

	if err != nil || response.Status != 0 {
		panic(response.Message)
	}

	return []string{}
}

// Starts a cool loading indicator.
func startSpinner(suffix string, finalMsg string) *spinner.Spinner {
	loader := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	loader.Suffix = suffix
	loader.FinalMSG = finalMsg
	loader.HideCursor = true
	loader.Color("cyan")
	loader.Start()

	return loader
}
