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
	"strconv"
	"time"
)

var (
	version    = "0.1.2"
	commitHash = "d9dad2041f4d8e8f3c409a8794b502d21ddf9e45"
	modifiedAt = "1597304062"
)

func lastModified() time.Time {
	return unixStringToTime(modifiedAt)
}

func unixStringToTime(unixStr string) time.Time {
	i, _ := strconv.ParseInt(unixStr, 10, 64)
	return time.Unix(i, 0).UTC()
}
