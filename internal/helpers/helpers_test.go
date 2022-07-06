package helpers

import (
	"os"
	"testing"

	"github.com/rbraddev/holly-go/internal/assert"
)

func TestParseTemplate(t *testing.T) {
	ParseTemplateSuccessData := map[string]any{"enable": false, "vlans": []string{"10", "20"}}
	ParseTemplateSuccessExpResult := `
interface vlan 10
  shutdown
interface vlan 20
  shutdown`

	testCases := []struct {
		name         string
		templateDir  string
		templateFile string
		templateData map[string]any
		expResult    string
		expErr       error
	}{
		{"ParseTemplateSuccess", "testdata", "testTemplate.tmpl", ParseTemplateSuccessData, ParseTemplateSuccessExpResult, nil},
		{"ParseTemplateNoFile", "testdata", "bogus.tmpl", nil, "", os.ErrNotExist},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file := os.DirFS(tc.templateDir)
			templateStr, err := ParseTemplate(file, tc.templateFile, tc.templateData)

			if tc.expErr == nil {
				assert.NilError(t, err)
			} else {
				assert.ExpError(t, tc.expErr)
				return
			}

			assert.Equal(t, templateStr, tc.expResult)

		})
	}
}
