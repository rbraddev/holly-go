package hosts

import "fmt"

type HostList struct {
	Hosts []Host
}

type Host struct {
	ip       string
	hostname string
	platform string
}

type SearchParams struct {
	Site     []string
	Device   []string
	Hostname []string
	Ip       []string
}

func (sp *SearchParams) constructQuery() string {
	where := "WHERE "
	params := make(map[string]any)

	if len(sp.Site) > 0 && len(sp.Device) > 0 {
		for i, s := range sp.Site {
			for _, d := range sp.Device {
				id := fmt.Sprint(('a' + (i - 1)))
				params[id] = fmt.Sprintf("%s%s%%", d, s)
				for k := range params {
					where += fmt.Sprintf("NodeName LIKE @%s OR ", k)
				}
			}
		}
	}

	fmt.Println(params)

	return fmt.Sprintf(`
		SELECT NodeID as nodeid, IPAddress as ip, NodeName as hostname
		FROM Orion.Nodes
		%s
	`, where)
}

func LoadSolarwindsHosts(sp SearchParams) {
	query := sp.constructQuery()
	fmt.Println(query)
}

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
