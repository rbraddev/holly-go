package helpers

import (
	"bytes"
	"fmt"
	"text/template"
)

func In[T comparable](value T, checklist []T) bool {
	for i := range checklist {
		if value == checklist[i] {
			return true
		}
	}
	return false
}

func ParseTemplate(f string, data map[string]any) (string, error) {
	t, err := template.ParseFiles(fmt.Sprintf("%s", f))
	if err != nil {
		return "", err
	}

	var tBytes bytes.Buffer
	err = t.Execute(&tBytes, data)
	if err != nil {
		return "", err
	}

	return tBytes.String(), nil
}
