package hosts

import (
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
)

type Host struct {
	Hostname string
	Ip       string
	Driver   *network.Driver
}

func (h *Host) SetDriver(u, p string) error {
	pl, err := platform.NewPlatform(
		"cisco_iosxe",
		h.Ip,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(u),
		options.WithAuthPassword(p),
	)
	if err != nil {
		return err
	}

	d, err := pl.GetNetworkDriver()
	if err != nil {
		return err
	}

	h.Driver = d

	return nil
}

func (h *Host) Open() error {
	if h.Driver == nil {
		return ErrDriverNotSet
	}

	err := h.Driver.Open()
	if err != nil {
		return err
	}
	return nil
}

func (h *Host) Close() {
	h.Driver.Close()
}
