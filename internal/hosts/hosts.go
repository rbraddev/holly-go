package hosts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrInvalidSearchType = errors.New("invalid search type")
	ErrSiteRequired      = errors.New("site required")
)

type HostList struct {
	Hosts []Host
}

type Host struct {
	ip       string
	hostname string
	platform string
}

type SearchParams struct {
	SearchType string
	Sites      []string
	Devices    []string
	Hostnames  []string
	Ips        []string
}

var constructorFuncs = map[string]func(sp SearchParams) (map[string]interface{}, error){
	"siteSearch": siteSearch,
}

func siteSearch(sp SearchParams) (map[string]interface{}, error) {
	if len(sp.Sites) == 0 {
		return map[string]interface{}{}, ErrSiteRequired
	}

	where := "WHERE "
	params := make(map[string]interface{})
	if len(sp.Sites) > 0 && len(sp.Devices) > 0 {
		id := int('a')
		for _, s := range sp.Sites {
			for _, d := range sp.Devices {
				key := fmt.Sprintf("%c", id)
				params[key] = fmt.Sprintf("%s%s%%", d, s)
				where += fmt.Sprintf("NodeName LIKE @%s OR ", key)
				id = id + 1
			}
		}
	}

	if where[len(where)-4:] == " OR " {
		where = where[:len(where)-4]
	}

	data := map[string]interface{}{
		"query":  fmt.Sprintf(`SELECT IPAddress as ip, NodeName as hostname FROM Orion.Nodes %s`, where),
		"params": params,
	}

	return data, nil
}

func (sp *SearchParams) constructQuery() (map[string]interface{}, error) {
	cf, ok := constructorFuncs[sp.SearchType]
	if !ok {
		return map[string]interface{}{}, ErrInvalidSearchType
	}

	query, err := cf(*sp)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return query, nil
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

func NewSolarwindsHosts(sp SearchParams) (HostList, error) {
	query, err := sp.constructQuery()
	if err != nil {
		return HostList{}, err
	}

	hl, err := getSolarwindsHosts(query)
	if err != nil {
		return HostList{}, err
	}
	// hl := HostList{
	// 	Hosts: []Host{
	// 		{ip: "10.0.0.1", hostname: "test1", platform: "ios"},
	// 		{ip: "10.0.0.2", hostname: "test2", platform: "ios"},
	// 	},
	// }

	return hl, nil
}

// func (h *Host) enableDisableSite(enable bool) error {

// }

// func (hl *HostList) LoadFile(hostsFile string) error {
// 	f, err := os.Open(hostsFile)
// 	if err != nil {
// 		if errors.Is(err, os.ErrNotExist) {
// 			return nil
// 		}
// 		return err
// 	}
// 	defer f.Close()

// 	scanner := bufio.NewScanner(f)

// 	for scanner.Scan() {
// 		hl.Hosts = append(hl.Hosts, scanner.Text())
// 	}

// 	return nil
// }
