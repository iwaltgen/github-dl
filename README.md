# GitHub-DL

![build](https://github.com/iwaltgen/github-dl/workflows/build/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/iwaltgen/github-dl)](https://pkg.go.dev/github.com/iwaltgen/github-dl)

`github-dl` downloads, search from GitHub release.

## Installation

```sh
go install github.com/iwaltgen/github-dl
```
Or download binaries from the [releases](https://github.com/iwaltgen/github-dl/releases) page.

## Usage

You'll need to export a GITHUB_TOKEN environment variable.
It will be used to fetch releases info from a GitHub repository.
You can create a token [here](https://github.com/settings/tokens) for GitHub.

```sh
export GITHUB_TOKEN="YOUR_GH_TOKEN"

github-dl --repo iwaltgen/github-dl [--tag, --asset, --dest, --target, --pick]
github-dl --repo iwaltgen/github-dl list [--page, --per-page]
github-dl --repo iwaltgen/github-dl info [--tag]
github-dl --repo iwaltgen/github-dl --asset github-dl --pick github-dl

github-dl help
github-dl help info
```
