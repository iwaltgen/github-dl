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
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
	"github.com/mholt/archiver/v3"
	"github.com/reactivex/rxgo/v2"
	"golang.org/x/oauth2"

	ggithub "github.com/google/go-github/github"
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
// first returns release asset info.
// second returns download progress info or error info use a stream.
// third returns initialize error info.
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
		matchedName := strings.Contains(name, strings.ToLower(opt.Name))

		matchedOS := strings.Contains(name, strings.ToLower(opt.OS))
		if !matchedOS {
			for _, v := range opt.OSAlias {
				if matchedOS = strings.Contains(name, v); matchedOS {
					break
				}
			}
		}

		matchedArch := strings.Contains(name, strings.ToLower(opt.Arch))
		if !matchedArch {
			for _, v := range opt.ArchAlias {
				if matchedArch = strings.Contains(name, v); matchedArch {
					break
				}
			}
		}

		if matchedName && matchedOS && matchedArch {
			return &asset
		}
	}
	return nil
}

func (c *Client) downloadAsset(asset *ReleaseAsset, opt *AssetOptions) (rxgo.Observable, error) {
	url := *asset.BrowserDownloadURL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		defer resp.Body.Close()

		filename := path.Base(url)
		destination := filepath.Join(opt.DestPath, filename)
		tempext := ".ghdownload"

		if err := os.MkdirAll(filepath.Dir(destination), os.ModePerm); err != nil {
			next <- rxgo.Error(err)
			return
		}
		file, err := os.Create(destination + tempext)
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

		if err := os.Rename(destination+tempext, destination); err != nil {
			next <- rxgo.Error(err)
		}
		defer func() {
			_ = os.Remove(destination + tempext)
			_ = os.Remove(destination)
		}()

		if !c.supportArchiveFile(filename) {
			if opt.Target != "" {
				newDestination := filepath.Join(opt.DestPath, opt.Target)
				if err := os.Rename(destination, newDestination); err != nil {
					next <- rxgo.Error(err)
				}
			}
			return
		}

		if opt.PickPattern != "" {
			c.extractFile(next, destination, opt)
			return
		}

		newDestination := filepath.Join(opt.DestPath, opt.Target)
		if err := archiver.Unarchive(destination, newDestination); err != nil {
			next <- rxgo.Error(err)
			return
		}
	}}), nil
}

func (c *Client) extractFile(ch chan<- rxgo.Item, filename string, opt *AssetOptions) {
	tempdir, err := ioutil.TempDir(os.TempDir(), "github-dl")
	if err != nil {
		ch <- rxgo.Error(err)
		return
	}
	defer func() {
		_ = os.RemoveAll(tempdir)
	}()

	if err := archiver.Unarchive(filename, tempdir); err != nil {
		ch <- rxgo.Error(err)
		return
	}

	matches, err := zglob.Glob(filepath.Join(tempdir, "**", opt.PickPattern))
	if err != nil {
		ch <- rxgo.Error(err)
		return
	}

	for i, path := range matches {
		if opt.Target != "" {
			suffix := ""
			if i != 0 {
				suffix = fmt.Sprintf(".%d", i)
			}

			destination := filepath.Join(opt.DestPath, opt.Target+suffix)
			if err := os.Rename(path, destination); err != nil {
				ch <- rxgo.Error(err)
				return
			}
			continue
		}

		filename := filepath.Base(path)
		destination := filepath.Join(opt.DestPath, filename)
		if err := os.Rename(path, destination); err != nil {
			ch <- rxgo.Error(err)
			return
		}
	}
}

func (c *Client) supportArchiveFile(filename string) bool {
	switch path.Ext(filename) {
	case ".zip", ".gz", ".tgz", ".br", ".zst", ".lz4", ".xz", ".sz":
		return true

	default:
		return false
	}
}
