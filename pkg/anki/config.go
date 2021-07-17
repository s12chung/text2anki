package anki

import (
	"log"
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

var cacheDir string

const appCacheDirName = "Text2Anki"
const filesDirName = "files"

func init() {
	var err error
	cacheDir, err = os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
}

// DefaultConfig returns a default config
func DefaultConfig() Config {
	return Config{
		ExportPrefix:  "t2a-",
		NotesCacheDir: path.Join(cacheDir, appCacheDirName, filesDirName),
	}
}

// SetupDefaultConfig sets the default config and ensures the file paths exist
func SetupDefaultConfig() error {
	config := DefaultConfig()
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil
	}
	if err := os.Mkdir(path.Join(cacheDir, appCacheDirName), 0750); err != nil {
		return nil
	}
	if err := os.Mkdir(path.Join(cacheDir, appCacheDirName, filesDirName), 0750); err != nil {
		return nil
	}
	SetConfig(config)
	return nil
}
