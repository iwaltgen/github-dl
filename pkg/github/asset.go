package github

import "math"

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

// AssetProgress is progress info of downloading a file.
type AssetProgress interface {
	Percentage() float64
}

// AssetTotalSize is the size of a file.
type AssetTotalSize uint64

// Percentage is downloaded size over total size.
func (p AssetTotalSize) Percentage() float64 {
	return 0
}

// AssetReceviedSize is the current downloaded size of a file.
type AssetReceviedSize struct {
	total    uint64
	received uint64
}

// Percentage is downloaded size over total size.
func (p AssetReceviedSize) Percentage() float64 {
	return math.Round(float64(p.received/p.total)*100) / 100
}
