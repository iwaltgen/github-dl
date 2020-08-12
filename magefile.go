// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	packageName = "github.com/iwaltgen/github-dl"
	version     = "0.1.0"
	ldflags     = "-ldflags=-s -w" +
		" -X $PACKAGE/cmd.version=$VERSION" +
		" -X $PACKAGE/cmd.commitHash=$COMMIT_HASH" +
		" -X $PACKAGE/cmd.buildDate=$BUILD_DATE"
)

type Version mg.Namespace

var (
	goexe = "go"
	git   = sh.RunCmd("git")

	workspace = packageName
	started   = time.Now().Unix()
)

func init() {
	goexe = mg.GoCmd()
	workspace, _ = os.Getwd()
}

func buildEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--verify", "HEAD")
	return map[string]string{
		"PACKAGE":     packageName,
		"WORKSPACE":   workspace,
		"VERSION":     version,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  fmt.Sprintf("%d", started),
	}
}

// Build build app
func Build() error {
	mg.Deps(Lint)

	args := []string{"build", "-trimpath", ldflags, "-o", "./build/github-dl"}
	return sh.RunWith(buildEnv(), goexe, args...)
}

// Clean clean build artifacts
func Clean() error {
	return sh.Rm("build")
}

// Test test app
func Test() error {
	return sh.RunV("sh", "-c", "go test ./pkg/... -cover -json | tparse -all")
}

// Lint lint app
func Lint() error {
	return sh.RunV("golangci-lint", "run", "--timeout", "3m", "-E", "misspell")
}

// All clean, lint, test, build app
func All() {
	mg.Deps(Clean, Lint, Test, Build)
}

type goget struct{}

func (g goget) installModule(uri string) error {
	env := map[string]string{"GO111MODULE": "off"}
	return sh.RunWith(env, goexe, "get", uri)
}

func existsFile(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// Tag creates tag (current version)
func (v Version) Tag() error {
	tag := "v" + version
	if err := git("tag", "-a", tag, "-m", tag+" release"); err != nil {
		return fmt.Errorf("add git tag error: %w", err)
	}
	return git("push", "origin", tag)
}

// Major bump major version
func (v Version) Major() error {
	curVer, _ := semver.NewVersion(version)
	nextVer := curVer.IncMajor()
	return v.bumpVersion(version, nextVer.String())
}

// Minor bump minor version
func (v Version) Minor() error {
	curVer, _ := semver.NewVersion(version)
	nextVer := curVer.IncMinor()
	return v.bumpVersion(version, nextVer.String())
}

// Patch bump patch version
func (v Version) Patch() error {
	curVer, _ := semver.NewVersion(version)
	nextVer := curVer.IncPatch()
	return v.bumpVersion(version, nextVer.String())
}

func (v Version) bumpVersion(old, new string) error {
	files := []string{"magefile.go"}
	for _, file := range files {
		if err := v.replaceFileText(file, old, new); err != nil {
			return fmt.Errorf("bump version `%s` error: %w", file, err)
		}
	}

	for _, file := range files {
		if err := git("add", file); err != nil {
			return fmt.Errorf("git add `%s` error: %w", file, err)
		}
	}

	color.Green("new version: %s", new)
	return git("commit", "-m", "chore: bump version")
}

func (Version) replaceFileText(path, old, new string) error {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	newContents := strings.Replace(string(read), old, new, -1)
	return ioutil.WriteFile(path, []byte(newContents), 0)
}
