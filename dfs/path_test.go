package dfs

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestIsRootedPath(t *testing.T) {
	absolute, err := filepath.Abs(filepath.Join("synthetic", "selected"))
	if err != nil {
		t.Fatal(err)
	}

	if runtime.GOOS == "windows" && strings.HasPrefix(absolute, RootDirectory) {
		t.Fatalf("expected a drive-qualified Windows path, got %q", absolute)
	}

	cases := []struct {
		name string
		path string
		want bool
	}{
		{name: "platform absolute path", path: absolute, want: true},
		{name: "root shortcut", path: RootDirectory, want: true},
		{name: "relative path", path: filepath.Join("synthetic", "selected"), want: false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if got := IsRootedPath(test.path); got != test.want {
				t.Fatalf("IsRootedPath(%q) = %t, want %t", test.path, got, test.want)
			}
		})
	}
}
