package hosts

import "fmt"

func newConstructor(s string) (func(i SWInv) (map[string]interface{}, error), error) {
	if s == "" {
		return nil, fmt.Errorf("search type required")
	}

	switch s {
	case "site":
		return siteSearch, nil
	default:
		return nil, fmt.Errorf("invalid search type: %s", s)
	}
}

func (i SWInv) constructQuery() (map[string]interface{}, error) {
	cf, err := newConstructor(i.Opts.SearchType)
	if err != nil {
		return map[string]interface{}{}, err
	}

	query, err := cf(i)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return query, nil
}

func siteSearch(i SWInv) (map[string]interface{}, error) {
	if len(i.Opts.Sites) == 0 {
		return map[string]interface{}{}, ErrSiteRequired
	}

	where := "WHERE "
	params := make(map[string]string)
	if len(i.Opts.Sites) > 0 && len(i.Opts.Devices) > 0 {
		id := int('a')
		for _, s := range i.Opts.Sites {
			for _, d := range i.Opts.Devices {
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
