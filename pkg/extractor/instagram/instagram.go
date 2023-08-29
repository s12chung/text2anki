// Package instagram extracts instagram posts to create Image db.SourcePartMedia
package instagram

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/util/archive"
	"github.com/s12chung/text2anki/pkg/util/stringutil"
)

// Factory generates Sources
type Factory struct{}

// NewSource returns a new Post
func (f Factory) NewSource(url string) extractor.Source {
	return &Post{url: url}
}

// Extensions returns the extensions the extractor returns
func (f Factory) Extensions() []string {
	return []string{".jpg"}
}

// Post represents an instagram post
type Post struct {
	url    string
	verify *bool
	id     string
}

const hostname = "www.instagram.com"
const pathPrefix = "/p/"

// Verify returns true if the url is an instagram post url
func (s *Post) Verify() bool {
	if s.verify != nil {
		return *s.verify
	}

	u, err := url.Parse(s.url)
	if err != nil {
		return false
	}

	ok := u.Hostname() == hostname && strings.HasPrefix(u.Path, pathPrefix)
	s.verify = &ok
	if !ok {
		return ok
	}
	s.id = strings.TrimSuffix(strings.TrimPrefix(u.Path, pathPrefix), "/")
	return ok
}

// ID is the ID of the post
func (s *Post) ID() string {
	s.Verify()
	return s.id
}

// ExtractToDir extracts the post to the directory
func (s *Post) ExtractToDir(cacheDir string) error {
	if ok := s.Verify(); !ok {
		return fmt.Errorf("url is not vertified for instagram: %v", s.url)
	}
	cmd := exec.Command("instaloader", "--dirname-pattern", ".", "--", "-"+s.ID()) //nolint:gosec //this is how it works
	cmd.Dir = cacheDir
	return cmd.Run()
}

const infoGlob = "*.xz"

type postInfo struct {
	Node struct {
		Owner struct {
			Username string `json:"username"`
		} `json:"owner"`

		EdgeMediaToCaption struct {
			Edges []struct {
				Node struct {
					Text string `json:"text"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"edge_media_to_caption"`
	} `json:"node"`
}

// Info returns the info given from the extraction
func (s *Post) Info(cacheDir string) (extractor.SourceInfo, error) {
	matches, err := filepath.Glob(filepath.Join(cacheDir, infoGlob))
	if err != nil {
		return extractor.SourceInfo{}, err
	}
	if len(matches) != 1 {
		return extractor.SourceInfo{}, fmt.Errorf("found != 1 files with glob (%v): %v", infoGlob, strings.Join(matches, ", "))
	}
	bytes, err := archive.XZBytes(matches[0])
	if err != nil {
		return extractor.SourceInfo{}, err
	}
	info := &postInfo{}
	if err := json.Unmarshal(bytes, info); err != nil {
		return extractor.SourceInfo{}, err
	}
	username := info.Node.Owner.Username
	title := stringutil.FirstUnbrokenSubstring(info.Node.EdgeMediaToCaption.Edges[0].Node.Text, 30)
	return extractor.SourceInfo{
		Name:      fmt.Sprintf("%v - %v", username, title),
		Reference: s.url,
	}, nil
}
