package github

// AssetOptions are parameters to download an asset file.
type AssetOptions struct {
	Tag      string
	Name     string
	OS       string
	Arch     string
	DestPath string
	DestFile string
	PickFile string
}
