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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Fetch github repository release list.",
	Long: `Fetch github repository release list.

Example:
github-dl --repo iwaltgen/github-dl list
github-dl --repo iwaltgen/github-dl list --page 1 --per-page 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client := github.NewClient(githubToken())
		opt := &github.ListOptions{
			Page:    page,
			PerPage: perPage,
		}

		resp, err := client.ListReleases(ctx, github.Repository(repo), opt)
		if err != nil {
			return err
		}

		if !verbose {
			var results []*repoRelease
			for _, v := range resp {
				results = append(results, &repoRelease{
					ID:   *v.ID,
					Name: *v.Name,
					Tag:  *v.TagName,
					URL:  *v.HTMLURL,
				})
			}

			return printPrettyJSON(Cyan, results)
		}

		color.Cyan("repository:\t%s", repo)
		color.Cyan("page-num:\t%d", page)
		color.Cyan("per-page:\t%d", perPage)

		return printPrettyJSON(Cyan, resp)
	},
}

type repoRelease struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Tag  string `json:"tag,omitempty"`
	URL  string `json:"url,omitempty"`
}

var (
	page    int
	perPage int
)

func init() {
	rootCmd.AddCommand(listCmd)

	flagSet := listCmd.Flags()
	flagSet.IntVar(&page, "page", 1, "request page number")
	flagSet.IntVar(&perPage, "per-page", 10, "request per page count")
}
