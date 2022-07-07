package hosts

import (
	"fmt"

	"github.com/rbraddev/holly-go/assets"
	"github.com/rbraddev/holly-go/internal/helpers"
	"github.com/scrapli/scrapligo/response"
)

const global bool = true

var tmplDir string = "templates"

func (hl *HostList) SendConfig(g bool) {
	var config string

	errCh := make(chan error)
	resCh := make(chan response.Response)
	doneCh := make(chan struct{})

	hl.wg.Add(len(hl.Hosts))

	for _, h := range hl.Hosts {
		go func(h Host) {
			defer hl.wg.Done()

			if g {
				config = hl.GlobalConfig
			} else {
				config = h.Config
			}

			err := h.SetDriver(hl.Username, hl.Password)
			if err != nil {
				errCh <- err
			}

			err = h.Open()
			if err != nil {
				errCh <- err
			}
			defer h.Driver.Close()

			res, err := h.Driver.SendConfig(config)
			if err != nil {
				errCh <- err
			}

			resCh <- *res

		}(h)
	}

	go func() {
		hl.wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			fmt.Println(err)
		case res := <-resCh:
			fmt.Println(res)
		case <-doneCh:
			return
		}
	}
}

func (hl *HostList) EnableDisableSite(enable bool, vlans []string) error {
	var err error
	data := map[string]any{
		"enable": enable,
		"vlans":  vlans,
	}
	hl.GlobalConfig, err = helpers.ParseTemplate(assets.Files, tmplDir+"/enableDisableSite.tmpl", data)
	if err != nil {
		return err
	}

	hl.SendConfig(global)

	return nil
}
