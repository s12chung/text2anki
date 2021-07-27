package anki

import (
	"os"
	"path"
)

var config Config

// SetConfig sets the Export and Cache config
func SetConfig(c Config) {
	config = c
}

// GetConfig returns the Export and Cache config
func GetConfig() Config {
	return config
}

// Config contains Export and Cache config for Anki
type Config struct {
	ExportPrefix  string
	NotesCacheDir string
}

const appCacheDirName = "Text2Anki"
const filesDirName = "files"

// DefaultConfig returns a default config
func DefaultConfig() (Config, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return Config{}, nil
	}
	return Config{
		ExportPrefix:  "t2a-",
		NotesCacheDir: path.Join(cacheDir, appCacheDirName, filesDirName),
	}, nil
}

// SetupDefaultConfig sets the default config and ensures the file paths exist
func SetupDefaultConfig() error {
	c, err := DefaultConfig()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(c.NotesCacheDir, 0750); err != nil {
		return err
	}
	SetConfig(c)
	return nil
}
