# semvergo
Semantic versioning with configurable prefixes/suffixes (pre-release labels) and branch-name labeling semantics.

## Installation

    go install github.com/adamhassel/semvergo

## Command line options

```
  -branch
    	use branch name as suffix. When used with -tags, the version number used as input is the latest tag suffixed with the branch name
  -gitdir value
    	git directory. Default is current directory.
  -major
    	increment major version
  -minor
    	increment minor version
  -patch
    	increment patch version. This is the default if no other increments are set.
  -prefix value
    	prefix to add to semver string
  -prefix-sep value
    	prefix separator used to separate prefix from  semver string. Used both for parsing and constructing. Default is empty
  -suffix value
    	suffix to add to semver string
  -suffix-sep value
    	suffix separator used to separate semver string from suffix. Used both for parsing and constructing. Default is '-'
  -tags
    	use latest tag on git repository as version string
  -v value
    	version string to use
```
# Examples

## Given/implied version
```
$ semvergo
0.0.1

$ semvergo -major
1.0.0

$ semvergo -suffix daily
0.0.1-daily

$ semvergo -v 7.0.54-dev -minor
7.1.0-dev

$ semvergo -v 7.0.54-dev -minor -patch
7.1.1-dev

$ semvergo -v 7.0.54-dev -prefix v
v7.1.0-dev

$ semvergo -v v7.0.54-dev
v7.0.55-dev

$ semvergo -v v7.0.54-dev -prefix '' -suffix ''
7.0.55

$ semvergo -v v7.0.54-dev -prefix '' -suffix 'daily'
7.0.55-daily

$ semvergo -v v7.0.54-dev -prefix '' -suffix 'daily' -suffix-sep '~~'
7.0.55~~daily
```

## Version based on (existing) git tags and/or branch names

```
# Show current git tags
$ git tag
v0.0.1-dev
v0.0.10-dev
v0.0.10-test
v0.0.11-dev
[...]
v0.0.19-dev

$ git checkout test
Switched to branch 'test'

# Generate next version number, retaining existing suffixes from latest version
$ semvergo -tags
v0.0.20-dev

# Generate next version number, retaining existing suffixes from latest version, remove prefix
$ semvergo -tags -prefix ''
0.0.20-dev

# Generate next version using only tags with current branch name as suffix
$ semvergo -tags -branch
v0.0.11-test

# Generate next version from all tags, but change the suffix 
$ semvergo -tags -suffix test
v0.0.20-test
```
