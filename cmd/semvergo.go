package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"

	"github.com/adamhassel/semvergo/pkg/flags"
	git2 "github.com/adamhassel/semvergo/pkg/git"
	"github.com/adamhassel/semvergo/pkg/semver"
)

var older, newer, incMajor, incMinor, incPatch, usetags, usebranch flags.Bool
var version, prefix, suffix, prefixSeparator, suffixSeparator, gitdir flags.String

func init() {
	flag.Var(&version, "v", "version string to use")
	flag.Var(&older, "o", "Print oldest of version strings (>=2) given in arguments without incrementing anything.\n Example: 'semvergo -o v1.2.3 v2.3.4 v0.1.2' yields v0.1.2")
	flag.Var(&newer, "n", "Print newest of version strings (>=2) given in arguments without incrementing anything.\n Example: 'semvergo -n v1.2.3 v2.3.4 v0.1.2' yields v2.3.4")
	flag.Var(&incMajor, "major", "increment major version")
	flag.Var(&incMinor, "minor", "increment minor version")
	flag.Var(&incPatch, "patch", "increment patch version. This is the default if no other increments are set.")
	flag.Var(&prefix, "prefix", "prefix to add to semver string")
	flag.Var(&suffix, "suffix", "suffix to add to semver string")
	flag.Var(&prefixSeparator, "prefix-sep", "prefix separator used to separate prefix from  semver string. Used both for parsing and constructing. Default is empty")
	flag.Var(&suffixSeparator, "suffix-sep", "suffix separator used to separate semver string from suffix. Used both for parsing and constructing. Default is '-'. Changing this breaks the semver standard.")

	flag.Var(&usetags, "tags", "use latest tag on git repository as version string")
	flag.Var(&usebranch, "branch", "use branch name as suffix. When used with -tags, the version number used as input is the latest tag suffixed with the branch name")
	flag.Var(&gitdir, "gitdir", "git directory. Default is current directory.")
}

func main() {
	flag.Parse()

	var sv semver.SemVer

	if !suffixSeparator.IsSet() {
		suffixSeparator.Set("-")
	}

	sv.Presep(prefixSeparator.String())
	sv.Sufsep(suffixSeparator.String())

	switch {
	case usetags.Bool():
		dir := gitdir.String()

		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
		}
		repo, err := git.PlainOpen(dir)
		if err != nil {
			log.Fatal(err)
		}
		sv, err = git2.LatestsGitVersionTag(repo, usebranch.Bool(), suffixSeparator.String())
	case version.IsSet() && version.String() != "":
		var err error
		sv, err = semver.ParseSeparated(version.String(), prefixSeparator.String(), suffixSeparator.String())
		if err != nil {
			log.Fatal(err)
		}
	case older.Bool(), newer.Bool():
		if len(flag.Args()) < 2 {
			log.Fatal("must provide at least two version strings to compare")
		}
		fmt.Printf(CompareVersions(flag.Args(), older.Bool()).String())
		return
	}

	if incMajor.IsSet() && incMajor.Bool() {
		sv.IncrementMajor()
	}
	if incMinor.IsSet() && incMinor.Bool() {
		sv.IncrementMinor()
	}
	if incPatch.IsSet() && incPatch.Bool() {
		sv.IncrementPatch()
	}

	if !incMajor.IsSet() && !incMinor.IsSet() {
		sv.IncrementPatch()
	}

	if suffix.IsSet() {
		sv.Suffix(suffix.String())
	}
	if prefix.IsSet() {
		sv.Prefix(prefix.String())
	}

	fmt.Printf(sv.String())
}

func CompareVersions(vs []string, older bool) semver.SemVer {
	sv := make([]semver.SemVer, 0, len(vs))
	for _, v := range vs {
		s, err := semver.Parse(v)
		if err != nil {
			log.Printf("error parsing version %q: %v", v, err)
			continue
		}
		sv = append(sv, s)
	}
	if older {
		return semver.MinSlice(sv)
	}
	return semver.MaxSlice(sv)
}
