# github-dl

download assets from github release

## Installation

```sh
go get -u github.com/iwaltgen/github-dl
```

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
