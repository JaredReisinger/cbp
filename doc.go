// Copyright 2018 Jared Reisinger. All rights reserved.
// DO NOT EDIT THIS FILE. GENERATED BY mkdoc.sh.

// cbp is a tool for simple Golang remote import path management.
//
// Usage: cbp -prefix IMPORT-PREFIX -root REPO-ROOT [-vcs VCS]]
//
//   -help
//     	show help
//   -prefix import-prefix
//     	the import-prefix hostname for the custom import path
//   -root string
//     	the actual hosting repo for the custom import path
//   -vcs string
//     	the VCS for the repos (default "git")
//   -version
//     	show the version
//
// Examples:
//
//     cbp -prefix go.example.org -root https://github.com/example
//
//   Starts a cbp server that points requests for 'go.example.org/pkg' to
//   'https://github.com/example/pkg'.
package main
