// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	version = "0.1.5"
	ldflags = "-ldflags=-s -w"
)

type Version mg.Namespace

var (
	goexe = "go"
	git   = sh.RunCmd("git")
)

func init() {
	goexe = mg.GoCmd()
}

// Build build app
func Build() error {
	mg.Deps(Lint)

	args := []string{"build", "-trimpath", "-ldflags", "-s -w", "-o", "./build/github-dl"}
	return sh.RunV(goexe, args...)
}

// Clean clean build artifacts
func Clean() {
	sh.Rm("build")
	sh.Rm("dist")
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

func existsFile(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// Release creates release (current version)
func (v Version) Release() error {
	mg.Deps(v.Tag)
	return sh.RunV("goreleaser", "--rm-dist")
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
	return v.bumpVersion(nextVer.String())
}

// Minor bump minor version
func (v Version) Minor() error {
	curVer, _ := semver.NewVersion(version)
	nextVer := curVer.IncMinor()
	return v.bumpVersion(nextVer.String())
}

// Patch bump patch version
func (v Version) Patch() error {
	curVer, _ := semver.NewVersion(version)
	nextVer := curVer.IncPatch()
	return v.bumpVersion(nextVer.String())
}

func (v Version) bumpVersion(version string) error {
	files := []string{"magefile.go", "cmd/version.go"}

	hash, _ := sh.Output("git", "rev-parse", "--verify", "HEAD")
	kv := map[*regexp.Regexp]string{
		regexp.MustCompile(`version\s+= \"([0-9]+)\.([0-9]+)\.([0-9]+)\"`): fmt.Sprintf(`version = "%s"`, version),
		regexp.MustCompile(`commitHash\s= \"[a-z0-9]+\"`):                  fmt.Sprintf(`commitHash = "%s"`, hash),
		regexp.MustCompile(`modifiedAt\s= \"[0-9]+\"`):                     fmt.Sprintf(`modifiedAt = "%d"`, time.Now().Unix()),
	}

	for _, file := range files {
		if err := v.replaceFileText(file, kv); err != nil {
			return fmt.Errorf("bump version replace `%s` error: %w", file, err)
		}
		if err := sh.Run(goexe, "fmt", file); err != nil {
			return fmt.Errorf("bump version fmt `%s` error: %w", file, err)
		}
	}

	for _, file := range files {
		if err := git("add", file); err != nil {
			return fmt.Errorf("git add `%s` error: %w", file, err)
		}
	}

	color.Green("new version: %s", version)
	return git("commit", "-m", "chore: bump version")
}

func (Version) replaceFileText(path string, kv map[*regexp.Regexp]string) error {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	contents := string(read)
	for k, v := range kv {
		contents = k.ReplaceAllString(contents, v)
	}
	return ioutil.WriteFile(path, []byte(contents), 0)
}
