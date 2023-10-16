// Package instagram extracts instagram posts to create Image db.SourcePartMedia
package instagram

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/util/archive/xz"
)

// GetLoginFromEnv gets the login from the default ENV var
func GetLoginFromEnv() string { return os.Getenv("INSTAGRAM_LOGIN") }

// NewFactory returns a new Factory
func NewFactory(login string) Factory { return Factory{login: login} }

// Factory generates Sources
type Factory struct{ login string }

// NewSource returns a new Post
func (f Factory) NewSource(url string) extractor.Source { return &Post{login: f.login, url: url} }

// Extensions returns the extensions the extractor returns
func (f Factory) Extensions() []string { return []string{".jpg"} }

// Post represents an instagram post
type Post struct {
	login  string
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
	cmd := exec.Command("instaloader", "--login", s.login, "--dirname-pattern", ".", "--", "-"+s.ID()) //nolint:gosec //this is how it works
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
func (s *Post) Info(cacheDir string) (db.PrePartInfo, error) {
	matches, err := filepath.Glob(filepath.Join(cacheDir, infoGlob))
	if err != nil {
		return db.PrePartInfo{}, err
	}
	if len(matches) != 1 {
		return db.PrePartInfo{}, fmt.Errorf("found != 1 files with glob (%v): %v", infoGlob, strings.Join(matches, ", "))
	}
	bytes, err := xz.Read(matches[0])
	if err != nil {
		return db.PrePartInfo{}, err
	}
	info := &postInfo{}
	if err := json.Unmarshal(bytes, info); err != nil {
		return db.PrePartInfo{}, err
	}
	username := info.Node.Owner.Username
	title := db.SourceDefaultedName(info.Node.EdgeMediaToCaption.Edges[0].Node.Text)
	return db.PrePartInfo{
		Name:      fmt.Sprintf("%v - %v", username, title),
		Reference: s.url,
	}, nil
}
