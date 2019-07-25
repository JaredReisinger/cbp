package main // import "github.com/JaredReisinger/cbp"

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

const addr = ":9090"

const repoInfo = `
<html>
  <head>
    <title>{{ .ImportPath }} - cbp</title>
	<meta name="go-import" content="{{ .ImportPath}} {{ .VCS }} {{ .RepoPath }}" />
  </head>
  <body>
    <h1>{{ .ImportPath }}</h1>
	<ul>
	  <li>VCS: {{ .VCS }}</li>
	  <li>Repo Root: {{ .RepoPath }}</li>
	</ul>
  </body>
</html>
`

var (
	importPrefix string
	vcs          string
	repoPrefix   string

	repoTemplate = template.Must(template.New("repoInfo").Parse(repoInfo))
)

func main() {
	flag.StringVar(&importPrefix, "import-prefix", "", "the hostname for the custom import path")
	flag.StringVar(&vcs, "vcs", "", "the VCS for the repos")
	flag.StringVar(&repoPrefix, "repo-prefix", "", "the actual hosting repo for the custom import path")
	flag.Parse()

	// ensure any trailing slashes have been removed from importPrefix and
	// repoPrefix... the calculated paths will *always* start with a slash.
	importPrefix = strings.TrimRight(importPrefix, "/")
	repoPrefix = strings.TrimRight(repoPrefix, "/")

	showHelp := false

	if importPrefix == "" {
		fmt.Println("-import-path required")
		showHelp = true
	}

	if vcs == "" {
		fmt.Println("-vcs required")
		showHelp = true
	}

	if repoPrefix == "" {
		fmt.Println("-repoPrefix required")
		showHelp = true
	}

	if showHelp {
		fmt.Println("")
		flag.Usage()
		return
	}

	log.Printf("Starting cbp at \"%s\"...", addr)

	http.HandleFunc("/", serveMeta)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(addr, nil))
}

func serveMeta(w http.ResponseWriter, r *http.Request) {
	importPath, repoPath := calculatePaths(r.URL.Path, importPrefix, repoPrefix)

	repoTemplate.Execute(w, map[string]string{
		"ImportPath": importPath,
		"VCS":        vcs,
		"RepoPath":   repoPath,
	})
}

func calculatePaths(requestPath string, importPrefix string, repoPrefix string) (importPath string, repoPath string) {
	pathPrefix := requestPath
	// Remember that real requests always start with "/", so we ignore that
	// character.  Also, we split into 3 since we don't use anything past the
	// second component.
	parts := strings.SplitN(pathPrefix[1:], "/", 3)

	// For zero, one, or two parts, the existing prefix (path) is good.  For
	// anything longer, we shorten it to just the first two (org/repo).
	if len(parts) > 2 {
		pathPrefix = fmt.Sprintf("/%s", strings.Join(parts[:2], "/"))
	}

	importPath = fmt.Sprintf("%s%s", importPrefix, requestPath)
	repoPath = fmt.Sprintf("%s%s", repoPrefix, pathPrefix)

	log.Printf("%s => %s => %s", requestPath, importPath, repoPath)

	return
}
