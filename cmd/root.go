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
	"fmt"
	"os"
	"runtime"
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
github-dl --repo iwaltgen/github-dl --asset github-dl
github-dl --repo uber/prototool --asset prototool --target prototool
github-dl --repo golangci/golangci-lint --asset golangci-lint --target golangci-lint --pick golangci-lint
github-dl --repo google/protobuf --asset protoc --os osx --dest ./bin --target protoc --pick bin/protoc`,
		version,
		commitHash,
		buildTime().Format(time.RFC3339),
	),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client := github.NewClient(githubToken())

		if verbose {
			color.Cyan("repository:\t%s", repo)
			color.Cyan("release:\t%s", tag)
		}

		if asset == "" {
			color.Magenta("require asset name: see flags --asset")
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}

		asset, observable, err := client.DownloadReleaseAsset(ctx, github.Repository(repo), &github.AssetOptions{
			Name:     asset,
			Tag:      tag,
			OS:       osname,
			Arch:     arch,
			DestPath: dest,
			Target:   target,
			PickFile: pick,
		})
		if err != nil {
			return err
		}

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
	},
}

var (
	verbose  bool
	tokenEnv string
	token    string
	repo     string
)

var (
	asset  string
	osname string
	arch   string
	dest   string
	target string
	pick   string
)

func init() {
	pflagSet := rootCmd.PersistentFlags()
	pflagSet.BoolVarP(&verbose, "verbose", "v", verbose, "verbose output")
	pflagSet.StringVar(&tokenEnv, "token-env", "GITHUB_TOKEN", "github oauth2 token environment name")
	pflagSet.StringVar(&token, "token", token, "github oauth2 token value (optional)")
	pflagSet.StringVar(&repo, "repo", repo, "github repository (owner/name)")

	flagSet := rootCmd.Flags()
	wd, _ := os.Getwd()
	flagSet.StringVar(&tag, "tag", tag, "release tag")
	flagSet.StringVar(&asset, "asset", asset, "asset name keyword")
	flagSet.StringVar(&osname, "os", runtime.GOOS, "os keyword")
	flagSet.StringVar(&arch, "arch", runtime.GOARCH, "arch keyword")
	flagSet.StringVar(&dest, "dest", wd, "destination path")
	flagSet.StringVar(&target, "target", target, "destination file (optional)")
	flagSet.StringVar(&pick, "pick", pick, "extract archive and pick a file (optional)")
}

func githubToken() string {
	if token != "" {
		return token
	}
	return os.Getenv(tokenEnv)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Magenta(err.Error())
		os.Exit(1)
	}
}
