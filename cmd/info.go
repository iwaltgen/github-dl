/*
Copyright © 2020 iwaltgen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/iwaltgen/github-dl/pkg/github"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show github repository release info.",
	Long: `Show github repository release info.

Example:
github-dl --repo iwaltgen/github-dl info
github-dl --repo iwaltgen/github-dl info --tag v0.1.0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client := github.NewClient(githubToken(), verbose)

		resp, err := client.GetRelease(ctx, github.Repository(repo), tag)
		if err != nil {
			return err
		}

		if verbose {
			color.Cyan("repository:\t%s", repo)
			color.Cyan("release tag:\t%s", tag)
		}

		return printPrettyJSON(Cyan, resp)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	flagSet := infoCmd.Flags()
	flagSet.StringVar(&tag, "tag", tag, "release tag")
}
