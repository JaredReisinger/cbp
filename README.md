# cbp

Super-simple Golang remote import path service.

## Usage

To provide resolution for `go.example.com/...` to a (for instance) Git repository host at `git.example.com` that uses SSH, try:

```sh
cbp --import-prefix go.example.com --vcs git --repo-prefix ssh://git@git.example.com
```
An incoming request for `go.example.com/path/name` will yield:

```
<meta name="go-import" content="go.example.com/path/name vcs ssh://git@git.example.com/path/name">
```

## Caveats

* It only handles mapping for a particular import host to a *single* backend repository server.

* It assumes that all repos are "GitHub-style", meaning that there's an owner/organization and a repo name (like "JaredReisinger/cbp", or "mozilla/mig").

On the plus side, because of these restrictions, no advance registration of repositories is needed; the service does not need to be informed/updated when a new repo is created.

## Why "cbp"?

CBP is the common initialism for the U.S. Customs and Border Protection, which deals with customs and importing, and this service is all about custom import paths.
