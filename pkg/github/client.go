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
	"context"

	ggithub "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client is github oauth2 client
type Client struct {
	client *ggithub.Client
}

// NewClient is create github client
func NewClient(accessToken string) *Client {
	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return &Client{
		client: ggithub.NewClient(oauth2.NewClient(ctx, tokenSource)),
	}
}

// ListReleases is get release list
func (c *Client) ListReleases(ctx context.Context,
	repo Repository,
	opt *ListOptions,
) ([]*RepositoryRelease, error) {
	if err := repo.valid(); err != nil {
		return nil, err
	}

	releases, _, err := c.client.Repositories.ListReleases(ctx, repo.Owner(), repo.Name(), opt)
	return releases, err
}

// GetRelease is get release info
func (c *Client) GetRelease(ctx context.Context,
	repo Repository,
	id int64,
) (*RepositoryRelease, error) {
	if err := repo.valid(); err != nil {
		return nil, err
	}

	if id == 0 {
		release, _, err := c.client.Repositories.GetLatestRelease(ctx, repo.Owner(), repo.Name())
		return release, err
	}

	release, _, err := c.client.Repositories.GetRelease(ctx, repo.Owner(), repo.Name(), id)
	return release, err
}
