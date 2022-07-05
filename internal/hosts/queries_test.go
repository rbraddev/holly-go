package hosts

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/rbraddev/holly-go/internal/assert"
)

func TestNewConstructor(t *testing.T) {
	testCases := []struct {
		name       string
		searchtype string
		expFunc    func(i SWInv) (map[string]interface{}, error)
		expErr     error
	}{
		{"SiteSearch", "site", siteSearch, nil},
		{"NullSearch", "", nil, fmt.Errorf("search type required")},
		{"BogusSearch", "bogus", nil, fmt.Errorf("invalid search type: bogus")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cf, err := newConstructor(tc.searchtype)

			assert.Equal(t, fmt.Sprintf("%v", reflect.TypeOf(cf)), fmt.Sprintf("%v", reflect.TypeOf(tc.expFunc)))
			if tc.expErr == nil {
				assert.NilError(t, err)
			} else {
				assert.ExpError(t, tc.expErr)
			}
		})
	}
}

func TestSiteSearch(t *testing.T) {
	testCases := []struct {
		name      string
		opts      Options
		expQuery  string
		expParams map[string]string
		expErr    error
	}{
		{
			"SiteSearchSingleSiteSingleDevice",
			Options{Sites: []string{"123"}, Devices: []string{"sw"}},
			"SELECT IPAddress as ip, NodeName as hostname FROM Orion.Nodes WHERE NodeName LIKE @a",
			map[string]string{"a": "sw123%"},
			nil,
		},
		{
			"SiteSearchSingleSiteMultiDevice",
			Options{Sites: []string{"123"}, Devices: []string{"sw", "rt"}},
			"SELECT IPAddress as ip, NodeName as hostname FROM Orion.Nodes WHERE NodeName LIKE @a OR NodeName LIKE @b",
			map[string]string{"a": "sw123%", "b": "rt123%"},
			nil,
		},
		{
			"SiteSearchMultiSiteSingleDevice",
			Options{Sites: []string{"123", "124"}, Devices: []string{"sw"}},
			"SELECT IPAddress as ip, NodeName as hostname FROM Orion.Nodes WHERE NodeName LIKE @a OR NodeName LIKE @b",
			map[string]string{"a": "sw123%", "b": "sw124%"},
			nil,
		},
		{
			"SiteSearchMultiSiteMultiDevice",
			Options{Sites: []string{"123", "124"}, Devices: []string{"sw", "rt"}},
			"SELECT IPAddress as ip, NodeName as hostname FROM Orion.Nodes WHERE NodeName LIKE @a OR NodeName LIKE @b OR NodeName LIKE @c OR NodeName LIKE @d",
			map[string]string{"a": "sw123%", "b": "rt123%", "c": "sw124%", "d": "rt124%"},
			nil,
		},
		{
			"SiteSearchNoSite",
			Options{Sites: []string{}, Devices: []string{"sw"}},
			"",
			map[string]string{},
			ErrSiteRequired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			i := SWInv{Opts: tc.opts}
			result, err := siteSearch(i)

			if tc.expErr == nil {
				assert.NilError(t, err)
			} else {
				assert.ExpError(t, tc.expErr)
				return
			}

			query := result["query"].(string)
			params := result["params"].(map[string]string)

			assert.Equal(t, query, tc.expQuery)
			assert.Equal(t, len(params), len(tc.expParams))

			for k, v := range params {
				_, ok := tc.expParams[k]
				if !ok || v != tc.expParams[k] {
					t.Errorf("got: %s, expected: %s", params, tc.expParams)
				}
			}

		})
	}
}
