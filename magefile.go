// +build mage

package main

import (
	"fmt"
	"os"
	"time"

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
