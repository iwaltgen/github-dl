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
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	ggithub "github.com/google/go-github/github"
	"github.com/reactivex/rxgo/v2"
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
) (*ReleaseAsset, rxgo.Observable, error) {
	if err := repo.valid(); err != nil {
		return nil, nil, err
	}

	release, err := c.GetRelease(ctx, repo, opt.Tag)
	if err != nil {
		return nil, nil, err
	}

	asset := c.findReleaseAsset(release, opt)
	if asset == nil {
		err := fmt.Errorf("not found asset: [name: %s, os: %s, arch: %s]", opt.Name, opt.OS, opt.Arch)
		return nil, nil, err
	}

	observable, err := c.downloadAsset(asset, opt)
	return asset, observable, err
}

func (c *Client) findReleaseAsset(release *RepositoryRelease, opt *AssetOptions) *ReleaseAsset {
	for _, asset := range release.Assets {
		name := strings.ToLower(*asset.Name)
		matchedName := strings.Contains(name, opt.Name)
		matchedOS := strings.Contains(name, opt.OS)
		matchedArch := strings.Contains(name, opt.Arch)
		if matchedName && matchedOS && matchedArch {
			return &asset
		}
	}
	return nil
}

func (c *Client) downloadAsset(asset *ReleaseAsset, opt *AssetOptions) (rxgo.Observable, error) {
	url := *asset.BrowserDownloadURL
	filename := path.Base(url)
	filepath := filepath.Join(opt.DestPath, filename)
	tempext := ".ghdownload"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		defer resp.Body.Close()

		file, err := os.Create(filepath + tempext)
		if err != nil {
			next <- rxgo.Error(err)
			return
		}
		defer file.Close()

		counter := NewWriteCounter(next, int64(*asset.Size))
		if _, err = io.Copy(file, io.TeeReader(resp.Body, counter)); err != nil {
			next <- rxgo.Error(err)
			return
		}

		if err := os.Rename(filepath+tempext, filepath); err != nil {
			next <- rxgo.Error(err)
		}
	}}), nil

	// tempdir := os.TempDir()
	// https://github.com/mholt/archiver
	// return counter.Chan, nil
}
