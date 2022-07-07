package hosts

import (
	"testing"
)

func TestSiteEnableDisableTask(t *testing.T) {
	hl := HostList{
		Hosts: []Host{
			{Ip: "10.0.0.1", Hostname: "10.0.0.1"},
			{Ip: "10.0.0.2", Hostname: "10.0.0.2"},
		},
	}

	err := hl.EnableDisableSite(true, []string{"10", "20", "30"})
	if err != nil {
		t.Error(err)
	}
}
