package testdb

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

type testSeeder struct{ seedCount int }

func (t *testSeeder) Name() string              { return "" }
func (t *testSeeder) Filename() string          { return "" }
func (t *testSeeder) ReadFile() ([]byte, error) { return nil, nil }
func (t *testSeeder) Seed(_ db.TxQs) error {
	t.seedCount++
	return nil
}

func TestSeedList(t *testing.T) {
	testName := "TestSeedList"

	testCases := []struct {
		name string
		list map[string]bool
	}{
		{name: "all"},
		{name: "blacklist", list: map[string]bool{"Notes": false, "Persons": false}},
		{name: "whitelist", list: map[string]bool{"Notes": true, "Persons": true}},
		{name: "mixed", list: map[string]bool{"Notes": true, "Persons": true, "Books": false, "Cups": false}},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			sMap := map[string]seeder{}
			for _, k := range []string{"Books", "Cups", "Entries", "Notes", "Persons"} {
				sMap[k] = &testSeeder{}
			}
			require.NoError(seedList(db.TxQs{}, tc.list, sMap))

			resultMap := map[string]int{}
			for k, v := range sMap {
				s, ok := v.(*testSeeder)
				require.True(ok)
				resultMap[k] = s.seedCount
			}
			fixture.CompareReadOrUpdateJSON(t, path.Join(testName, tc.name), resultMap)
		})
	}
}
