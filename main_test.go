package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCalculatePath(t *testing.T) {
	importPrefix := "go.example.com"
	repoPrefix := "ssh://git@git.example.com"

	examples := []struct {
		requestPath    string
		expectedImport string
		expectedRepo   string
	}{
		{
			"/orgname/reponame",
			"go.example.com/orgname/reponame",
			"ssh://git@git.example.com/orgname/reponame",
		},
		{
			"/orgname/reponame/subdir",
			"go.example.com/orgname/reponame/subdir",
			"ssh://git@git.example.com/orgname/reponame",
		},
		{
			"/shortpath",
			"go.example.com/shortpath",
			"ssh://git@git.example.com/shortpath",
		},
	}

	for i, ex := range examples {
		actualImport, actualRepo := calculatePaths(ex.requestPath, importPrefix, repoPrefix)

		if actualImport != ex.expectedImport {
			t.Errorf("import path [%d]: expected %s, got %s", i, ex.expectedImport, actualImport)
		}

		if actualRepo != ex.expectedRepo {
			t.Errorf("repo path [%d]: expected %s, got %s", i, ex.expectedRepo, actualRepo)
		}
	}
}

func TestValidateOptions(t *testing.T) {
	examples := []struct {
		importPrefix string
		vcs          string
		repoPrefix   string
		ok           bool
	}{
		{"", "", "", false},
		{"x", "", "", false},
		{"", "x", "", false},
		{"", "", "x", false},
		{"x", "x", "", false},
		{"x", "", "x", true}, // vcs defaults to 'git' when empty
		{"", "x", "x", false},
		{"x", "x", "x", true},
		{"/", "x", "x", false},
		{"/", "x", "/", false},
		{"x", "x", "/", false},
	}

	origUsage := flag.Usage
	defer func() { flag.Usage = origUsage }()
	flag.Usage = func() {}

	for i, ex := range examples {
		*importPrefix = ex.importPrefix
		*vcs = ex.vcs
		*repoPrefix = ex.repoPrefix

		ok := validateOptions()
		if ok != ex.ok {
			t.Errorf("options [%d]: expected %v, got %v", i, ex.ok, ok)
		}
	}
}

type fakeWriter struct {
	buf bytes.Buffer
}

func (f *fakeWriter) Header() http.Header {
	log.Fatal("unexpected")
	return http.Header{}
}

func (f *fakeWriter) Write(b []byte) (int, error) {
	// save output....
	return f.buf.Write(b)
}

func (f *fakeWriter) WriteHeader(int) {
	log.Fatal("unexpected")
}

func TestServeMeta(t *testing.T) {
	*importPrefix = "IMPORT"
	*repoPrefix = "REPO"
	*vcs = "VCS"

	examples := []struct {
		requestPath string
		importPath  string
		repoPath    string
	}{
		{"/foo", "IMPORT/foo", "REPO/foo"},
		{"/foo/bar", "IMPORT/foo/bar", "REPO/foo/bar"},
		{"/foo/bar/baz", "IMPORT/foo/bar/baz", "REPO/foo/bar"},
		{"/foo/bar/baz/quux", "IMPORT/foo/bar/baz/quux", "REPO/foo/bar"},
	}

	for i, ex := range examples {
		url, err := url.ParseRequestURI(ex.requestPath)
		if err != nil {
			t.Errorf("error in example %d, %q not parsed as a URL", i, ex.requestPath)
			continue
		}
		r := &http.Request{URL: url}
		w := &fakeWriter{}
		serveMeta(w, r)

		expectedMeta := fmt.Sprintf("<meta name=\"go-import\" content=\"%s %s %s\" />", ex.importPath, *vcs, ex.repoPath)
		if !strings.Contains(w.buf.String(), expectedMeta) {
			t.Errorf("request [%d]: %q did not generate expected meta info %q", i, ex.requestPath, expectedMeta)
		}
	}
}

// TODO: figure out better testing here?
func TestMain(t *testing.T) {
	testingOnly = true
	main()
}
