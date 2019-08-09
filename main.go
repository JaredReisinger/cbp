//go:generate go run -v ./generate_build_info.go

package main // import "github.com/JaredReisinger/cbp"

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	fullVersion string
	rootCmd     *cobra.Command

	addr         string
	depth        int
	importPrefix string
	vcs          string
	repoPrefix   string

	testingOnly  = false
	repoTemplate = template.Must(template.New("repoInfo").Parse(`
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
`))
)

func init() {
	fullVersion := fmt.Sprintf("%s, built: %s", gitVersion, buildDate)

	rootCmd = &cobra.Command{
		Use:     "cbp [flags] repo-prefix",
		Short:   "Super-simple Golang remote import path service",
		Long:    "Super-simple Golang remote import path service. Complete documentation is available at https://github.com/JaredReisinger/cbp.",
		Args:    cobra.ExactArgs(1),
		Version: fullVersion,
		// Weird output formatting on this... is there a better way to do this?
		Example: `
  cbp https://github.com/your-name-or-org

	The simplest possible example; if hosted at "http://go.your-name-or-org.com",
	an import requests for "go.your-name-or-org.com/project" would resolve to
	"https://github.com/your-name-or-org/project".`,
		Run: run,
	}

	rootCmd.Flags().StringVarP(&addr, "addr", "a", "", "address/port on which to listen")
	rootCmd.Flags().IntVarP(&depth, "depth", "d", 0, "number of path segments after the import prefix to the repository root (defaults to 2 minus any segments from the repo-prefix, but a minimum of 1)")
	rootCmd.Flags().StringVarP(&importPrefix, "import-prefix", "i", "", "hostname and/or leading path for the custom import path; the 'Host' header of each incoming request is used by default")
	rootCmd.Flags().StringVarP(&vcs, "vcs", "", "", "version control system (VCS) name for the repos (inferred from repo-prefix, or \"git\" if uncertain)")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	repoPrefix = args[0]

	err := normalizeArgs()
	if err != nil {
		log.Fatal(err)
	}

	// log all of the options in play:
	if importPrefix == "" {
		log.Print("using the Host header on incoming requests as the import prefix")
	} else {
		log.Printf("using %q as the import prefix", importPrefix)
	}

	log.Printf("using %d path segments from import for repository roots", depth)
	log.Printf("using %q as the version control system", vcs)
	log.Printf("using %q as the repository prefix for import requests", repoPrefix)

	if addr == "" {
		log.Print("using the default address:port (:http)")
	} else {
		log.Printf("using address:port %q", addr)
	}

	if !testingOnly {
		log.Print("starting cbp...")
		http.HandleFunc("/", serveMeta)
		http.Handle("/favicon.ico", http.NotFoundHandler())
		log.Fatal(http.ListenAndServe(addr, nil))
	}
}

func normalizeArgs() (err error) {
	if repoPrefix == "" {
		err = errors.New("a repository root/prefix is required")
		return
	}
	// parse the repo prefix to count paths...(?)
	url, err := url.Parse(repoPrefix)
	if err != nil {
		return
	}

	if url.Opaque != "" {
		err = fmt.Errorf("unexpected opaque data %q; missing %q after scheme %q?", url.Opaque, "//", url.Scheme)
		return
	}

	if url.RawQuery != "" {
		log.Printf("ignoring query %q in repo prefix", url.RawQuery)
		url.RawQuery = ""
	}

	if url.Fragment != "" {
		log.Printf("ignoring fragment %q in repo prefix", url.Fragment)
		url.Fragment = ""
	}

	// we will be stripping leading "/" on incoming reqests, so we need to be
	// sure that the repo prefix url is normalized to include a trailing one.
	if !strings.HasSuffix(url.Path, "/") {
		url.Path = fmt.Sprintf("%s/", url.Path)
	}

	// save the normalized prefix
	repoPrefix = url.String()

	// default the VCS if the user didn't specify
	if vcs == "" {
		scheme := strings.TrimSuffix(url.Scheme, "+ssh")
		if scheme != "http" && scheme != "https" {
			vcs = scheme
		} else {
			vcs = "git"
		}
	}

	if depth <= 0 {
		// we could do math on the path, but it's either 2 if the path is '/', or 1 otherwise.
		if url.Path == "/" {
			depth = 2
		} else {
			depth = 1
		}
	}

	err = nil
	return
}

func serveMeta(w http.ResponseWriter, r *http.Request) {
	importPrefixAuto := r.Host
	if importPrefix != "" {
		importPrefixAuto = importPrefix
	}
	importPath, repoPath := calculatePaths(r.URL.Path, importPrefixAuto, repoPrefix)

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
