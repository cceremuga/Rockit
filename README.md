# Salesforce Package Installer

[![Build Status](https://travis-ci.org/cceremuga/sf-package-installer.svg?branch=master)](https://travis-ci.org/cceremuga/sf-package-installer)

A Salesforce unlocked package installer with support for automatic dependency detection and installation.

## Requirements

* Salesforce CLI

## Usage

* Be in a terminal.
* Authorize the destination org with the Salesfordce CLI, make note of the username.
* `cd` to wherever you cloned the repository.
* `chmod +x ./sf-package-installer`
* `./sf-package-installer -u target_org_username@example.com -p packageIdGoesHere -k optionalInstallKeyGoesHere`

## Release Notes

* 8/17/19 - Initial release with base functionality including dependency installation.

## Known Issues

* No support for package install keys yet.
* Only compiled binary is for Darwin (Mac OS) architecture.

## Building

* `cd` to wherever you cloned the repository.
* `go get` to install depencencies.
* `go build` to compile.

## License

MIT License. See LICENSE for more info.