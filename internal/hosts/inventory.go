package hosts

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

type HostList struct {
	Username string
	Password string
	Hosts    []Host
	wg       *sync.WaitGroup
}

type Options struct {
	Username   string
	Password   string
	SearchType string
	Sites      []string
	Devices    []string
	File       string
}

type InventoryModule interface {
	Get() (*HostList, error)
}

type SWInv struct {
	Opts Options
}

type HFInv struct {
	Opts Options
}

type Inventory struct {
	InventoryType string
	Module        InventoryModule
}

func getSolarwindsHosts(query map[string]interface{}) (HostList, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return HostList{}, err
	}

	requestBody := bytes.NewBuffer(queryJson)

	fmt.Println(requestBody)

	return HostList{}, nil
}

func (i *SWInv) Get() (*HostList, error) {
	query, err := i.constructQuery()
	if err != nil {
		return nil, err
	}
	hl, err := getSolarwindsHosts(query)
	if err != nil {
		return nil, err
	}
	return &hl, nil
}

func getFileHosts(f io.Reader) HostList {
	scanner := bufio.NewScanner(f)

	var hl HostList
	for scanner.Scan() {
		ip := scanner.Text()
		hl.Hosts = append(hl.Hosts, Host{Hostname: ip, Ip: ip})
	}
	return hl
}

func (i *HFInv) Get() (*HostList, error) {
	f, err := os.Open(i.Opts.File)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}
	defer f.Close()

	hl := getFileHosts(f)
	hl.Username = i.Opts.Username
	hl.Password = i.Opts.Password

	return &hl, nil
}

func NewInventory(i string, o Options) (*Inventory, error) {
	if i == "" {
		return nil, ErrMissingInventoryType
	}
	switch i {
	case "solarwinds":
		return &Inventory{
			InventoryType: "solarwinds",
			Module:        &SWInv{Opts: o},
		}, nil
	case "file":
		return &Inventory{
			InventoryType: "hostfile",
			Module:        &HFInv{Opts: o},
		}, nil
	default:
		return nil, fmt.Errorf("invalid inventory type: %s", i)
	}
}
