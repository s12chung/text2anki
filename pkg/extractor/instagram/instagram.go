// Package instagram extracts instagram posts to create Image db.SourcePartMedia
package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/util/archive/xz"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

var extensions = []string{".jpg"}

// GetLoginFromEnv gets the login from the default ENV var
func GetLoginFromEnv() string { return os.Getenv("INSTAGRAM_LOGIN") }

// NewFactory returns a new Factory
func NewFactory(login string) Factory { return Factory{login: login} }

// Factory generates Sources
type Factory struct{ login string }

// NewSource returns a new Post
func (f Factory) NewSource(url string) extractor.Source { return &Post{login: f.login, url: url} }

// Extensions returns the extensions the extractor returns
func (f Factory) Extensions() []string { return append([]string{}, extensions...) }

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

var extractToDirArgs = func(login, id string) []string {
	return []string{"instaloader", "--login", login, "--dirname-pattern", ".", "--", "-" + id}
}

// ExtractToDir extracts the post to the directory
func (s *Post) ExtractToDir(cacheDir string) error {
	if ok := s.Verify(); !ok {
		return fmt.Errorf("url is not verified for instagram: %v", s.url)
	}
	args := extractToDirArgs(s.login, s.ID())
	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec //this is how it works
	cmd.Dir = cacheDir
	if err := cmd.Run(); err != nil {
		return err
	}
	return numberPadFilenames(cacheDir)
}

func numberPadFilenames(cacheDir string) error {
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return err
	}
	filenames := ioutil.FilenamesWithExtensions(entries, extensions)

	for _, filename := range filenames {
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			return errors.New("file found with no underscore")
		}
		numberPart := strings.Split(parts[len(parts)-1], ".")[0]
		number, err := strconv.Atoi(numberPart)
		if err != nil {
			return err
		}
		newFilename := strings.Join(parts[:len(parts)-1], "_") + "_" + fmt.Sprintf("%03d", number) + filepath.Ext(filename)

		if err = os.Rename(filepath.Join(cacheDir, filename), filepath.Join(cacheDir, newFilename)); err != nil {
			return err
		}
	}
	return nil
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
