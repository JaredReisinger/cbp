package main

import (
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