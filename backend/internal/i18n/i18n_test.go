package i18n_test

import (
	"os"
	"path/filepath"
	"testing"

	"posiflora/backend/internal/i18n"
)

func writeLocale(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "en.json"), []byte(content), 0o600); err != nil {
		t.Fatalf("write locale: %v", err)
	}
	return dir
}

func TestLoad_Success(t *testing.T) {
	dir := writeLocale(t, `{"hello": "Hello {{.Name}}"}`)
	tr, err := i18n.Load(dir, "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil translator")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := i18n.Load("/nonexistent", "en")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := writeLocale(t, `{invalid json}`)
	_, err := i18n.Load(dir, "en")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestRender_Success(t *testing.T) {
	dir := writeLocale(t, `{"greeting": "Hello {{.Name}}"}`)
	tr, _ := i18n.Load(dir, "en")

	result, err := tr.Render("greeting", struct{ Name string }{"World"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", result)
	}
}

func TestRender_KeyNotFound(t *testing.T) {
	dir := writeLocale(t, `{"greeting": "Hello"}`)
	tr, _ := i18n.Load(dir, "en")

	_, err := tr.Render("missing_key", nil)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	dir := writeLocale(t, `{"bad": "{{.Unclosed"}`)
	tr, _ := i18n.Load(dir, "en")

	_, err := tr.Render("bad", nil)
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestRender_NoInterpolation(t *testing.T) {
	dir := writeLocale(t, `{"static": "static text"}`)
	tr, _ := i18n.Load(dir, "en")

	result, err := tr.Render("static", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "static text" {
		t.Errorf("expected 'static text', got %q", result)
	}
}

func TestRender_MultipleKeys(t *testing.T) {
	dir := writeLocale(t, `{"a": "A", "b": "B {{.X}}"}`)
	tr, _ := i18n.Load(dir, "en")

	a, _ := tr.Render("a", nil)
	b, _ := tr.Render("b", struct{ X string }{"ok"})

	if a != "A" {
		t.Errorf("unexpected a: %q", a)
	}
	if b != "B ok" {
		t.Errorf("unexpected b: %q", b)
	}
}
