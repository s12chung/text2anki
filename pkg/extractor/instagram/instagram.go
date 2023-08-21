// Package instagram extracts instagram posts to create Image db.SourcePartMedia
package instagram

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/s12chung/text2anki/pkg/extractor"
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
