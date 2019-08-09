## [1.0.2](https://github.com/JaredReisinger/cbp/compare/v1.0.1...v1.0.2) (2019-08-09)


### Bug Fixes

* **typos:** fix spelling errors in UI text ([996c917](https://github.com/JaredReisinger/cbp/commit/996c917))

## [1.0.1](https://github.com/JaredReisinger/cbp/compare/v1.0.0...v1.0.1) (2019-08-09)


### Bug Fixes

* **docker build:** ensure git information is availble for docker build ([85a0a5b](https://github.com/JaredReisinger/cbp/commit/85a0a5b))

# [1.0.0](https://github.com/JaredReisinger/cbp/compare/v0.1.0...v1.0.0) (2019-08-09)


### Continuous Integration

* **drone:** improve clone step ([479fa94](https://github.com/JaredReisinger/cbp/commit/479fa94))


### Features

* resolve all open issues (import prefix, address, depth), change CLI ([6775ddf](https://github.com/JaredReisinger/cbp/commit/6775ddf)), closes [#2](https://github.com/JaredReisinger/cbp/issues/2) [#7](https://github.com/JaredReisinger/cbp/issues/7) [#9](https://github.com/JaredReisinger/cbp/issues/9)


### BREAKING CHANGES

* **drone:** now expose port 80 instead of 9090 in the Docker image
* Old `--repo-prefix` option now superseded by required positional argument.

# [0.1.0](https://github.com/JaredReisinger/cbp/compare/v0.0.1...v0.1.0) (2019-07-25)


### Features

* use "-vcs git" by default ([a4cf12e](https://github.com/JaredReisinger/cbp/commit/a4cf12e)), closes [#1](https://github.com/JaredReisinger/cbp/issues/1)
