package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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

func TestNormalizeArgs(t *testing.T) {
	examples := []struct {
		repoIn   string
		vcsIn    string
		importIn string
		depthIn  int

		repoOut   string
		vcsOut    string
		importOut string
		depthOut  int

		ok bool
	}{
		// valid/expected args
		{"http://REPO", "", "", 0, "http://REPO/", "git", "", 2, true},
		{"https://REPO", "", "", 0, "https://REPO/", "git", "", 2, true},

		// vcs inference and overrides
		{"git://REPO", "", "", 0, "git://REPO/", "git", "", 2, true},
		{"bzr://REPO", "", "", 0, "bzr://REPO/", "bzr", "", 2, true},
		{"svn://REPO", "", "", 0, "svn://REPO/", "svn", "", 2, true},
		{"git+ssh://REPO", "", "", 0, "git+ssh://REPO/", "git", "", 2, true},
		{"bzr+ssh://REPO", "", "", 0, "bzr+ssh://REPO/", "bzr", "", 2, true},
		{"svn+ssh://REPO", "", "", 0, "svn+ssh://REPO/", "svn", "", 2, true},
		{"http://REPO", "git", "", 0, "http://REPO/", "git", "", 2, true},
		{"http://REPO", "bzr", "", 0, "http://REPO/", "bzr", "", 2, true},
		{"http://REPO", "svn", "", 0, "http://REPO/", "svn", "", 2, true},
		{"bzr://REPO", "svn", "", 0, "bzr://REPO/", "svn", "", 2, true},

		// depth inference and overrides
		{"http://REPO", "", "", -1, "http://REPO/", "git", "", 2, true},
		{"http://REPO/path", "", "", 0, "http://REPO/path/", "git", "", 1, true},
		{"http://REPO/path", "", "", 2, "http://REPO/path/", "git", "", 2, true},
		{"http://REPO/path", "", "", 3, "http://REPO/path/", "git", "", 3, true},
		{"http://REPO/path", "", "", -1, "http://REPO/path/", "git", "", 1, true},

		// allowed/ignored args
		{"http://REPO/?query=ignored", "", "", 0, "http://REPO/", "git", "", 2, true},
		{"http://REPO/#fragment-ignored", "", "", 0, "http://REPO/", "git", "", 2, true},

		// invalid args...
		{"", "", "", 0, "", "", "", 0, false},
		{"BOGUS:\ncontrol-character", "", "", 0, "", "", "", 0, false},
		{"blob:opaque-not-allowed", "", "", 0, "", "", "", 0, false},
	}

	origUsage := flag.Usage
	defer func() { flag.Usage = origUsage }()
	flag.Usage = func() {}

	for i, ex := range examples {
		repoPrefix = ex.repoIn
		vcs = ex.vcsIn
		importPrefix = ex.importIn
		depth = ex.depthIn

		err := normalizeArgs()
		if (err != nil && ex.ok) || (err == nil && !ex.ok) {
			t.Errorf("options [%d]: expected %v, got %v", i, ex.ok, err)
		}

		if err == nil && ex.ok {
			expect(t, "repoPrefix", ex.repoOut, repoPrefix)
			expect(t, "vcs", ex.vcsOut, vcs)
			expect(t, "importPrefix", ex.importOut, importPrefix)
			expectInt(t, "depthPrefix", ex.depthOut, depth)
		}
	}
}

func expect(t *testing.T, label string, expected string, actual string) {
	if actual != expected {
		t.Errorf("%s expected %q, got %q", label, expected, actual)
	}
}

func expectInt(t *testing.T, label string, expected int, actual int) {
	if actual != expected {
		t.Errorf("%s expected %d, got %d", label, expected, actual)
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
	importPrefix = "IMPORT"
	repoPrefix = "REPO"
	vcs = "VCS"

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

		expectedMeta := fmt.Sprintf("<meta name=\"go-import\" content=\"%s %s %s\" />", ex.importPath, vcs, ex.repoPath)
		if !strings.Contains(w.buf.String(), expectedMeta) {
			t.Errorf("request [%d]: %q did not generate expected meta info %q", i, ex.requestPath, expectedMeta)
		}
	}
}

// TODO: figure out better testing here?
func TestRun(t *testing.T) {
	testingOnly = true

	importPrefix = ""
	addr = ""
	run(nil, []string{"http://REPO"})

	// run(nil, []string{"\n"})

	importPrefix = "IMPORT"
	addr = ":addr"
	run(nil, []string{"http://REPO"})
}

// TODO: figure out better testing here?
func TestMain(t *testing.T) {
	testingOnly = true
	importPrefix = ""
	vcs = ""
	depth = 0
	addr = ""
	os.Args = []string{"cbp", "http://REPO"}
	main()
}
