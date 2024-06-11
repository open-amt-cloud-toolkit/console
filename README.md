![Build](https://img.shields.io/github/actions/workflow/status/open-amt-cloud-toolkit/console/ci.yml?style=for-the-badge&label=Build&logo=github)
![Codecov](https://img.shields.io/codecov/c/github/open-amt-cloud-toolkit/console?style=for-the-badge&logo=codecov)
[![OSSF-Scorecard Score](https://img.shields.io/ossf-scorecard/github.com/open-amt-cloud-toolkit/console?style=for-the-badge&label=OSSF%20Score)](https://api.securityscorecards.dev/projects/github.com/open-amt-cloud-toolkit/console)
[![Discord](https://img.shields.io/discord/1063200098680582154?style=for-the-badge&label=Discord&logo=discord&logoColor=white&labelColor=%235865F2&link=https%3A%2F%2Fdiscord.gg%2FDKHeUNEWVH)](https://discord.gg/DKHeUNEWVH)
# Console


> Disclaimer: Production viable releases are tagged and listed under 'Releases'.  All other check-ins and pre-releases should be considered 'in-development' and should not be used in production

## Overview

This is an application that packages the UI, RPS, and MPS into a single executable for use in an enterprise environment.

## Quick start 

### For users

If you're looking for the latest release of console visit [Github Releases](https://github.com/open-amt-cloud-toolkit/console/releases/latest) and download the appropriate binary assets for your OS and Architecture: https://github.com/open-amt-cloud-toolkit/console/releases/latest. This is the quickest way to get up and running for non-developers.

## For Developers

Local development (in Linux or WSL):

To start the service with Postgres: 

```sh
# Postgres
$ make compose-up
# Run app with migrations
$ make run
```

Download and check out the sample-web-ui:
```
git clone https://github.com/open-amt-cloud-toolkit/sample-web-ui
```

Ensure that the environment file has cloud set to `false` and that the URLs for RPS and MPS are pointing to where you have `Console` running. The default is `http://localhost:8181`. Follow the instructions for launching and running the UI in the sample-web-ui readme.






## Dev tips for passing CI Checks

- Install gofumpt `go install mvdan.cc/gofumpt@latest` (replaces gofmt)
- Install gci `go install github.com/daixiang0/gci@latest` (organizes imports)
- Ensure code is formatted correctly with `gofumpt -l -w -extra ./`
- Ensure code is gci'd with `gci.exe write --skip-generated -s standard -s default .`
- Ensure all unit tests pass with `go test ./...`
- Ensure code has been linted with `docker run --rm -v ${pwd}:/app -w /app golangci/golangci-lint:latest golangci-lint run -v`
