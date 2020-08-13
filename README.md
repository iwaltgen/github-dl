# GitHub-DL

![build](https://github.com/iwaltgen/github-dl/workflows/build/badge.svg)

`github-dl` downloads, search from GitHub release.

## Installation

```sh
go get -u github.com/iwaltgen/github-dl
```
Or download binaries from the [releases](https://github.com/iwaltgen/github-dl/releases) page.

## Usage

You'll need to export a GITHUB_TOKEN environment variable.
It will be used to fetch releases info from a GitHub repository.
You can create a token [here](https://github.com/settings/tokens) for GitHub.

```sh
export GITHUB_TOKEN="YOUR_GH_TOKEN"

github-dl --repo iwaltgen/github-dl [--tag, --asset, --os, --arch, --dest, --target, --pick]
github-dl --repo iwaltgen/github-dl list [--page, --per-page]
github-dl --repo iwaltgen/github-dl info [--tag]

github-dl help
github-dl help info
```
