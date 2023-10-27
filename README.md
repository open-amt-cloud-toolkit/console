# Go Github Action Workflow Templates

This template repository contains recommended github action workflows for go based applications.
## What Workflows are included?

- [CodeQL Analysis](#codeql-analysis)
- [CodeQL Report](#codeql-report)
- [Go Build](#go-build)
- [Docker Build/Publish](#docker-buildpublish)
- [Conventional Commits](#conventional-commits)
- [Trivy Container Scan](#trivy-scan)
- [Dependabot Configuration](#dependabot-configuration)

### CodeQL Analysis - `codeql-analysis.yml`
#### Description

Used for static code analysis. See https://codeql.github.com/ for more information.

### CodeQL Report - `codeql-report.yml`
#### Prerequistes/Configuration Required

- A Github Secret called `SECURITY_TOKEN` populated with a Github PAT is required with scopes of `public_repo` and `security_events`.

#### Description

Used for automatically uploading CodeQL Analysis to Github Artifacts. Useful for providing evidence of scans.

### Go Build - `go.yml`

#### Description

This is the main file for Go builds. This is a matrixed build, meaning it will simultaneously run the steps on multiple versions of Go. For example 1.18.x,1.19.x,1.20.x of Go. We do this to ensure compability with the supported versions of go see https://go.dev/doc/devel/release#policy. 

The build includes the following tasks:

- Running Unit Tests
- Ensuring code is formatted using gofmt
- Ensuring code is vetted using go vet 
- Ensuring code is linted using [golangci/golangci-lint-action@v3](https://github.com/golangci/golangci-lint-action/tree/v3/)
- Uploading Test Results (in JUnit format)
- Uploading Code Coverage results using codecov.io

### Docker Build/Publish -`docker-build.yml`

#### Prerequistes/Configuration Required

The following files must be available in your repository:

- `Dockerfile` - contains the instructions for building the Docker image
- `docker-compose.yml` (optional) - if you are using docker-compose for multi-container Docker applications
  You should also have your container registry username and password stored as GitHub Secrets as `DOCKER_USERNAME` and `DOCKER_PASSWORD`. This allows the GitHub Actions to push the Docker image to your container repository.

Ensure the registry of your choice and the app name is specified in the defaults:

```yaml
 ...
 inputs:
      docker_registry:
            description: 'Registry URL'
            required: true
            default: 'docker.io/username' # update to your private registry
      image_name:
            description: 'Name you wish to use on the docker image (ex. myapp). This will be tagged with :latest, and the git sha'
            required: true
            default: 'app' # update to your app name
```

#### Description

This workflow builds a Docker image from the Dockerfile in your repository, and then pushes that image to a specified container registry.

The build and publish includes the following tasks:

- Building the Docker image
- Logging in to DockerHub
- Tagging the Docker image with the commit hash and 'latest' tag
- Pushing the Docker image to DockerHub

### Trivy Container Scan - `trivy-scan.yml`

#### Prerequistes/Configuration Required

Ensure a `Dockerfile` is available to build and provide to trivvy.

#### Description

Trivy is a comprehensive open-source vulnerability scanner for containers. It detects vulnerabilities in OS packages (Alpine, RHEL, CentOS, etc.) and application dependencies (NPM, pip, etc.). This GitHub Actions workflow uses Trivy to scan your Docker image for any known vulnerabilities and provides a report which can be viewed directly in the GitHub Actions interface.

### Dependabot Configuration - `dependabot.yml`

#### Prerequistes/Configuration Required

No specific configuration is required for this workflow to run. However, you need to ensure that your project has a valid package.json file, as Dependabot relies on it to check for outdated dependencies. Optionally, you may change the schedule as needed from daily to something else.

#### Description

Dependabot is a tool that checks your project dependencies for any known security vulnerabilities or updates. It can automatically create pull requests to update your dependencies to the latest versions. This GitHub Actions workflow configures Dependabot for your Go project. It is highly recommended to keep your dependencies up to date not just to benefit from the latest features and improvements, but also to avoid potential security risks associated with outdated packages.

The Dependabot Configuration includes the following tasks:

- Daily check for outdated packages
- Automatic pull request creation for outdated packages
- Optional automatic merge for minor and patch updates of packages
- Security advisories notifications for packages.

### Semantic Pull Request - `semantic.yml`

#### Prerequistes/Configuration Required

The following files must be available in your repository:

`commitlint.config.js` - configuration file for CommitLint which contains the rules for commit messages. You can customize the list of allowed scopes in the `scope-enum` rule:

```js
'scope-enum': [
    2,
    'always',
    [] // add scopes here to enforce when scope is provided (ex. ['core','api','startup'])
]
```

#### Description

The Semantic Pull Request workflow enforces a set of standards as defined by conventional commits for all pull requests and commit messages in your repository. It helps ensure that your project maintains a consistent and clean commit history. The standards are documented here https://www.conventionalcommits.org/en/v1.0.0/

This GitHub Actions workflow is triggered whenever a pull request is opened, edited, reopened, or synchronized. It uses the CommitLint tool to validate the commit messages and pull request title against the rules defined in your commitlint.config.js file.
