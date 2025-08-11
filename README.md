# gh-runtime-cli

This is the GitHub CLI extension that adds support for the GitHub Spark and GitHub Spark Runtime products. It enables the creation and deletion of Spark Runtime apps (which back all Sparks), as well as retrieving information about the apps, and also the deployment/publishing of a bundled application to the respective app.

We make use of this in the GitHub Spark deployment process, where it's aliased as `gh spark` from `gh runtime`.

### Developing and Releasing

With our current mechanisms for pulling in this CLI extension into the relevant tooling for Spark, releases in this repository are not necessary. The only thing that needs to be done is merging to `main`.

### Versioning

The version of the CLI app is defined in the `cmd/version.go` and can be accessed via the `version` command.
Please update accordingly when making changes to the app.
