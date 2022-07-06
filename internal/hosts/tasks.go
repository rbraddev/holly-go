package hosts

import (
	"fmt"

	"github.com/rbraddev/holly-go/assets"
	"github.com/rbraddev/holly-go/internal/helpers"
)

type Result struct{}

type taskfunc func(map[string]any) (Result, error)

var tmplDir string = "templates"

func getTask(t string) (taskfunc, error) {
	switch t {
	case "enableDisableSite":
		return enableDisableSite, nil
	default:
		return nil, fmt.Errorf("invalid task: %s", t)
	}
}

func (hl *HostList) Run(t string, data map[string]any) error {

}

func enableDisableSite(data map[string]any) (Result, error) {
	_, err := helpers.ParseTemplate(assets.Files, tmplDir+"/enableDisableSite.tmpl", data)
	if err != nil {
		return Result{}, err
	}

	// err := h.SetDriver(hl.Username, hl.Password)
	// if err != nil {
	// 	return err
	// }

	// err = h.Open()
	// if err != nil {
	// 	return err
	// }
	// defer h.Driver.Close()

	// res, err := h.Driver.SendConfig(config)
	// if err != nil {
	// 	return err
	// }

	return Result{}, nil
}
