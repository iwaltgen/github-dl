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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github-dl",
	Short: "Download github repository release assets.",
	Long: `Download github repository release assets.

Example: TBD
λ> github-dl iwaltgen/github-dl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root called")
		fmt.Println(githubToken())
	},
}

var (
	verbose  bool
	tokenEnv string
	token    string
	repo     string
)

func init() {
	flagSet := rootCmd.PersistentFlags()
	flagSet.BoolVarP(&verbose, "verbose", "v", verbose, "verbose output")
	flagSet.StringVar(&tokenEnv, "token-env", "GITHUB_TOKEN", "github oauth2 token environment name")
	flagSet.StringVar(&token, "token", "", "github oauth2 token value")
	flagSet.StringVar(&repo, "repo", "", "github repository (owner/name)")
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
		fmt.Println(err)
		os.Exit(1)
	}
}
