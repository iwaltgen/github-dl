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
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show github-dl version",
	Long:  `Show github-dl version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", version)
		fmt.Printf("commit: %s\n", commitHash)
		fmt.Printf("build: %s\n", buildTime().Format(time.RFC3339))
	},
}

var (
	version    = "dev"
	commitHash = "dev"
	buildDate  = ""
	startTime  time.Time
)

func init() {
	startTime = time.Now()
}

func buildTime() time.Time {
	buildTime, err := unixStringToTime(buildDate)
	if err != nil {
		return startTime
	}
	return buildTime
}

func unixStringToTime(unixStr string) (t time.Time, err error) {
	i, err := strconv.ParseInt(unixStr, 10, 64)
	if err != nil {
		return t, fmt.Errorf("parse unix timestamp string: %w", err)
	}
	return time.Unix(i, 0).UTC(), nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
