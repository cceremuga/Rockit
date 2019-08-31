package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

const waitTime = "15"
const publishWait = "10"

// DependencyRetrieveResult wraps the JSON representation of dependency results in a SOQL query.
type DependencyRetrieveResult struct {
	Status  int
	Message string
	Result  struct {
		Records []struct {
			Attributes struct {
				Type string
			}
			Dependencies struct {
				IDs []struct {
					SubscriberPackageVersionID string
				}
			}
		}
	}
}

// Manages all aspects of the package installation process.
func main() {
	// Pull runtime options and validate them.
	username, packageIDs, installationKey := getAndValidateOpts()

	// Loop through all top-level packages, install them.
	installPackages(username, packageIDs, installationKey)
}

func installPackages(username string, packageIDs []string, installationKey string) {
	for _, packageID := range packageIDs {
		// Retrieve a list of dependencies for the top-level package.
		dependencies := getDependencies(username, packageID, installationKey)

		fmt.Printf("Found %d dependencies to install.\n", len(dependencies))

		if len(dependencies) > 0 {
			// Install all dependencies.
			installDependencies(dependencies, username, installationKey)
		}

		fmt.Println("Dependencies installed. Preparing to install top-level package.")

		// Install the top-level package.
		installPackage(username, packageID, installationKey)

		fmt.Println("You're all set, have a lovely day!")
	}
}

// Pulls all expected command-line flags, validates them.
func getAndValidateOpts() (string, []string, string) {
	username, packageIDs, installationKey := getOpts()
	validateOpts(username, packageIDs)

	return username, strings.Split(packageIDs, ","), installationKey
}

// Gets all command-line flags.
func getOpts() (string, string, string) {
	username := flag.String("u", "", "The target org username to install packages to.")
	packageIDs := flag.String("p", "", "The comma-separated Id(s) of the top-level package(s) to install.")
	installationKey := flag.String("k", "", "An optional installation key for packages.")
	flag.Parse()

	return *username, *packageIDs, *installationKey
}

// Validates expected command-line flags.
func validateOpts(username string, packageIDs string) {
	if username == "" {
		panic("Target org username must be specified with the -u command-line flag.")
	}

	if packageIDs == "" {
		panic("Top-level package(s) must be specified with the -p command-line flag, separated by commas.")
	}
}

// Retrieves all dependency packages for the top-level package Id.
func getDependencies(username string, packageID string, installationKey string) []string {
	loader := startSpinner(" Retrieving dependencies...", "")

	keyAppend := ""

	if installationKey != "" {
		keyAppend = fmt.Sprintf(" AND InstallationKey = '%s'", installationKey)
	}

	soqlQuery :=
		fmt.Sprintf("SELECT Dependencies FROM SubscriberPackageVersion WHERE Id = '%s'%s", packageID, keyAppend)

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

	// Extract the list of dependent packages from the query result.
	dependentPackages := response.Result.Records[0].Dependencies.IDs

	if len(dependentPackages) == 0 {
		// No dependencies found? Early exit.
		return []string{}
	}

	// Result had dependencies, build a slice of them.
	dependencies := []string{}

	for _, dependentPackage := range dependentPackages {
		dependencies = append(dependencies, dependentPackage.SubscriberPackageVersionID)
	}

	return dependencies
}

// Installs all dependent packages.
func installDependencies(packageIDs []string, username string, installationKey string) {
	for _, packageID := range packageIDs {
		installPackage(username, packageID, installationKey)
	}
}

// Installs a package.
func installPackage(username string, packageID string, installationKey string) {
	fmt.Printf("Starting install for %s.\n", packageID)

	args := []string{
		"force:package:install",
		"--package",
		packageID,
		"-u",
		username,
		"-w",
		waitTime,
		"--publishwait",
		publishWait,
		"--securitytype",
		"AllUsers"} // TODO Support security types.

	if installationKey != "" {
		args = append(args, "--installationkey", installationKey)
	}

	packageInstallCommand := exec.Command("sfdx", args...)
	packageInstallCommand.Stdout = os.Stdout
	packageInstallCommand.Stderr = os.Stderr
	err := packageInstallCommand.Run()

	if err != nil {
		log.Fatalf("Package install failed with %s.\n", err)
	}
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
