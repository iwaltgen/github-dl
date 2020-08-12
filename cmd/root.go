/*
Copyright © 2020 iwaltgen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
	"github.com/reactivex/rxgo/v2"
	"github.com/spf13/cobra"

	"github.com/iwaltgen/github-dl/pkg/github"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github-dl",
	Short: "Download a github repository release asset.",
	Long: fmt.Sprintf(`Download a github repository release asset.

version: %s
commit: %s
build: %s

Example:
github-dl --repo cli/cli --asset gh --dest bin --pick gh
github-dl --repo golangci/golangci-lint --asset golangci-lint --pick golangci-lint
github-dl --repo uber/prototool --asset prototool --target prototool
github-dl --repo google/protobuf --asset protoc --target protoc --pick protoc`,
		version,
		commitHash,
		buildTime().Format(time.RFC3339),
	),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client := github.NewClient(githubToken())

		opt, err := makeAssetOptions()
		if err != nil {
			color.Magenta(err.Error())
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}

		if verbose {
			color.Cyan("repository:\t%s", repo)
			color.Cyan("release:\t%s", tag)
		}

		asset, observable, err := client.DownloadReleaseAsset(ctx, github.Repository(repo), opt)
		if err != nil {
			return err
		}

		return showDownloadProgress(ctx, asset, observable)
	},
}

var (
	verbose  bool
	tokenEnv string
	token    string
	repo     string
)

var (
	asset     string
	tag       = "latest"
	osname    = runtime.GOOS
	osAlias   = "darwin:macos,osx;windows:win"
	arch      = runtime.GOARCH
	archAlias = "amd64:x86_64"
	dest, _   = os.Getwd()
	target    string
	pick      string
)

func init() {
	pflagSet := rootCmd.PersistentFlags()
	pflagSet.BoolVarP(&verbose, "verbose", "v", verbose, "verbose output")
	pflagSet.StringVar(&tokenEnv, "token-env", "GITHUB_TOKEN", "github oauth2 token environment name")
	pflagSet.StringVar(&token, "token", token, "github oauth2 token value (optional)")
	pflagSet.StringVar(&repo, "repo", repo, "github repository (owner/name)")

	flagSet := rootCmd.Flags()
	flagSet.StringVar(&asset, "asset", asset, "asset name keyword")
	flagSet.StringVar(&tag, "tag", tag, "release tag")
	flagSet.StringVar(&osname, "os", osname, "os keyword")
	flagSet.StringVar(&osAlias, "os-alias", osAlias, "os keyword alias")
	flagSet.StringVar(&arch, "arch", arch, "arch keyword")
	flagSet.StringVar(&archAlias, "arch-alias", archAlias, "arch keyword alias")
	flagSet.StringVar(&dest, "dest", dest, "destination path")
	flagSet.StringVar(&target, "target", target, "rename destination file (optional)")
	flagSet.StringVar(&pick, "pick", pick, "extract archive and pick a file name pattern (optional)")
}

func githubToken() string {
	if token != "" {
		return token
	}
	return os.Getenv(tokenEnv)
}

func makeAssetOptions() (*github.AssetOptions, error) {
	if asset == "" {
		return nil, errors.New("require asset name: see flags --asset")
	}

	osAliasMap, err := parseAlias(osAlias)
	if err != nil {
		return nil, errors.New("parse alias error: see flags --os-alias")
	}

	archAliasMap, err := parseAlias(archAlias)
	if err != nil {
		return nil, errors.New("parse alias error: see flags --arch-alias")
	}

	return &github.AssetOptions{
		Name:        asset,
		Tag:         tag,
		OS:          osname,
		OSAlias:     osAliasMap[osname],
		Arch:        arch,
		ArchAlias:   archAliasMap[arch],
		DestPath:    dest,
		Target:      target,
		PickPattern: pick,
	}, nil
}

func parseAlias(flagAlias string) (map[string][]string, error) {
	ret := map[string][]string{}
	aliases := strings.Split(flagAlias, ";")
	for _, alias := range aliases {
		kv := strings.Split(alias, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("parse alias: %v", kv)
		}
		k, v := kv[0], strings.Split(kv[1], ",")
		ret[k] = v
	}
	return ret, nil
}

func showDownloadProgress(ctx context.Context,
	asset *github.ReleaseAsset,
	observable rxgo.Observable,
) error {
	totalSize := int64(*asset.Size)
	pbbar := pb.Full.New(int(totalSize))
	pbbar.Set(pb.Bytes, true)
	pbbar.Set(pb.Terminal, true)

	pbbar.Start()
	for item := range observable.Observe(rxgo.WithContext(ctx)) {
		if item.Error() {
			return item.E
		}

		progress := item.V.(*github.DownloadProgress)
		pbbar.SetCurrent(progress.Received)
	}
	pbbar.SetCurrent(totalSize)
	pbbar.Finish()

	return nil
}
