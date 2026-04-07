package test

import (
	"testing"

	"posiflora/backend/internal/i18n"
)

func LoadTranslator(t testing.TB, localesDir, lang string) *i18n.Translator {
	t.Helper()
	tr, err := i18n.Load(localesDir, lang)
	if err != nil {
		t.Fatalf("load translator: %v", err)
	}
	return tr
}
