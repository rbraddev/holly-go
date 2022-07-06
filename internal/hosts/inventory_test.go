package hosts

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/rbraddev/holly-go/internal/assert"
)

func TestNewInventory(t *testing.T) {
	testCases := []struct {
		name            string
		inventoryType   string
		expInventoryMod InventoryModule
		expErr          error
	}{
		{"SolarwindsInventory", "solarwinds", &SWInv{}, nil},
		{"HostFileInventory", "file", &HFInv{}, nil},
		{"NullInventory", "", nil, ErrMissingInventoryType},
		{"BogusInventory", "bogus", nil, fmt.Errorf("invalid inventory type: bogus")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			i, err := NewInventory(tc.inventoryType, Options{})

			if tc.expErr == nil {
				assert.NilError(t, err)
			} else {
				assert.ExpError(t, tc.expErr)
				return
			}

			assert.Equal(t, fmt.Sprintf("%v", reflect.TypeOf(i.Module)), fmt.Sprintf("%v", reflect.TypeOf(tc.expInventoryMod)))
		})
	}
}

func createHostsFile(t *testing.T) *os.File {
	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("error creating temp file: %s", err)
	}

	hosts := "10.0.0.1\n10.0.0.2\n10.0.0.3\n"
	err = os.WriteFile(tf.Name(), []byte(hosts), 0644)
	if err != nil {
		os.Remove(tf.Name())
		t.Fatal(err)
	}
	return tf
}

func TestGetFileHosts(t *testing.T) {
	tf := createHostsFile(t)
	defer os.Remove(tf.Name())

	hl := getFileHosts(tf)

	expHl := HostList{
		Hosts: []Host{
			{Hostname: "10.0.0.1", Ip: "10.0.0.1"},
			{Hostname: "10.0.0.2", Ip: "10.0.0.2"},
			{Hostname: "10.0.0.3", Ip: "10.0.0.3"},
		},
	}

	assert.Equal(t, len(hl.Hosts), len(expHl.Hosts))
	assert.Equal(t, fmt.Sprint(hl), fmt.Sprint(expHl))
}

func TestHostFileInventoryGet(t *testing.T) {
	tf := createHostsFile(t)
	defer os.Remove(tf.Name())

	fileHl := HostList{
		Hosts: []Host{
			{Hostname: "10.0.0.1", Ip: "10.0.0.1"},
			{Hostname: "10.0.0.2", Ip: "10.0.0.2"},
			{Hostname: "10.0.0.3", Ip: "10.0.0.3"},
		},
	}

	testCases := []struct {
		name   string
		file   string
		expHl  HostList
		expErr error
	}{
		{"ValidHostFile", tf.Name(), fileHl, nil},
		{"InvalidHostFile", "bogusfile", HostList{}, os.ErrNotExist},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := Options{File: tc.file, Username: "testuser", Password: "testpassword"}
			i, err := NewInventory("file", opts)
			if err != nil {
				t.Fatal(err)
			}

			hl, _ := i.Module.Get()

			if tc.expErr == nil {
				assert.NilError(t, err)
			} else {
				assert.ExpError(t, tc.expErr)
				return
			}

			assert.Equal(t, len(hl.Hosts), len(tc.expHl.Hosts))
			assert.Equal(t, fmt.Sprint(hl.Hosts), fmt.Sprint(tc.expHl.Hosts))
			assert.Equal(t, hl.Username, "testuser")
			assert.Equal(t, hl.Password, "testpassword")
		})
	}
}
