# Salesforce Package Installer

[![Build Status](https://travis-ci.org/cceremuga/sf-package-installer.svg?branch=master)](https://travis-ci.org/cceremuga/sf-package-installer)

A Salesforce unlocked package installer with support for automatic dependency detection and installation.

## Requirements

* Salesforce CLI

## Usage

* Be in a terminal.
* Authorize the destination org with the Salesfordce CLI, make note of the username.
* `cd` to wherever you cloned the repository.
* `cd` to `bin`
* `ls` to find the binary in the appropriate architecture for your system.
* `chmod +x ./sf-package-installer-darwin`
* `./sf-package-installer-darwin -u target_org_username@example.com -p packageIdGoesHere -k optionalInstallKeyGoesHere`

## Release Notes

* 8/23/19 - Untested beta support for installation keys. Compiled binaries for all supported architectures.
* 8/17/19 - Initial release with base functionality including dependency installation.

## Known Issues

* None

## Building

* `cd` to wherever you cloned the repository.
* `go get` to install depencencies.
* `chmod +x build.sh`
* `./build.sh` to compile for all architectures. Output is to `bin`.

## License

MIT License. See LICENSE for more info.