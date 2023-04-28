// Package krdict manages schema and seeding for krdict dictionary
package krdict

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// DefaultRscPath is the default path of the rsc XML files
const DefaultRscPath = "seed/rsc/krdict"

// RscXMLPaths returns the paths of all XML files in rscPath
func RscXMLPaths(rscPath string) ([]string, error) {
	paths := []string{}
	files, err := os.ReadDir(rscPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".xml" {
			paths = append(paths, path.Join(rscPath, file.Name()))
		}
	}
	sort.Slice(paths, func(i, j int) bool {
		return xmlPathSort(paths[i]) < xmlPathSort(paths[j])
	})
	return paths, nil
}

var xmlPathSortRegexp = regexp.MustCompile("[0-9]+$")

func xmlPathSort(path string) int {
	s := xmlPathSortRegexp.FindString(strings.TrimSuffix(path, filepath.Ext(path)))
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i
}
