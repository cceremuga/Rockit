# Salesforce Package Installer

[![Build Status](https://travis-ci.org/cceremuga/sf-package-installer.svg?branch=master)](https://travis-ci.org/cceremuga/sf-package-installer)

A Salesforce unlocked package installer with support for automatic dependency detection and installation.

## Requirements

* Salesforce CLI
* Golang

## Usage

* Authorize the destination org with the Salesfordce CLI.
* `go build` in the directory where the repository lives.
* `chmod +x ./sf-package-installer`
* `./sf-package-installer -u target_org_username@example.com -p packageIdGoesHere -k optionalInstallKeyGoesHere`

## License

MIT License. See LICENSE for more info.