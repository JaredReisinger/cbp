package main

import (
	"fmt"
	"runtime"
)

var (
	version    = "(development build)"
	buildDate  string
	commitHash string
)

func printVersion() {
	fmt.Printf(`%s
version     : %s
build date  : %s
git hash    : %s
go version  : %s
go compiler : %s
platform    : %s/%s
`, name, version, buildDate, commitHash,
		runtime.Version(), runtime.Compiler,
		runtime.GOOS, runtime.GOARCH)

}
