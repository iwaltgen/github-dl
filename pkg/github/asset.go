package github

// AssetOptions are parameters to download an asset file.
type AssetOptions struct {
	Tag         string
	Name        string
	OS          string
	OSAlias     []string
	Arch        string
	ArchAlias   []string
	DestPath    string
	Target      string
	PickPattern string
}
