package db

import (
	"encoding/json"
	"io"
	"io/fs"
)

// PrePartListFile is the fileTree for PreParts
type PrePartListFile struct {
	InfoFile fs.File               `json:"info_file,omitempty"`
	PreParts []SourcePartMediaFile `json:"pre_parts"`
}

// PrePartList is a KeyTree for PreParts
type PrePartList struct {
	ID       string            `json:"id"`
	InfoKey  string            `json:"info_key,omitempty"`
	PreParts []SourcePartMedia `json:"pre_parts"`
}

// PrePartInfo contains Source related info within PrePartList.InfoKey
type PrePartInfo struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

// Info returns the PrePartInfo that InfoKey represents
func (p PrePartList) Info() (PrePartInfo, error) {
	if p.InfoKey == "" {
		return PrePartInfo{}, nil
	}

	f, err := dbStorage.Get(p.InfoKey)
	if err != nil {
		return PrePartInfo{}, err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return PrePartInfo{}, err
	}
	info := PrePartInfo{}
	if err := json.Unmarshal(b, &info); err != nil {
		return PrePartInfo{}, err
	}
	return info, nil
}

// PrePartListURL represents all the Source parts together for a given id
type PrePartListURL struct {
	ID       string            `json:"id"`
	PreParts []PrePartMediaURL `json:"pre_parts"`
}

// StaticCopy returns a copy without fields that variate
func (p PrePartListURL) StaticCopy() PrePartListURL { return p }

// PrePartMediaURL represents a SourcePartMedia before it is created, only stored via. Routes.Storage.Storer
type PrePartMediaURL struct {
	ImageURL string `json:"image_url,omitempty"`
	AudioURL string `json:"audio_url,omitempty"`
}
