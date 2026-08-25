package xess

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

func TestEmbeddedFontAssets(t *testing.T) {
	cssBytes, err := fs.ReadFile(Static, "xess.css")
	if err != nil {
		t.Fatal(err)
	}
	css := string(cssBytes)

	for _, reference := range regexp.MustCompile(`url\("(\./static/[^"]+\.woff2)"\)`).FindAllStringSubmatch(css, -1) {
		path := strings.TrimPrefix(reference[1], "./")
		if _, err := fs.Stat(Static, path); err != nil {
			t.Errorf("font reference %q does not resolve in the embedded filesystem: %v", reference[1], err)
		}
	}

	for _, font := range []struct {
		weight string
		path   string
	}{
		{weight: "400", path: "static/inter-regular.woff2"},
		{weight: "600", path: "static/inter-semibold.woff2"},
		{weight: "800", path: "static/inter-extrabold.woff2"},
	} {
		if _, err := fs.Stat(Static, font.path); err != nil {
			t.Errorf("missing Inter %s font: %v", font.weight, err)
		}
		if !strings.Contains(css, `font-weight: `+font.weight+`;`) {
			t.Errorf("missing Inter %s font-weight declaration", font.weight)
		}
		if !strings.Contains(css, `src: url("./`+font.path+`") format("woff2");`) {
			t.Errorf("missing CSS reference to %s", font.path)
		}
	}

	license, err := fs.ReadFile(Static, "static/OFL.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(license), "SIL OPEN FONT LICENSE Version 1.1") {
		t.Error("OFL.txt does not contain the SIL Open Font License")
	}

	for _, want := range []string{
		"--body-sans-font: Inter, sans-serif;",
		"--body-title-font: Inter, sans-serif;",
		"--body-preformatted-font: Iosevka Curly Iaso, monospace;",
		"text-transform: uppercase;",
		"font-weight: 800;",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("missing typography rule %q", want)
		}
	}
}
