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
)

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
