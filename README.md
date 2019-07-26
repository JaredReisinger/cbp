# cbp

[![build status](https://img.shields.io/drone/build/JaredReisinger/cbp?logo=drone)](https://cloud.drone.io/JaredReisinger/cbp)
[![test coverage](https://img.shields.io/codecov/c/github/JaredReisinger/cbp?logo=codecov)](https://codecov.io/gh/JaredReisinger/cbp)
[![commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen)](http://commitizen.github.io/cz-cli/)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079)](https://github.com/semantic-release/semantic-release)\
[![version](https://images.microbadger.com/badges/version/jaredreisinger/cbp.svg)](https://microbadger.com/images/jaredreisinger/cbp)
[![commit](https://images.microbadger.com/badges/commit/jaredreisinger/cbp.svg)](https://microbadger.com/images/jaredreisinger/cbp)
[![layers](https://images.microbadger.com/badges/image/jaredreisinger/cbp.svg)](https://microbadger.com/images/jaredreisinger/cbp)
[![license](https://images.microbadger.com/badges/license/jaredreisinger/cbp.svg)](https://microbadger.com/images/jaredreisinger/cbp)


Super-simple Golang remote import path service.

## Usage

Let's say you have your source hosted on GitHub (`https://github.com/example/my-awesome-go-package`), but you really want the package exposed via your vanity domain, as `go.example.org/my-awesome-go-package`.  All you need to do is:

```sh
cbp -import-prefix go.example.org -repo-prefix https://github.com/example
```

Once you ensure that HTTP traffic for `go.example.org` is connected to the `cbp` started above, requests for `go.example.org/my-awesome-go-package` will result in a response containing the meta tag:

```html
<meta name="go-import" content="go.example.org/my-awesome-go-package git https://github.com/example/my-awesome-go-package">
```

_(If the block above is blank, it’s because the markdown renderer isn’t escaping the angle brackets in the HTML example—I’m looking at you, cloud.docker.com—and you should see [this README on GitHub](https://github.com/JaredReisinger/cbp/#readme) for a more-accurate rendering.)_

Or, if you have (say) a private organization-wide source-code tool hosted at `http://git.internal.example.org`, with GitHub-style `user/repo` or `team/repo` structure, and you again want that exposed (again internally) as `go.internal.example.org/...`, you can use:

```sh
cbp -import-prefix go.internal.example.org -repo-prefix http://git.internal.example.org
```

And again with proper traffic routing, requests to import `go.internal.example.org/team/project` will be redirected to `http://git.internal.example.org/team/project`.

## Caveats

* It only handles mapping for a particular import host to a *single* backend repository server.

* It assumes that all repos are "GitHub-style", meaning that there's an owner/organization and a repo name (like "JaredReisinger/cbp", or "mozilla/mig").  More directly, it assumes that the first two components of the request path are the repository root, and any subsequent path components are simply directories (sub-packages) inside the repository.

On the plus side, because of these restrictions, no advance registration of repositories is needed; the service does not need to be informed/updated when a new repo is created.  It blindly assumes that given a request for path `foo...` all it really has to do is glue the repo-prefix onto it.

## Why "cbp"?

CBP is the common initialism for the U.S. Customs and Border Protection, which deals with customs and importing, and this service is all about custom import paths.

## Developing / contributing

This repo is "commitizen-friendly", despite being a golang project and not a Node.js/npm project.  The caveat is that you need to have `commitizen` and `cz-conventional-changelog` installed locally/globally and on your path for `git-cz` to give you the interactive prompt.