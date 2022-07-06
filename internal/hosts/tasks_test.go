package hosts

import (
	"fmt"
	"testing"
)

func TestSiteEnableDisableTask(t *testing.T) {
	h := Host{
		Ip:       "10.0.0.1",
		Hostname: "10.0.0.2",
	}

	config, err := h.EnableDisableSite("test", "test", true, []string{"10", "20", "30"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(config)
}
