package i18n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

type Translator struct {
	messages map[string]string
}

func Load(localesDir, lang string) (*Translator, error) {
	path := fmt.Sprintf("%s/%s.json", localesDir, lang)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read locale %s: %w", path, err)
	}

	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("parse locale %s: %w", path, err)
	}

	return &Translator{messages: messages}, nil
}

func (t *Translator) Render(key string, data any) (string, error) {
	tmplStr, ok := t.messages[key]
	if !ok {
		return "", fmt.Errorf("i18n key not found: %s", key)
	}

	tmpl, err := template.New(key).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("parse template %s: %w", key, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render template %s: %w", key, err)
	}

	return buf.String(), nil
}
