# cbp

[![build status](https://img.shields.io/drone/build/JaredReisinger/cbp?logo=drone)](https://cloud.drone.io/JaredReisinger/cbp)
[![test coverage](https://img.shields.io/codecov/c/github/JaredReisinger/cbp?logo=codecov)](https://codecov.io/gh/JaredReisinger/cbp)
[![Go Report Card](https://goreportcard.com/badge/github.com/JaredReisinger/cbp)](https://goreportcard.com/report/github.com/JaredReisinger/cbp)
[![commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen)](http://commitizen.github.io/cz-cli/)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079)](https://github.com/semantic-release/semantic-release)\
[![version](https://images.microbadger.com/badges/version/jaredreisinger/cbp.svg)](https://microbadger.com/images/jaredreisinger/cbp)
[![commit](https://images.microbadger.com/badges/commit/jaredreisinger/cbp.svg)](https://microbadger.com/images/jaredreisinger/cbp)
[![layers](https://images.microbadger.com/badges/image/jaredreisinger/cbp.svg)](https://microbadger.com/images/jaredreisinger/cbp)
[![license](https://images.microbadger.com/badges/license/jaredreisinger/cbp.svg)](https://microbadger.com/images/jaredreisinger/cbp)

Super-simple Golang remote import path service.

## Usage

Let's say you have your source hosted on GitHub (`https://github.com/example/my-awesome-go-package`), but you really want the package exposed via your vanity domain, as `go.example.org/my-awesome-go-package`. All you need to do is:

```sh
cbp https://github.com/example
```

Once you ensure that HTTP traffic for `go.example.org` is connected to the `cbp` started above, requests for `go.example.org/my-awesome-go-package` will result in a response containing the meta tag:

```html
<meta
  name="go-import"
  content="go.example.org/my-awesome-go-package git https://github.com/example/my-awesome-go-package"
/>
```

_(If the block above is blank, it’s because the markdown renderer isn’t escaping the angle brackets in the HTML example—I’m looking at you, docker.com—you should instead see [this same README on GitHub](https://github.com/JaredReisinger/cbp/#readme) for a more-accurate rendering.)_

In the case of this example, if you were to visit `http://go.example.org/my-awesome-go-package` in a web browser, you'd see something that looks like:

> #### go.example.org/my-awesome-go-package
>
> - VCS: git
> - Repository root: https://github.com/example/my-awesome-go-package

This is the same information that's in the `meta` tag, but in an easier-to-read format.

Or, if you have (say) a private company-wide Subversion service with three levels to get to an individual project (for `/organization/team/project`), hosted at `svn://svn.internal.example.org`, and you again want that exposed (again internally) as `go.internal.example.org/...`, you can use:

```sh
cbp --vcs svn --depth 3 svn://svn.internal.example.org
```

And again with proper traffic routing, requests to import `go.internal.example.org/group/team/project` will be redirected to `svn://svn.internal.example.org/group/team/project`:

```html
<meta
  name="go-import"
  content="go.internal.example.org/group/team/project git svn://svn.internal.example.org/group/team/project"
/>
```

> #### go.internal.example.org/group/team/project
>
> - VCS: svn
> - Repository root: svn://svn.internal.example.org/group/team/project

### Options

| long                                           | short | description                                                                                                            |
| ---------------------------------------------- | ----- | ---------------------------------------------------------------------------------------------------------------------- |
| <code>\-\-addr _address_</code>                | `-a`  | address/port on which to listen (defaults to ":http")                                                                  |
| <code>\-\-depth _number_</code>                | `-d`  | number of path segments to respository root (defaults to 2, minus any segments from `--import-prefix`)                 |
| <code>\-\-help</code>                          | `-h`  | help for cbp                                                                                                           |
| <code>\-\-import-prefix _nameAndOrPath_</code> | `-i`  | hostname and/or leading path for the custom import path; the 'Host' header of each incoming request is used by default |
| <code>\-\-vcs _vcsType_</code>                 |       | version control system (VCS) name for the repos (default "git")                                                        |
| <code>\-\-version</code>                       |       | show version for cbp                                                                                                   |

### Further reading

For a detailed description of the mechanism at play, please see [the `go` documentation for remote import paths](https://golang.org/cmd/go/#hdr-Remote_import_paths).

## Caveats

- It only handles mapping for a particular import host to a _single_ backend repository server.

- It assumes (and requires) that all revision-control roots are the same number of path components long (2, by default). It has to do so because it dead-reckons the repository root rather than verifying against it or requiring pre-registration.

On the plus side, because of these restrictions, no advance registration of repositories is needed; the service does not need to be informed/updated when a new repo is created. It blindly assumes that given a request for path `foo...` all it really has to do is glue the repo-prefix onto it.

## Why "cbp"?

CBP is the common initialism for the U.S. Customs and Border Protection, which deals with customs and importing, and this service is all about custom import paths.

## Developing / contributing

This repo is "commitizen-friendly", despite being a golang project and not a Node.js/npm project. The caveat is that you need to have `commitizen` and `cz-conventional-changelog` installed locally/globally and on your path for `git-cz` to give you the interactive prompt.

If you are building locally, you will need to run `go generate` at least once to create the version information. There's no `Makefile` or `magefile.go` because that seems like overkill for something like this. Besides, the `go generate` is run by the CI, so any "real" releases will include the proper version information.

### Guidelines / Notes to self...

As I re-work the internals, here are the guidelines I'm following:

- **Defaults should accomodate the 80/20 rule.** In the majority of cases, the service will run on the root endpoint for a given hostname. Further, all modern tools send the HTTP `Host` header. Therefore, the import path prefix to strip doesn't need to be specified at all; `cbp` can simply parse the request path as the path on the VCS.

- **Testing should be robust.** Creating good, robust tests is an art. The goal is to test "what do we expect as a result", not "was this implemented a particular way". The former will detect changes in behavior, while the latter often need to be "fixed" for even minor code changes.
