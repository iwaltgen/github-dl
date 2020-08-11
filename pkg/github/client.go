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
	"fmt"

	ggithub "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client is a github oauth2 client.
type Client struct {
	client *ggithub.Client
}

// NewClient creates github client.
func NewClient(accessToken string) *Client {
	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return &Client{
		client: ggithub.NewClient(oauth2.NewClient(ctx, tokenSource)),
	}
}

// ListReleases get release list.
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

// GetRelease gets release info.
func (c *Client) GetRelease(ctx context.Context,
	repo Repository,
	tag string,
) (*RepositoryRelease, error) {
	if err := repo.valid(); err != nil {
		return nil, err
	}

	if tag == "latest" {
		release, _, err := c.client.Repositories.GetLatestRelease(ctx, repo.Owner(), repo.Name())
		return release, err
	}

	release, _, err := c.client.Repositories.GetReleaseByTag(ctx, repo.Owner(), repo.Name(), tag)
	return release, err
}

// DownloadReleaseAsset downloads a release asset file.
func (c *Client) DownloadReleaseAsset(ctx context.Context,
	repo Repository,
	opt *AssetOptions,
) (<-chan AssetProgress, error) {
	if err := repo.valid(); err != nil {
		return nil, err
	}

	release, err := c.GetRelease(ctx, repo, opt.Tag)
	if err != nil {
		return nil, err
	}

	asset := c.findReleaseAsset(release, opt)
	if asset == nil {
		return nil, fmt.Errorf("not found asset: [name: %s, os: %s, arch: %s]", opt.Name, opt.OS, opt.Arch)
	}

	// TODO(iwaltgen): download file
	// https://golangcode.com/download-a-file-with-progress/
	// https://github.com/mholt/archiver
	return nil, nil
}

func (c *Client) findReleaseAsset(release *RepositoryRelease, opt *AssetOptions) *ReleaseAsset {
	for _, asset := range release.Assets {
		// TODO(iwaltgen): find asset match opt
		_ = asset
	}
	return nil
}
