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

package github

import (
	"fmt"
	"strings"

	ggithub "github.com/google/go-github/github"
)

// Repository is github repository(owner/name)
type Repository string

// Owner is repository owner part
func (r Repository) Owner() string {
	parts := strings.Split(string(r), "/")
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}

// Name is repository name part
func (r Repository) Name() string {
	parts := strings.Split(string(r), "/")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func (r Repository) valid() error {
	parts := strings.Split(string(r), "/")
	if len(parts) != 2 {
		return fmt.Errorf("malformed github repository: %s", r)
	}
	return nil
}

// ListOptions specifies the optional parameters to various List methods that support pagination.
type ListOptions = ggithub.ListOptions

// RepositoryRelease represents a GitHub release in a repository.
type RepositoryRelease = ggithub.RepositoryRelease
